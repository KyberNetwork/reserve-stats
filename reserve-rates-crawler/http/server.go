package http

import (
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/common"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// Server is the engine to serve reserve-rate API query
type Server struct {
	r    *gin.Engine
	db   rateStorage
	host string
}

func (sv *Server) getReserveRate(c *gin.Context) {
	fromTime, _ := strconv.ParseUint(c.Query("fromTime"), 10, 64)
	toTime, _ := strconv.ParseUint(c.Query("toTime"), 10, 64)
	if toTime == 0 {
		toTime = common.TimeToTimestampMs(time.Now())
	}
	reserveAddr := ethereum.HexToAddress(c.Query("reserveAddr"))
	if reserveAddr.Big().Cmp(ethereum.Big0) == 0 {
		httputil.ResponseFailure(c, httputil.WithReason("Reserve address is invalid"))
		return
	}
	result, err := sv.db.GetRatesByTimePoint(reserveAddr, fromTime, toTime)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

// Run starts HTTP server on preconfigure-host. Return error if occurs
func (sv *Server) Run() error {
	sv.r.GET("/reserve-rate", sv.getReserveRate)
	return sv.r.Run(sv.host)
}

// NewServer create an instance of Server to serve API query
func NewServer(host string, db rateStorage) (*Server, error) {
	r := gin.Default()
	return &Server{
		r:    r,
		db:   db,
		host: host,
	}, nil
}
