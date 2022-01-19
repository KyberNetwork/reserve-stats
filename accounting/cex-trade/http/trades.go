package http

import (
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	_ "github.com/KyberNetwork/reserve-stats/accounting/common/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/accounting/zerox"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

const (
	maxTimeFrame     = time.Hour * 24 * 30 // 30 days
	defaultTimeFrame = time.Hour * 24      // 1 day
	usdt             = "USDT"
	eth              = "ETH"
	weth             = "WETH"
	sellType         = "sell"
	buyType          = "buy"
	askSide          = "ask"
)

type getTradesQuery struct {
	httputil.TimeRangeQuery
	Exchanges []string `form:"cex"`
}

type getTradesResponse struct {
	Huobi   map[string][]huobi.TradeHistory   `json:"huobi,omitempty"`
	Binance map[string][]binance.TradeHistory `json:"binance,omitempty"`
}

// getTrades returns list of trades from centralized exchanges.
func (s *Server) getTrades(c *gin.Context) {
	var (
		logger        = s.sugar.With("func", caller.GetCurrentFunctionName())
		query         getTradesQuery
		huobiTrades   = make(map[string][]huobi.TradeHistory)
		binanceTrades = make(map[string][]binance.TradeHistory) // map account with its trades
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		s.sugar.Errorw("failed to validate query", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if len(query.Exchanges) == 0 {
		query.Exchanges = []string{
			common.Huobi.String(),
			common.Binance.String()}
	}

	fromTime, toTime, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)

	if err != nil {
		s.sugar.Errorw("faield to validate time range query", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	logger = logger.With("from", fromTime, "to", toTime, "exchanges", query.Exchanges)
	logger.Debug("querying trades from database")

	for _, cex := range query.Exchanges {
		switch cex {
		case common.Huobi.String():
			huobiTrades, err = s.hs.GetTradeHistory(fromTime, toTime)
			if err != nil {
				s.sugar.Errorw("failed to get huobi trade history", "error", err)
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
		case common.Binance.String():
			binanceTrades, err = s.bs.GetTradeHistory(fromTime, toTime)
			if err != nil {
				s.sugar.Errorw("failed to get binance trade history", "error", err)
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
			binanceMarginTrades, err := s.bs.GetMarginTradeHistory(fromTime, toTime)
			if err != nil {
				s.sugar.Errorw("failed to get binance margin trade history", "error", err)
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
			for account := range binanceMarginTrades {
				binanceTrades[account] = append(binanceTrades[account], binanceMarginTrades[account]...) // append margin trades into spot trades
			}
		}
	}

	c.JSON(http.StatusOK, getTradesResponse{
		Huobi:   huobiTrades,
		Binance: binanceTrades,
	})
}

type getSpecialTradesQuery struct {
	httputil.TimeRangeQuery
	Sort string `form:"sort"`
}

func (s *Server) getConvertToETHPrice(c *gin.Context) {
	var (
		query getSpecialTradesQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		s.sugar.Errorw("failed to validate query", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	result, err := s.bs.GetConvertToETHPrice(query.From, query.To)
	if err != nil {
		s.sugar.Errorw("failed to get convert eth price", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(
		http.StatusOK,
		result,
	)
}

// ConvertTrade ...
type ConvertTrade struct {
	Timestamp    int64   `json:"timestamp"`
	Rate         float64 `json:"rate"`
	AccountName  string  `json:"account_name"`
	Pair         string  `json:"pair"`
	Type         string  `json:"type"`
	Qty          float64 `json:"qty"` // amount of base token
	ETHChange    float64 `json:"eth_change"`
	TokenChange  float64 `json:"token_change"`
	Hash         string  `json:"hash"`
	TakerAddress string  `json:"taker_address"`
}

func process(trade zerox.ConvertTradeInfo) []ConvertTrade {
	var (
		result    = []ConvertTrade{}
		ethAmount float64
	)

	// find eth amount
	if trade.InToken == usdt {
		ethAmount = trade.InTokenAmount / trade.ETHRate
	} else if trade.OutToken == usdt {
		ethAmount = trade.OutTokenAmount / trade.ETHRate
	} else {
		ethAmount = trade.InTokenAmount * trade.InTokenRate / trade.ETHRate
	}

	// find side and rate
	if trade.InToken != usdt {
		symbol, side, rate := convertRateToBinance(trade.InTokenAmount, ethAmount, trade.InToken, eth)
		tradeType, ethChange, tokenChange := getAmountAndType(symbol, side, ethAmount, trade.InTokenAmount)
		result = append(result, ConvertTrade{
			AccountName:  trade.AccountName,
			Timestamp:    trade.Timestamp,
			Pair:         symbol,
			Type:         tradeType,
			Rate:         rate,
			ETHChange:    ethChange,
			TokenChange:  tokenChange,
			Qty:          trade.InTokenAmount,
			Hash:         trade.TxHash,
			TakerAddress: trade.Taker,
		})
	}
	if trade.OutToken != usdt {
		symbol, side, rate := convertRateToBinance(ethAmount, trade.OutTokenAmount, eth, trade.OutToken)
		tradeType, ethChange, tokenChange := getAmountAndType(symbol, side, ethAmount, trade.OutTokenAmount)
		result = append(result, ConvertTrade{
			AccountName:  trade.AccountName,
			Timestamp:    trade.Timestamp,
			Pair:         symbol,
			Type:         tradeType,
			Rate:         rate,
			ETHChange:    ethChange,
			TokenChange:  tokenChange,
			Qty:          trade.OutTokenAmount,
			Hash:         trade.TxHash,
			TakerAddress: trade.Taker,
		})
	}

	return result
}

func (s *Server) getConvertTrades(c *gin.Context) {
	var (
		query    getSpecialTradesQuery
		response []ConvertTrade
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		s.sugar.Errorw("failed to validate query", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	// on chain convert trades
	result, err := s.zs.GetConvertTradeInfo(int64(query.From), int64(query.To))
	if err != nil {
		s.sugar.Errorw("failed to get convert trade info", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	zeroxTrades, err := s.zs.Get0xTrades(int64(query.From), int64(query.To))
	if err != nil {
		s.sugar.Errorw("failed to get zerox trades", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	var r []ConvertTrade
	// trade with ETH already
	for _, t := range zeroxTrades {
		if t.InputToken != weth && t.OutputToken != weth {
			continue
		}
		qty := t.InputAmount
		ethAmount := t.OutputAmount
		symbol, side, rate := convertRateToBinance(t.InputAmount, t.OutputAmount, t.InputToken, eth)
		if t.InputToken == weth {
			qty = t.OutputAmount
			symbol, side, rate = convertRateToBinance(t.InputAmount, t.OutputAmount, eth, t.OutputToken)
		}
		tradeType, ethChange, tokenChange := getAmountAndType(symbol, side, ethAmount, qty)
		r = append(r, ConvertTrade{
			AccountName:  "0xRFQ",
			Timestamp:    t.Timestamp * 1000,
			Rate:         rate,
			Pair:         symbol,
			Type:         tradeType,
			Qty:          qty,
			ETHChange:    ethChange,
			TokenChange:  tokenChange,
			Hash:         t.Tx,
			TakerAddress: t.TakerAddress,
		})
	}
	response = append(response, r...)

	for _, trade := range result {
		r := process(trade)
		response = append(response, r...)
	}

	// off chain (binance) convert trades
	result, err = s.zs.GetBinanceConvertTradeInfo(int64(query.From), int64(query.To))
	if err != nil {
		s.sugar.Errorw("failed to get binance convert trade info", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	fromTime := time.UnixMilli(int64(query.From))
	toTime := time.UnixMilli(int64(query.To))
	originalTrades, err := s.bs.GetTradeHistory(fromTime, toTime)
	if err != nil {
		s.sugar.Errorw("failed to get original trades", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	convertTrades := []ConvertTrade{}
	var (
		ethChange, tokenChange float64
		tradeType              string
	)
	for accountName, oTrades := range originalTrades {
		for _, t := range oTrades {
			if strings.HasPrefix(t.Symbol, eth) || strings.HasSuffix(t.Symbol, eth) {
				rate, err := strconv.ParseFloat(t.Price, 64)
				if err != nil {
					s.sugar.Errorw("failed to parse rate", "err", err)
					break
				}
				qty, err := strconv.ParseFloat(t.Quantity, 64)
				if err != nil {
					s.sugar.Errorw("failed to parse quantity", "err", err)
					break
				}
				tradeType = "sell"
				if t.IsBuyer {
					tradeType = "buy"
				}
				if strings.HasPrefix(t.Symbol, eth) {
					ethChange = qty
					inTokenAmount := qty * rate
					tokenChange = inTokenAmount * -1
					if !t.IsBuyer {
						ethChange *= -1
						tokenChange = inTokenAmount
					}
				} else {
					inTokenAmount := qty
					ethAmount := inTokenAmount * rate
					tokenChange = inTokenAmount * -1
					ethChange = ethAmount
					if t.IsBuyer {
						ethChange = ethChange * -1
						tokenChange = inTokenAmount
					}
				}
				convertTrades = append(convertTrades, ConvertTrade{
					AccountName: accountName,
					Timestamp:   int64(t.Time),
					Pair:        t.Symbol,
					Type:        tradeType,
					Rate:        rate,
					ETHChange:   ethChange,
					TokenChange: tokenChange,
					Qty:         qty,
				})
				continue
			}
		}
	}
	response = append(response, convertTrades...)
	for _, trade := range result {
		r := s.processBinanceConvertTrade(trade, originalTrades)
		response = append(response, r...)
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Timestamp > response[j].Timestamp
	})

	if strings.ToLower(query.Sort) == "asc" {
		sort.Slice(response, func(i, j int) bool {
			return response[i].Timestamp < response[j].Timestamp
		})
	}

	c.JSON(
		http.StatusOK,
		response,
	)
}

func convertRateToBinance(inAmount, outAmount float64, inToken, outToken string) (string, string, float64) {
	var (
		in, out     = math.MaxInt64, math.MaxInt64
		quoteTokens = []string{"DAI", "USDT", "BUSD", "USDC", "BTC", "WBTC", "WETH", "ETH"}
	)
	for i, t := range quoteTokens {
		if inToken == t {
			in = i
			continue
		}
		if outToken == t {
			out = i
			continue
		}
	}
	if in == out {
		return "", "", 0
	}
	if in < out {
		symbol := fmt.Sprintf("%s%s", outToken, inToken)
		side := "ask"
		rate := inAmount / outAmount
		return symbol, side, rate
	}
	symbol := fmt.Sprintf("%s%s", inToken, outToken)
	side := "bid"
	rate := outAmount / inAmount
	return symbol, side, rate
}

// return ethChange, tokenChange and ty
func getAmountAndType(symbol, side string, ethAmount, qty float64) (string, float64, float64) {
	var (
		ethChange, tokenChange float64
	)
	tradeType := sellType
	tokenChange = qty * -1
	ethChange = ethAmount
	if side == askSide {
		tradeType = buyType
		if strings.HasSuffix(symbol, eth) {
			ethChange = ethAmount * -1
			tokenChange = qty
		}
	} else if strings.HasPrefix(symbol, eth) {
		ethChange = ethAmount * -1
		tokenChange = qty
	}
	return tradeType, ethChange, tokenChange
}

func (s *Server) processBinanceConvertTrade(trade zerox.ConvertTradeInfo, originalTrades map[string][]binance.TradeHistory) []ConvertTrade {
	var (
		result       = []ConvertTrade{}
		ethAmount    float64
		quoteString  = []string{"DAI", "USDT", "BUSD", "USDC", "BTC", "WBTC", "WETH", "ETH"}
		symbol, side string
		rate         float64
	)
	regexpString := fmt.Sprintf(".*(%s)$", strings.Join(quoteString, "|"))
	re := regexp.MustCompile(regexpString)

	for _, oTrades := range originalTrades {
		for _, t := range oTrades {
			if t.Time == uint64(trade.Timestamp) {
				// find eth amount
				quote := quoteFromOriginalSymbol(re, t.Symbol)
				inToken := strings.TrimSuffix(t.Symbol, quote)
				inTokenAmount, err := strconv.ParseFloat(t.Quantity, 64)
				if err != nil {
					s.sugar.Errorw("failed to parse token amount", "err", err)
					break
				}
				ethAmount = (inTokenAmount * trade.InTokenRate) / trade.ETHRate

				// find side and rate
				if !trade.IsBuyer {
					symbol, side, rate = convertRateToBinance(inTokenAmount, ethAmount, inToken, eth)
				} else {
					symbol, side, rate = convertRateToBinance(ethAmount, inTokenAmount, eth, inToken)
				}
				tradeType, ethChange, tokenChange := getAmountAndType(symbol, side, ethAmount, inTokenAmount)
				result = append(result, ConvertTrade{
					AccountName: trade.AccountName, //
					Timestamp:   trade.Timestamp,
					Pair:        symbol,
					Type:        tradeType,
					Rate:        rate,
					ETHChange:   ethChange,
					TokenChange: tokenChange,
					Qty:         inTokenAmount,
				},
				)
			}
		}
	}

	return result
}

func quoteFromOriginalSymbol(re *regexp.Regexp, symbol string) string {
	res := re.FindAllStringSubmatch(symbol, -1)
	if len(res) == 0 {
		return ""
	}
	return res[0][1]
}
