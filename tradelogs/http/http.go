package http

import (
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/app-names"
	lipappnames "github.com/KyberNetwork/reserve-stats/lib/appnames"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	libhttputil "github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	hourlyBurnFeeMaxDuration = time.Hour * 24 * 180 // 180 days
)

// NewServer returns an instance of HttpApi to serve trade logs.
func NewServer(storage storage.Interface, host string, sugar *zap.SugaredLogger, sett core.Interface, options ...ServerOption) *Server {
	var (
		logger = sugar.With("func", "tradelogs/http/NewServer")
		sv     = &Server{
			storage:     storage,
			host:        host,
			sugar:       sugar,
			coreSetting: sett,
		}
	)

	for _, opt := range options {
		opt(sv)
	}

	if sv.getAddrToAppName == nil {
		logger.Warn("application names integration is not configured")
		sv.getAddrToAppName = func() (map[ethereum.Address]string, error) { return nil, nil }
	}

	return sv
}

// ServerOption configures the behaviour of Server constructor.
type ServerOption func(server *Server)

// WithApplicationNames configures the Server instance to use appname integration.
func WithApplicationNames(an lipappnames.AddrToAppName) ServerOption {
	return func(sv *Server) {
		sv.getAddrToAppName = an.GetAddrToAppName
	}
}

// Server serve trade logs through http endpoint.
type Server struct {
	storage          storage.Interface
	host             string
	sugar            *zap.SugaredLogger
	coreSetting      core.Interface
	getAddrToAppName func() (map[ethereum.Address]string, error)
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
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
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

	addrToAppName, err := sv.getAddrToAppName()
	if err != nil {
		libhttputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	for i, log := range tradeLogs {
		if (log.IntegrationApp != appname.KyberSwapAppName) && (len(log.WalletFees) > 0) {
			name, avai := addrToAppName[log.WalletFees[0].WalletAddress]
			if avai {
				tradeLogs[i].IntegrationApp = name
			}
		}
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
	r.GET("/trade-summary", sv.getTradeSummary)
	r.GET("/user-volume", sv.getUserVolume)
	r.GET("/user-list", sv.getUserList)
	r.GET("/wallet-stats", sv.getWalletStats)
	r.GET("/country-stats", sv.getCountryStats)
	r.GET("/heat-map", sv.getTokenHeatMap)
	r.GET("/integration-volume", sv.getIntegrationVolume)
	return r
}

// Start running http server to serve trade logs data
func (sv *Server) Start() error {
	r := sv.setupRouter()
	return r.Run(sv.host)
}
