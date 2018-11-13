package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

type tokenHeatmapQuery struct {
	httputil.TimeRangeQuery
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
}
