package http

import (
	"fmt"
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/gin-gonic/gin"
)

const (
	addressesTableName = "addresses"
)

type tradeSummaryQuery struct {
	httputil.TimeRangeQueryFreq
	Timezone uint64 `form:"timezone" binding:"required,isValidTimezone"`
}

func (sv *Server) countKYCEDAddresses(ts uint64) (uint64, error) {
	var (
		result uint64
		err    error
	)
	fromTime := timeutil.TimestampMsToTime(ts)
	// one day time
	toTime := timeutil.TimestampMsToTime(ts + 86400000)
	if err = sv.userPostgres.Get(&result, fmt.Sprintf(`SELECT COUNT(1) FROM "%s" WHERE timestamp >= $1 AND timestamp < $2`, addressesTableName), fromTime.UTC(), toTime.UTC()); err != nil {
		return result, err
	}
	return result, err
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

	_, _, err := query.Validate()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	tradeSummary, err := sv.storage.GetTradeSummary(query.From, query.To)
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
		kycedAddresses, err := sv.countKYCEDAddresses(ts)
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
