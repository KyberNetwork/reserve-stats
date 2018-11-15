package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/gin-gonic/gin"
)

type countryStatsQuery struct {
	httputil.TimeRangeQueryFreq
	CountryCode string `form:"country" binding:"required,isValidCountryCode"`
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

	from, to, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	countryCode := query.CountryCode
	if countryCode == common.UnknownCountry {
		countryCode = ""
	}

	countryStats, err := sv.storage.GetCountryStats(countryCode, from, to)
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
