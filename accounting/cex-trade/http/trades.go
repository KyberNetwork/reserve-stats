package http

import (
	"fmt"
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

func process(trade zerox.ConvertTradeInfo) []ConvertTrade {
	var (
		result    = []ConvertTrade{}
		ethAmount float64
	)

	// find eth amount
	if trade.ETHRate == 0 {
		ethAmount = 0
	} else {
		if trade.InToken == usdt {
			ethAmount = trade.InTokenAmount / trade.ETHRate
		} else if trade.OutToken == usdt {
			ethAmount = trade.OutTokenAmount / trade.ETHRate
		} else {
			ethAmount = trade.InTokenAmount * trade.InTokenRate / trade.ETHRate
		}
	}

	// find side and rate
	if trade.InToken != usdt {
		symbol, side, rate := convertRateToBinance(trade.InTokenAmount, ethAmount, trade.InToken, eth)
		tradeType, ethChange, tokenChange := getAmountAndTypeOnchain(symbol, side, ethAmount, trade.InTokenAmount)
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
		tradeType, ethChange, tokenChange := getAmountAndTypeOnchain(symbol, side, ethAmount, trade.OutTokenAmount)
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
			ethAmount = t.InputAmount
			symbol, side, rate = convertRateToBinance(t.InputAmount, t.OutputAmount, eth, t.OutputToken)
		}
		tradeType, ethChange, tokenChange := getAmountAndTypeOnchain(symbol, side, ethAmount, qty)
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
		return response[i].Timestamp < response[j].Timestamp
	})

	if err := s.updatePricingGood(response); err != nil {
		s.sugar.Errorw("failed to get original trades", "error", err)
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	if strings.ToLower(query.Sort) == "desc" {
		sort.Slice(response, func(i, j int) bool {
			return response[i].Timestamp > response[j].Timestamp
		})
	}

	c.JSON(
		http.StatusOK,
		response,
	)
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
				if trade.ETHRate != 0 {
					ethAmount = (inTokenAmount * trade.InTokenRate) / trade.ETHRate
				}

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

func (s *Server) detectPricingGood(pair string, onchainTrades, rebalanceTrades []*ConvertTrade) (float64, bool) {
	var (
		pnlRate     float64
		pricingGood bool
	)
	const pnlGoodConstValue = 1.001
	if len(rebalanceTrades) > 0 {
		s.sugar.Infow("update pricing good or not", "pair", pair)
		// detect is pricing good
		lastOnchainTrades := onchainTrades[len(onchainTrades)-1]
		onchainAVGPrice := avgPrice(onchainTrades)
		rebalanceAVGPrice := avgPrice(rebalanceTrades)
		s.sugar.Infow("price", "onchain", onchainAVGPrice, "rebalance", rebalanceAVGPrice)

		if lastOnchainTrades.Type == buyType {
			if onchainAVGPrice != 0 {
				pnlRate = rebalanceAVGPrice / onchainAVGPrice
			}
		} else {
			if rebalanceAVGPrice != 0 {
				pnlRate = onchainAVGPrice / rebalanceAVGPrice
			}
		}
		s.sugar.Infow("rate", "pnlRate", pnlRate)
	} else {
		pnlRate = 0
	}
	if pnlRate > pnlGoodConstValue {
		pricingGood = true
	}
	return pnlRate, pricingGood
}

func (s *Server) updatePricingGood(trades []ConvertTrade) error {
	var (
		onchainTrades     = make(map[string][]*ConvertTrade) // save the index of the convert trades
		rebalanceTrades   = make(map[string][]*ConvertTrade) // save the index of the convert trades
		rebalancePair     string
		lastOnchainTrades ConvertTrade
	)
	s.sugar.Infow("update pricing good", "length", len(trades))
	for i, t := range trades {
		if t.AccountName == "0xRFQ" {
			s.sugar.Infow("check", "pair", t.Pair)
			ok := false
			if o, e := onchainTrades["ETHUSDT"]; e && len(rebalanceTrades["ETHUSDT"]) > 0 { // edge case on USDT, if trade ETHUSDT will match with any on chain trades
				rebalancePair = "ETHUSDT"
				ok = true
				lastOnchainTrades = *o[len(o)-1]
			}
			if o, exist := onchainTrades[t.Pair]; exist || ok {
				if exist {
					s.sugar.Debugw("pair", "pair", t.Pair, "len", len(o))
					lastOnchainTrades = *o[len(o)-1]
					if t.Type == lastOnchainTrades.Type { // if 2 on-chain have the same type, it could count as one trade
						timeDiff := t.Timestamp - lastOnchainTrades.Timestamp
						s.sugar.Infow("timediff", "diff", timeDiff, "diff number", timeDiff-30*time.Minute.Milliseconds())
						if timeDiff < 30*time.Minute.Milliseconds() { // if 2 trades happen within 30 minutes, we combine it
							onchainTrades[t.Pair] = append(onchainTrades[t.Pair], &trades[i])
							continue
						}
					} else { // if 2 on-chain trades has different types, they rebalance themselves
						rebalanceTrades[t.Pair] = append(rebalanceTrades[t.Pair], &trades[i])
					}
					rebalancePair = t.Pair
				}
				pnlRate, pricingGood := s.detectPricingGood(rebalancePair, onchainTrades[rebalancePair], rebalanceTrades[rebalancePair])
				for _, tt := range onchainTrades[t.Pair] {
					tt.PricingGood = pricingGood
					if pnlRate != 0 {
						tt.PnLBPS = pnlRate - 1
					}
				}

				onchainTrades[t.Pair] = []*ConvertTrade{&trades[i]} // reset
				rebalanceTrades[t.Pair] = []*ConvertTrade{}         // reset
				continue
			}
			onchainTrades[t.Pair] = []*ConvertTrade{&trades[i]} // init
			continue
		}
		rebalanceTrades[t.Pair] = append(rebalanceTrades[t.Pair], &trades[i])
	}

	// calculate all the left trade
	for pair, t := range onchainTrades {
		s.sugar.Infow("left over pair", "name", pair)
		pnlRate, pricingGood := s.detectPricingGood(pair, t, rebalanceTrades[pair])
		for index := range t {
			t[index].PricingGood = pricingGood
			if pnlRate != 0 {
				t[index].PnLBPS = pnlRate - 1
			}
		}
	}
	return nil
}
