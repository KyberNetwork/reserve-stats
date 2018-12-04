package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/gin-gonic/gin"
)

type countryStatsQuery struct {
	httputil.TimeRangeQuery
	CountryCode string `form:"country" binding:"required,isValidCountryCode"`
	Timezone    int8   `form:"timezone" binding:"isSupportedTimezone"`
}

func (sv *Server) getCountryStats(c *gin.Context) {
	var query countryStatsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	from, to, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	countryCode := query.CountryCode
	if countryCode == common.UnknownCountry {
		countryCode = ""
	}

	countryStats, err := sv.storage.GetCountryStats(countryCode, from, to, query.Timezone)
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
		countryStats,
	)
}
