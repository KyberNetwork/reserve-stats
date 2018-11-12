package http

import (
	"net/http"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/gin-gonic/gin"
)

type userListQuery struct {
	httputil.TimeRangeQuery
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
	_, _, err := query.Validate()
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	userList, err := sv.storage.GetUserList(query.From, query.To)
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
