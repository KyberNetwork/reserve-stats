package http

import (
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

const (
	maxTimeFrame     = time.Hour * 24 * 365 * 3 // 3 years
	defaultTimeFrame = time.Hour * 24 * 3       // 3 days
)

func (sv *Server) getIntegrationVolume(c *gin.Context) {
	var query httputil.TimeRangeQuery
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

	integrationVolume, err := sv.storage.GetIntegrationVolume(fromTime, toTime)
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
		integrationVolume,
	)
}
