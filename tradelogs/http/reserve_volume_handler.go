package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"

	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type reserveVolumeQuery struct {
	httputil.TimeRangeQueryFreq
	Asset   string `form:"asset" binding:"required"`
	Reserve string `form:"reserve" binding:"isAddress"`
}

func (sv *Server) getReserveVolume(c *gin.Context) {
	var query reserveVolumeQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	_, _, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	token, err := core.LookupToken(sv.coreSetting, query.Asset)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}

	result, err := sv.storage.GetReserveVolume(ethereum.HexToAddress(query.Reserve), token, query.From, query.To, query.Freq)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
