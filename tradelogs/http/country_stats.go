package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/gin-gonic/gin"
)

type countryStatsQuery struct {
	FromTime    uint64 `form:"from" binding:"required"`
	ToTime      uint64 `form:"to" binding:"required"`
	CountryCode string `form:"country" binding:"required,isValidCountryCode"`
}

func (ha *Server) getCountryStats(c *gin.Context) {
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

	countryStats, err := ha.storage.GetCountryStats(query.CountryCode, query.FromTime, query.ToTime)
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
