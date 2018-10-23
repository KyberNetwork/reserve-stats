package http

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const limitedTimeRange = 24 * time.Hour

// Server serve trade logs through http endpoint
type Server struct {
	storage storage.Interface
	addr    string
	sugar   *zap.SugaredLogger
}

type tradeLogsQuery struct {
	From uint64 `form:"from"`
	To   uint64 `form:"to"`
}

type burnFeeQuery struct {
	From        uint64 `form:"from"`
	To          uint64 `form:"to"`
	Freq        string `form:"freq"`
	ReserveAddr string `form:"reserveAddr"`
}

func validateTimeWindow(fromTime, toTime time.Time, freq string) error {
	switch strings.ToLower(freq) {
	case "h":
		if toTime.After(fromTime.Add(time.Hour * 24 * 180)) {
			return errors.New("hourly frequency limit is 180 days")
		}
	case "d":
		if toTime.After(fromTime.Add(time.Hour * 24 * 365 * 3)) {
			return errors.New("daily frequency limit is 3 years")
		}
	default:
		return errors.New("invalid frequency")
	}
	return nil
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
		ha.sugar.Errorw(err.Error(), "fromTime", fromTime, "toTime", toTime)
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

func (ha *Server) getBurnFee(c *gin.Context) {
	var query burnFeeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	if !ethereum.IsHexAddress(query.ReserveAddr) {
		ha.sugar.Debug("invalid reserve address")
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid reserve address"},
		)
		return
	}

	fromTime := time.Unix(0, int64(query.From)*int64(time.Millisecond))
	toTime := time.Unix(0, int64(query.To)*int64(time.Millisecond))

	if err := validateTimeWindow(fromTime, toTime, query.Freq); err != nil {
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

	reserveAddr := ethereum.HexToAddress(query.ReserveAddr)

	burnFee, err := ha.storage.GetAggregatedBurnFee(fromTime, toTime, query.Freq, reserveAddr)
	if err != nil {
		ha.sugar.Errorw(err.Error(), "parameter", query)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		burnFee,
	)
}

func (ha *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/trade-logs", ha.getTradeLogs)
	r.GET("/burn-fee", ha.getBurnFee)
	return r
}

// Start running http server to serve trade logs data
func (ha *Server) Start() error {
	r := ha.setupRouter()
	return r.Run(ha.addr)
}

// NewServer returns an instance of HttpApi to serve trade logs
func NewServer(storage storage.Interface, addr string, sugar *zap.SugaredLogger) *Server {
	return &Server{storage: storage, addr: addr, sugar: sugar}
}
