package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type countryStatsQuery struct {
	FromTime    uint64 `form:"fromTime" binding:"required"`
	ToTime      uint64 `form:"toTime" binding:"required"`
	CountryCode string `form:"country" binding:"required,isValidCountryCod"`
	Timezone    int64  `form:"timezone" binding:"required,isSupportedTimezone"`
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
	fromTime := time.Unix(0, int64(query.FromTime)*int64(time.Millisecond))
	toTime := time.Unix(0, int64(query.ToTime)*int64(time.Millisecond))

	countryStats, err := ha.storage.GetCountryStats(query.CountryCode, query.Timezone, fromTime, toTime)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		countryStats,
	)
}
