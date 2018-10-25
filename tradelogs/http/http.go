package http

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	libhttputil "github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const (
	hourlyBurnFeeMaxDuration = time.Hour * 24 * 180 // 180 days
)

// Server serve trade logs through http endpoint
type Server struct {
	storage     storage.Interface
	host        string
	sugar       *zap.SugaredLogger
	coreSetting core.Interface
}

type burnFeeQuery struct {
	libhttputil.TimeRangeQueryFreq
	ReserveAddrs []string `form:"reserve" binding:"dive,isAddress"`
}

func (sv *Server) getTradeLogs(c *gin.Context) {
	var query libhttputil.TimeRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		libhttputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	fromTime, toTime, err := query.Validate()
	if err != nil {
		libhttputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	tradeLogs, err := sv.storage.LoadTradeLogs(fromTime, toTime)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "fromTime", fromTime, "toTime", toTime)
		libhttputil.ResponseFailure(
			c,
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
		libhttputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	fromTime, toTime, err := query.Validate()
	if err != nil {
		libhttputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	for _, rsvAddr := range query.ReserveAddrs {
		rsvAddrs = append(rsvAddrs, ethereum.HexToAddress(rsvAddr))
	}

	burnFee, err := sv.storage.GetAggregatedBurnFee(fromTime, toTime, query.Freq, rsvAddrs)
	if err != nil {
		sv.sugar.Errorw(err.Error(), "parameter", query)
		libhttputil.ResponseFailure(
			c,
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
	r.GET("/trade-logs", sv.getTradeLogs)
	r.GET("/burn-fee", sv.getBurnFee)
	r.GET("/asset-volume", sv.getAssetVolume)
	r.GET("/reserve-volume", sv.getReserveVolume)
	r.GET("/wallet-fee", sv.getWalletFee)
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
