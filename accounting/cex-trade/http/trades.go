package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	_ "github.com/KyberNetwork/reserve-stats/accounting/common/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

const maxTimeFrame = time.Hour * 24 * 365 * 1 // 1 year

type getTradesQuery struct {
	httputil.TimeRangeQuery
	Exchanges []string `form:"cex" binding:"dive,isValidCEXName"`
}

type getTradesResponse struct {
	Huobi   []huobi.TradeHistory   `json:"huobi,omitempty"`
	Binance []binance.TradeHistory `json:"binance,omitempty"`
}

// getTrades returns list of trades from centralized exchanges.
func (s *Server) getTrades(c *gin.Context) {
	var (
		logger        = s.sugar.With("func", "accounting/cex-trade/http/Server.getTrades")
		query         getTradesQuery
		huobiTrades   []huobi.TradeHistory
		binanceTrades []binance.TradeHistory
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
		}
	}

	c.JSON(http.StatusOK, getTradesResponse{
		Huobi:   huobiTrades,
		Binance: binanceTrades,
	})
}
