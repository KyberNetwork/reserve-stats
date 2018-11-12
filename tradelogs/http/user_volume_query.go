package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type userVolumeQuery struct {
	httputil.TimeRangeQueryFreq
	UserAddress string `form:"userAddr" binding:"required,isAddress"`
}

func (sv *Server) getUserVolume(c *gin.Context) {
	var query userVolumeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	_, _, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	volume, err := sv.storage.GetUserVolume(ethereum.HexToAddress(query.UserAddress), query.From, query.To, query.Freq)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, volume)
}
