package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

const (
	maxTimeFrame     = time.Hour * 24 * 365 * 1 // 1 year
	defaultTimeFrame = time.Hour * 24           // 1 day

	huobiName   = "huobi"
	binanceName = "binance"
)

var validCEXs = map[string]struct{}{
	huobiName:   {},
	binanceName: {},
}

type getTradesQuery struct {
	httputil.TimeRangeQuery
	Exchanges []string `form:"cex"`
}

func (q *getTradesQuery) validate() (time.Time, time.Time, map[string]struct{}, error) {
	var cexs = make(map[string]struct{})

	fromTime, toTime, err := q.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)
	if err != nil {
		return time.Time{}, time.Time{}, nil, err
	}
	for _, cex := range q.Exchanges {
		_, ok := validCEXs[cex]
		if !ok {
			return time.Time{}, time.Time{}, nil, fmt.Errorf("invalid CEX %s", cex)
		}
		cexs[cex] = struct{}{}
	}
	if len(cexs) == 0 {
		cexs = validCEXs
	}

	return fromTime, toTime, cexs, nil
}

type getTradesResponse struct {
	Huobi   []huobi.TradeHistory   `json:"huobi,omitempty"`
	Binance []binance.TradeHistory `json:"binance,omitempty"`
}

// getTrades returns list of trades from centralized exchanges.
func (s *Server) getTrades(c *gin.Context) {
	var (
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

	fromTime, toTime, cexs, err := query.validate()
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	_, ok := cexs[huobiName]
	if ok {
		huobiTrades, err = s.hs.GetTradeHistory(fromTime, toTime)
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusBadRequest,
				err,
			)
			return
		}
	}

	_, ok = cexs[binanceName]
	if ok {
		binanceTrades, err = s.bs.GetTradeHistory(fromTime, toTime)
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusBadRequest,
				err,
			)
			return
		}
	}

	c.JSON(http.StatusOK, getTradesResponse{
		Huobi:   huobiTrades,
		Binance: binanceTrades,
	})
}
