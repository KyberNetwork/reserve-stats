package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type reserveVolumeQuery struct {
	httputil.TimeRangeQueryFreq
	Asset   string `form:"asset" binding:"required"`
	Reserve string `form:"reserve"`
}

func (sv *Server) getReserveVolume(c *gin.Context) {
	var query reserveVolumeQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	from, to, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	token := ethereum.HexToAddress(query.Asset)

	result, err := sv.storage.GetReserveVolume(ethereum.HexToAddress(query.Reserve), token,
		from, to, query.Freq)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
