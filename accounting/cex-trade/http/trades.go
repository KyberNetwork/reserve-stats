package http

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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
				httputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
			binanceMarginTrades, err := s.bs.GetMarginTradeHistory(fromTime, toTime)
			if err != nil {
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
}

func (s *Server) getConvertToETHPrice(c *gin.Context) {
	var (
		query getSpecialTradesQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	result, err := s.bs.GetConvertToETHPrice(query.From, query.To)
	if err != nil {
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
		symbol, side, rate := convertRateToBinance(trade.InTokenAmount, ethAmount, trade.InToken, "ETH")
		tradeType := sellType
		tokenChange := trade.InTokenAmount * -1
		if side == askSide {
			tradeType = buyType
			ethAmount *= -1
			tokenChange = trade.InTokenAmount
		}
		result = append(result, ConvertTrade{
			AccountName:  trade.AccountName,
			Timestamp:    trade.Timestamp,
			Pair:         symbol,
			Type:         tradeType,
			Rate:         rate,
			ETHChange:    ethAmount,
			TokenChange:  tokenChange,
			Qty:          trade.InTokenAmount,
			Hash:         trade.TxHash,
			TakerAddress: trade.Taker,
		})
	}
	if trade.OutToken != usdt {
		symbol, side, rate := convertRateToBinance(ethAmount, trade.OutTokenAmount, "ETH", trade.OutToken)
		tradeType := sellType
		tokenChange := trade.InTokenAmount * -1
		if side == askSide {
			tradeType = buyType
			ethAmount *= -1
			tokenChange = trade.InTokenAmount
		}
		result = append(result, ConvertTrade{
			AccountName:  trade.AccountName,
			Timestamp:    trade.Timestamp,
			Pair:         symbol,
			Type:         tradeType,
			Rate:         rate,
			ETHChange:    ethAmount,
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
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	result, err := s.zs.GetConvertTradeInfo(int64(query.From), int64(query.To))
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	for _, trade := range result {
		r := process(trade)
		response = append(response, r...)
	}

	result, err = s.zs.GetBinanceConvertTradeInfo(int64(query.From), int64(query.To))
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	for _, trade := range result {
		r := processBinanceConvertTrade(trade)
		response = append(response, r...)
	}
	c.JSON(
		http.StatusOK,
		response,
	)
}

func convertRateToBinance(inAmount, outAmount float64, inToken, outToken string) (string, string, float64) {
	var (
		in, out     = 100, 100
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
		log.Printf("%d - %d", in, out)
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

func processBinanceConvertTrade(trade zerox.ConvertTradeInfo) []ConvertTrade {
	var (
		result       = []ConvertTrade{}
		ethAmount    float64
		quoteString  = []string{"BTC", "USDT", "USDC", "WETH", "ETH"}
		symbol, side string
		rate         float64
	)

	regexpString := fmt.Sprintf(".*(%s)$", strings.Join(quoteString, "|"))
	re := regexp.MustCompile(regexpString)
	// find eth amount
	quote := quoteFromOriginalSymbol(re, trade.InToken)
	inToken := strings.ReplaceAll(trade.InToken, quote, "")

	ethAmount = (trade.InTokenAmount * trade.InTokenRate) / trade.ETHRate

	// find side and rate
	if trade.IsBuyer {
		symbol, side, rate = convertRateToBinance(trade.InTokenAmount, ethAmount, inToken, "ETH")
	} else {
		symbol, side, rate = convertRateToBinance(ethAmount, trade.InTokenAmount, "ETH", inToken)
	}
	tradeType := "sell"
	tokenChange := trade.InTokenAmount * -1
	if side == "ask" {
		tradeType = "buy"
		ethAmount *= -1
		tokenChange = trade.InTokenAmount
	}
	result = append(result, ConvertTrade{
		AccountName: trade.AccountName,
		Timestamp:   trade.Timestamp,
		Pair:        symbol,
		Type:        tradeType,
		Rate:        rate,
		ETHChange:   ethAmount,
		TokenChange: tokenChange,
		Qty:         trade.InTokenAmount,
	},
	)

	return result
}

func quoteFromOriginalSymbol(re *regexp.Regexp, symbol string) string {
	res := re.FindAllStringSubmatch(symbol, -1)
	return res[0][1]
}
