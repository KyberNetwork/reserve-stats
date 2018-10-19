package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const limitedTimeRange = 24 * time.Hour

// Server serve trade logs through http endpoint
type Server struct {
	storage storage.Interface
	router  *gin.Engine
	addr    string
}

type tradeLogsQuery struct {
	From uint64 `form:"from"`
	To   uint64 `form:"to"`
}

func (ha *Server) getTradeLogs(c *gin.Context) {
	var query tradeLogsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	fromTime := time.Unix(0, int64(query.From)*int64(time.Millisecond))
	toTime := time.Unix(0, int64(query.To)*int64(time.Millisecond))

	if toTime.After(fromTime.Add(limitedTimeRange)) {
		err := fmt.Errorf("time range is too broad, must be smaller or equal to %d milliseconds", limitedTimeRange/time.Millisecond)
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if toTime.Equal(time.Unix(0, 0)) {
		toTime = time.Now()
		fromTime = toTime.Add(-time.Hour)
	}

	tradeLogs, err := ha.storage.LoadTradeLogs(fromTime, toTime)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		tradeLogs,
	)
}

// Start running http server to serve trade logs data
func (ha *Server) Start() error {
	ha.router.GET("/trade-logs", ha.getTradeLogs)

	return ha.router.Run(ha.addr)
}

// NewServer returns an instance of HttpApi to serve trade logs
func NewServer(storage storage.Interface, addr string) *Server {
	r := gin.Default()
	return &Server{storage: storage, router: r, addr: addr}
}
