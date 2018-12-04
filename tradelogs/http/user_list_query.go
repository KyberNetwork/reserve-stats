package http

import (
	"net/http"
	"sort"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/gin-gonic/gin"
)

type userListQuery struct {
	httputil.TimeRangeQuery
	Timezone int8 `form:"timezone" binding:"isSupportedTimezone"`
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
	fromTime, toTime, err := query.Validate(
		httputil.TimeRangeQueryWithMaxTimeFrame(maxTimeFrame),
		httputil.TimeRangeQueryWithDefaultTimeFrame(defaultTimeFrame),
	)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	userList, err := sv.storage.GetUserList(fromTime, toTime, query.Timezone)
	sort.Sort(sort.Reverse(common.UserList(userList)))
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
