package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

type tradeSummaryQuery struct {
	httputil.TimeRangeQueryFreq
	// Timezone uint64 `form:"timezone" binding:"required"`
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

	fromTime, toTime, err := query.Validate()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	tradeSummary, err := sv.storage.GetTradeSummary(fromTime, toTime)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	// update kyced addresses
	for ts, trade := range tradeSummary {
		kycedAddresses, err := sv.userPostgres.CountKYCEDAddresses(ts)
		if err != nil {
			httputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
			return
		}
		trade.KYCEDAddresses = kycedAddresses
	}
	c.JSON(
		http.StatusOK,
		tradeSummary,
	)
}
