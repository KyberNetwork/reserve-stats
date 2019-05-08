package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type assetVolumeQuery struct {
	httputil.TimeRangeQueryFreq
	Asset string `form:"asset" binding:"required,isAddress"`
}

func (sv *Server) getAssetVolume(c *gin.Context) {
	var (
		query assetVolumeQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	from, to, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	token := common.HexToAddress(query.Asset)

	result, err := sv.storage.GetAssetVolume(token, from, to, query.Freq)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, result)
}
