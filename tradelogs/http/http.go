package http

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

// HTTPApi serve trade logs through http endpoint
type HTTPApi struct {
	storage storage.Interface
	router  *gin.Engine
	addr    string
}

func (ha *HTTPApi) getTradeLogs(c *gin.Context) {
	from, err := strconv.ParseInt(c.Query("from"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"success": false, "reason": fmt.Errorf("invalid from time %s", err.Error())},
		)
		return
	}
	to, err := strconv.ParseInt(c.Query("to"), 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"success": false, "reason": fmt.Errorf("invalid to time %s", err.Error())},
		)
		return
	}

	fromTime := time.Unix(from/1000, 0)
	toTime := time.Unix(to/1000, 0)
	tradeLogs, err := ha.storage.LoadTradeLogs(fromTime, toTime)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{"success": false, "reason": fmt.Errorf("invalid to time %s", err.Error())},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"success": true, "data": tradeLogs},
	)
}

// Start running http server to serve trade logs data
func (ha *HTTPApi) Start() {
	ha.router.GET("/trade-logs", ha.getTradeLogs)
	ha.router.Run(ha.addr)
}

// NewHTTPApi returns an instance of HttpApi to serve trade logs
func NewHTTPApi(storage storage.Interface, addr string) *HTTPApi {
	r := gin.Default()
	return &HTTPApi{storage: storage, router: r, addr: addr}
}
