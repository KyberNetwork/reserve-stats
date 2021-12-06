package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type tokenHeatmapQuery struct {
	httputil.TimeRangeQuery
	Asset    string `form:"asset" binding:"required"`
	Timezone int8   `form:"timezone"`
}

func (sv *Server) getTokenHeatMap(c *gin.Context) {
	var (
		query tokenHeatmapQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
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

	asset := common.HexToAddress(query.Asset)

	heatmap, err := sv.storage.GetTokenHeatmap(asset, from, to, query.Timezone)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, heatmap)
}
