package http

import (
	"errors"
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

type tradeSummaryQuery struct {
	From     uint64 `form:"fromTime"`
	To       uint64 `form:"toTime"`
	Timezone uint64 `form:"timezone" binding:"required,isValidTimezone"`
}

func (sv *Server) getTradeSummary(c *gin.Context) {
	var (
		query  tradeSummaryQuery
		logger = sv.sugar.With("func", "tradelogs/http/Server.getTradeSummary")
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	if !timeValidation(&query.From, &query.To, c, logger) {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			errors.New("time input is not valid"),
		)
		return
	}
}
