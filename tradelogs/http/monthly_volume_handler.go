package http

import (
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var (
	monthlyMaxTimeFrame = time.Hour * 24 * 365 * 10 // ~ 10 years
	monthlyTimeFrame    = time.Hour * 31 * 24       // 31 days
)

type monthlyVolumeQuery struct {
	httputil.TimeRangeQuery
	Reserve string `form:"reserve" binding:"isAddress"`
}

func (sv *Server) getMonthlyVolume(c *gin.Context) {
	var (
		query monthlyVolumeQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	from, to, err := query.Validate(httputil.TimeRangeQueryWithMaxTimeFrame(monthlyMaxTimeFrame), httputil.TimeRangeQueryWithDefaultTimeFrame(monthlyTimeFrame))
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	result, err := sv.storage.GetMonthlyVolume(common.HexToAddress(query.Reserve), from, to)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
