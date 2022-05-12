package zerox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/hasura/go-graphql-client"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/marketdata"
)

// TradelogClient ...
type TradelogClient struct {
	sugar            *zap.SugaredLogger
	httpClient       *http.Client
	graphqlClient    *graphql.Client
	baseURL          string
	binanceClient    *binance.Client
	marketDataClient *marketdata.Client
	symbols          []string
}

func updateBinaceSupportedSymbol(binanceClient *binance.Client) ([]string, error) {
	var (
		symbols []string
	)
	exchangeInfo, err := binanceClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return nil, err
	}
	for _, s := range exchangeInfo.Symbols {
		symbols = append(symbols, s.Symbol)
	}
	return symbols, nil
}

// NewZeroXTradelogClient ...
func NewZeroXTradelogClient(baseURL string, binanceClient *binance.Client, marketDataClient *marketdata.Client, sugar *zap.SugaredLogger) *TradelogClient {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	graphqlClient := graphql.NewClient("https://gateway.thegraph.com/api/6c2fe0823843837e65ab636c0a861158/subgraphs/id/0x36c057dd1850fad3c075ba83105e67d2448dedaf-0", &client)

	symbols, err := updateBinaceSupportedSymbol(binanceClient)
	if err != nil {
		sugar.Errorw("failed to get symbols", "error", err)
	}
	return &TradelogClient{
		graphqlClient:    graphqlClient,
		httpClient:       &client,
		baseURL:          baseURL,
		sugar:            sugar,
		binanceClient:    binanceClient,
		marketDataClient: marketDataClient,
		symbols:          symbols,
	}
}

// GetTradelogsFromHTTP ...
func (z *TradelogClient) GetTradelogsFromHTTP(fromTime, toTime time.Time) ([]Tradelog, error) {
	var (
		tradelogsResponse TradelogsResponse
		tradelogs         []Tradelog
	)
	for {
		z.sugar.Infow("get tradelogs", "from time", fromTime.Unix(), "to time", toTime.Unix())
		requestBody := []byte(`{"query": "{ maker(id: \"0xbc33a1f908612640f2849b56b67a4de4d179c151\") { nativeOrderFills( orderBy: timestamp, orderDirection: desc, first: 1000, skip: 0, where: { timestamp_lt: ` + fmt.Sprintf("%d", toTime.Unix()) + `, timestamp_gt: ` + fmt.Sprintf("%d", fromTime.Unix()) + ` } ) { timestamp taker{ id } transaction{ id } inputToken{ id, symbol, decimals } outputToken{ id, symbol, decimals } inputTokenAmount, outputTokenAmount } } }"}`)
		endpoint := fmt.Sprintf("%s/0x-exchange", z.baseURL)
		request, err := http.NewRequest(
			http.MethodPost,
			endpoint,
			bytes.NewBuffer(requestBody),
		)
		if err != nil {
			return nil, err
		}
		resp, err := z.httpClient.Do(request)
		if err != nil {
			return nil, err
		}
		switch resp.StatusCode {
		case http.StatusOK:
			response, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(response, &tradelogsResponse); err != nil {
				return nil, err
			}
		default:
			z.sugar.Errorw("failed get tradelogs with http code", "error code", resp.StatusCode)
		}
		tradelogs = append(tradelogs, tradelogsResponse.Data.Maker.NativeOrderFills...)
		if len(tradelogsResponse.Data.Maker.NativeOrderFills) < 1000 {
			break
		}
		timestampVal, err := strconv.ParseInt(tradelogsResponse.Data.Maker.NativeOrderFills[0].Timestamp, 10, 64)
		if err != nil {
			return tradelogs, err
		}
		fromTime = time.Unix(timestampVal+1, 0)
	}
	return tradelogs, nil
}

// ConvertTrades ...
type ConvertTrades struct {
	OriginalSymbols []string
	Symbols         []string
	Prices          []float64
	Timestamps      []int64
	InToken         []string
	InTokenAmount   []float64
	OutToken        []string
	OutTokenAmount  []float64
	OriginalTrades  [][]byte
	Trades          [][]byte
}

func (z *TradelogClient) checkSymbolIsSupported(symbol string) bool {
	for _, s := range z.symbols {
		if symbol == s {
			return true
		}
	}
	return false
}

