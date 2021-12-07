package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

type tradeSummaryQuery struct {
	httputil.TimeRangeQuery
	Timezone int8 `form:"timezone"`
}

func (sv *Server) getTradeSummary(c *gin.Context) {
	var query tradeSummaryQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	fromTime, toTime, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	tradeSummary, err := sv.storage.GetTradeSummary(fromTime, toTime, query.Timezone)
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
		tradeSummary,
	)
}
