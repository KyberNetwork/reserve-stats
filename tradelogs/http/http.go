package http

import (
	"net/http"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appname "github.com/KyberNetwork/reserve-stats/app-names"
	lipappnames "github.com/KyberNetwork/reserve-stats/lib/appnames"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	libhttputil "github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/userprofile"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const (
	hourlyBurnFeeMaxDuration = time.Hour * 24 * 180 // 180 days
)

// NewServer returns an instance of HttpApi to serve trade logs.
func NewServer(
	storage storage.Interface,
	host string,
	sugar *zap.SugaredLogger,
	symbolResolver blockchain.TokenSymbolResolver, options ...ServerOption) *Server {
	var (
		logger = sugar.With("func", "tradelogs/http/NewServer")
		sv     = &Server{
			storage:        storage,
			host:           host,
			sugar:          sugar,
			symbolResolver: symbolResolver,
		}
	)

	for _, opt := range options {
		opt(sv)
	}

	if sv.getAddrToAppName == nil {
		logger.Warn("application names integration is not configured")
		sv.getAddrToAppName = func() (map[ethereum.Address]string, error) { return nil, nil }
	}

	if sv.getUserProfile == nil {
		logger.Warn("user profile integration is not configured")
		sv.getUserProfile = func(ethereum.Address) (userprofile.UserProfile, error) { return userprofile.UserProfile{}, nil }
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

// WithUserProfile configures the Server instance to use user profile lookup
func WithUserProfile(up userprofile.Interface) ServerOption {
	return func(sv *Server) {
		sv.getUserProfile = up.LookUpUserProfile
	}
}

// Server serve trade logs through http endpoint.
type Server struct {
	storage          storage.Interface
	host             string
	sugar            *zap.SugaredLogger
	getAddrToAppName func() (map[ethereum.Address]string, error)
	getUserProfile   func(ethereum.Address) (userprofile.UserProfile, error)
	symbolResolver   blockchain.TokenSymbolResolver
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
		// get user profile
		up, err := sv.getUserProfile(tradeLogs[i].UserAddress)
		if err != nil {
			sv.sugar.Errorw(err.Error(), "fromTime", fromTime, "toTime", toTime)
			libhttputil.ResponseFailure(
				c,
				http.StatusInternalServerError,
				err,
			)
			return
		}
		tradeLogs[i].UserName = up.UserName
		tradeLogs[i].ProfileID = up.ProfileID
		if (tradeLogs[i].IntegrationApp != appname.KyberSwapAppName) && (len(log.WalletFees) > 0) {
			name, avai := addrToAppName[log.WalletFees[0].WalletAddress]
			if avai {
				tradeLogs[i].IntegrationApp = name
			}
		}

		// resolve token symbol
		if !blockchain.IsZeroAddress(log.SrcAddress) {
			srcSymbol, err := sv.symbolResolver.Symbol(log.SrcAddress)
			if err != nil {
				libhttputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
			tradeLogs[i].SrcSymbol = srcSymbol
		}

		if !blockchain.IsZeroAddress(log.DestAddress) {
			dstSymbol, err := sv.symbolResolver.Symbol(log.DestAddress)
			if err != nil {
				libhttputil.ResponseFailure(
					c,
					http.StatusInternalServerError,
					err,
				)
				return
			}
			tradeLogs[i].DestSymbol = dstSymbol
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
	logger := sv.sugar.Desugar()
	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
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