// ConvertTrades convert trades to trades with ETH
func (z *TradelogClient) ConvertTrades(tradelogs []Tradelog) (ConvertTrades, error) {
	var (
		delta  int64 = 60 * 60 * 1000 // default 1 hour
		result ConvertTrades
	)
	for _, trade := range tradelogs {
		endTime, err := strconv.ParseInt(trade.Timestamp, 10, 64)
		if err != nil {
			return result, err
		}
		endTime *= 1000
		startTime := endTime - delta + 1000
		originalSymbol := trade.InputToken.Symbol + trade.OutputToken.Symbol
		symbol := "ETHUSDT"
		result, err = z.updateTrade(result, originalSymbol, symbol, startTime, endTime, trade)
		if err != nil {
			return result, err
		}
		if trade.InputToken.Symbol == "USDT" || trade.OutputToken.Symbol == "USDT" {
			continue
		}
		symbol = trade.InputToken.Symbol + "USDT"
		if trade.InputToken.Symbol == "DAI" {
			symbol = "USDT" + trade.InputToken.Symbol
		}
		if !z.checkSymbolIsSupported(symbol) { // WBTC does not have pair with USDT on binance
			symbol = trade.InputToken.Symbol + "ETH"
		}
		result, err = z.updateTrade(result, originalSymbol, symbol, startTime, endTime, trade)
		if err != nil {
			return result, err
		}

		symbol = trade.OutputToken.Symbol + "USDT"
		if trade.OutputToken.Symbol == "DAI" {
			symbol = "USDT" + trade.OutputToken.Symbol
		}
		if !z.checkSymbolIsSupported(symbol) { // WBTC does not have pair with USDT on binance
			symbol = trade.OutputToken.Symbol + "ETH"
		}
		result, err = z.updateTrade(result, originalSymbol, symbol, startTime, endTime, trade)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (z *TradelogClient) updateTrade(result ConvertTrades, originalSymbol, symbol string, startTime, endTime int64, trade Tradelog) (ConvertTrades, error) {
	z.sugar.Infow("update trade", "symbol", symbol, "start time", startTime, "end time", endTime)
	aggTrades, err := z.getGetAggregatedTradesWithRetry(symbol, startTime, endTime)
	if err != nil {
		return result, err
	}
	if len(aggTrades) == 0 {
		z.sugar.Warnw("there is no trade for 1 hour", "symbol", symbol)
		return result, nil
	}
	atrade := aggTrades[0]
	result.OriginalSymbols = append(result.OriginalSymbols, originalSymbol)
	result.Symbols = append(result.Symbols, symbol)
	priceF, err := strconv.ParseFloat(atrade.Price, 64)
	if err != nil {
		return result, err
	}
	result.Prices = append(result.Prices, priceF)
	result.Timestamps = append(result.Timestamps, endTime)
	tradeB, err := json.Marshal(atrade)
	if err != nil {
		return result, err
	}
	result.Trades = append(result.Trades, tradeB)
	originalTradeB, err := json.Marshal(trade)
	if err != nil {
		return result, err
	}
	result.OriginalTrades = append(result.OriginalTrades, originalTradeB)
	result.InToken = append(result.InToken, trade.InputToken.Symbol)
	result.OutToken = append(result.OutToken, trade.OutputToken.Symbol)
	inTokenAmount, err := z.convertAmount(trade.InputTokenAmount, trade.InputToken.Decimals)
	if err != nil {
		return result, err
	}
	outTokenAmount, err := z.convertAmount(trade.OutputTokenAmount, trade.OutputToken.Decimals)
	if err != nil {
		return result, err
	}
	result.InTokenAmount = append(result.InTokenAmount, inTokenAmount)
	result.OutTokenAmount = append(result.OutTokenAmount, outTokenAmount)
	return result, nil
}

func (z *TradelogClient) getGetAggregatedTradesWithRetry(symbol string, startTime, endTime int64) ([]*binance.AggTrade, error) {
	var (
		aggregatedTrades []*binance.AggTrade
		err              error
	)
	attempt := 3
	for a := 0; a < attempt; a++ {
		aggTrades, err := z.binanceClient.NewAggTradesService().Symbol(symbol).StartTime(startTime).EndTime(endTime).Do(context.Background())
		if err != nil {
			return nil, err
		}
		aggregatedTrades = append(aggregatedTrades, aggTrades...)
	}
	return aggregatedTrades, err
}

func (z *TradelogClient) convertAmount(amountStr, decimalsStr string) (float64, error) {
	decimals, err := strconv.ParseInt(decimalsStr, 10, 64)
	if err != nil {
		z.sugar.Errorw("failed to parse decimals", "error", err)
		return 0, err
	}
	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	amountFloat, ok := big.NewFloat(0).SetString(amountStr)
	if !ok {
		z.sugar.Errorw("failed to parse amount", "amount str", amountStr)
		return 0, fmt.Errorf("failed to parse amount")
	}
	amountF := new(big.Float).Quo(amountFloat, big.NewFloat(0).SetInt(pow))
	amount, _ := amountF.Float64()
	return amount, nil
}
