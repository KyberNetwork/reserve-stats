package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	libhttputil "github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

const (
	limitedTimeRange         = 24 * time.Hour
	hourlyBurnFeeMaxDuration = time.Hour * 24 * 180 // 180 days
	hourlyFreq               = "h"
	dailyBurnFeeMaxDuration  = time.Hour * 24 * 365 * 3 // ~ 3 years
	dailyFreq                = "d"
)

// Server serve trade logs through http endpoint
type Server struct {
	storage     storage.Interface
	host        string
	sugar       *zap.SugaredLogger
	coreSetting core.Interface
}

type tradeLogsQuery struct {
	From uint64 `form:"from"`
	To   uint64 `form:"to"`
}

type burnFeeQuery struct {
	From         uint64   `form:"from"`
	To           uint64   `form:"to"`
	Freq         string   `form:"freq"`
	ReserveAddrs []string `form:"reserve" binding:"dive,isAddress"`
}

func validateTimeWindow(fromTime, toTime time.Time, freq string) error {
	switch strings.ToLower(freq) {
	case hourlyFreq:
		if toTime.After(fromTime.Add(hourlyBurnFeeMaxDuration)) {
			return fmt.Errorf("your query time range exceeds the duration limit %s", hourlyBurnFeeMaxDuration)
		}
	case dailyFreq:
		if toTime.After(fromTime.Add(dailyBurnFeeMaxDuration)) {
			return fmt.Errorf("your query time range exceeds the duration limit %s", dailyBurnFeeMaxDuration)
		}
	default:
		return fmt.Errorf("your query frequency is not supported, use %s or %s", hourlyFreq, dailyFreq)
	}
	return nil
}

func (sv *Server) getTradeLogs(c *gin.Context) {
	var query tradeLogsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			err,
		)
		return
	}

	fromTime := time.Unix(0, int64(query.From)*int64(time.Millisecond))
	toTime := time.Unix(0, int64(query.To)*int64(time.Millisecond))

	if toTime.After(fromTime.Add(limitedTimeRange)) {
		err := fmt.Errorf("time range is too broad, must be smaller or equal to %d milliseconds", limitedTimeRange/time.Millisecond)
		c.AbortWithError(
			http.StatusBadRequest,
			err,
		)
		return
	}

	if toTime.Equal(time.Unix(0, 0)) {
		toTime = time.Now()
		fromTime = toTime.Add(-time.Hour)
	}

	tradeLogs, err := sv.storage.LoadTradeLogs(fromTime, toTime)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "fromTime", fromTime, "toTime", toTime)
		c.AbortWithError(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		tradeLogs,
	)
}

func (sv *Server) getBurnFee(c *gin.Context) {
	var (
		query    burnFeeQuery
		rsvAddrs []ethereum.Address
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			err,
		)
		return
	}

	fromTime := timeutil.TimestampMsToTime(query.From)
	toTime := timeutil.TimestampMsToTime(query.To)

	if err := validateTimeWindow(fromTime, toTime, query.Freq); err != nil {
		c.AbortWithError(
			http.StatusBadRequest,
			err,
		)
		return
	}

	if toTime.IsZero() {
		toTime = time.Now()
	}

	if fromTime.IsZero() {
		fromTime = toTime.Add(-time.Hour)
	}

	for _, rsvAddr := range query.ReserveAddrs {
		rsvAddrs = append(rsvAddrs, ethereum.HexToAddress(rsvAddr))
	}

	burnFee, err := sv.storage.GetAggregatedBurnFee(fromTime, toTime, query.Freq, rsvAddrs)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "parameter", query)
		c.AbortWithError(
			http.StatusInternalServerError,
			err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		burnFee,
	)
}

func (sv *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(libhttputil.MiddlewareHandler)
	r.GET("/trade-logs", sv.getTradeLogs)
	r.GET("/burn-fee", sv.getBurnFee)
	r.GET("/asset-volume", sv.getAssetVolume)
	return r
}

// Start running http server to serve trade logs data
func (sv *Server) Start() error {
	r := sv.setupRouter()
	return r.Run(sv.host)
}

// NewServer returns an instance of HttpApi to serve trade logs
func NewServer(storage storage.Interface, host string, sugar *zap.SugaredLogger, sett core.Interface) *Server {
	return &Server{storage: storage, host: host, sugar: sugar, coreSetting: sett}
}
