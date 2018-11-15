package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

type userListQuery struct {
	httputil.TimeRangeQueryFreq
}

func (sv *Server) getUserList(c *gin.Context) {
	var (
		query userListQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	fromTime, toTime, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	userList, err := sv.storage.GetUserList(fromTime, toTime)
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}
	c.JSON(http.StatusOK, userList)
}
