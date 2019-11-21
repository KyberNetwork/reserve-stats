package http

import (
	"net/http"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	appname "github.com/KyberNetwork/reserve-stats/app-names"
	lipappnames "github.com/KyberNetwork/reserve-stats/lib/appnames"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
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
		logger = sugar.With("func", caller.GetCurrentFunctionName())
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

func (sv *Server) getTokenSymbol(tokenAddress ethereum.Address) (string, error) {
	symbol, err := sv.storage.GetTokenSymbol(tokenAddress.Hex())
	if err != nil {
		// we only log error here to notify about failed interact with db
		// we can get symbol from blockchain later
		sv.sugar.Errorw("get token symbol from database failed", "error", err)
	}
	if symbol == "" { // if cannot get symbol from database, then we get it from blockchain and save to db
		symbol, err = sv.symbolResolver.Symbol(tokenAddress)
		if err != nil {
			return "", err
		}
		// save srcSymbol token to db
		addresses := []string{tokenAddress.Hex()}
		symbols := []string{symbol}
		if err := sv.storage.UpdateTokens(addresses, symbols); err != nil {
			// we only log here as we already get token symbol from blockchain
			// we log error then sentry can trigger to notify us, the http flow should not
			// be broken
			sv.sugar.Errorw("cannot update token symbol", "error", err)
		}
	}
	return symbol, nil
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
			srcSymbol, err := sv.getTokenSymbol(log.SrcAddress)
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
			dstSymbol, err := sv.getTokenSymbol(log.DestAddress)
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

type getSymbolRequest struct {
	Address string `form:"address" binding:"required"`
}

func (sv *Server) getSymbol(c *gin.Context) {
	var query getSymbolRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		libhttputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	symbol, err := sv.storage.GetTokenSymbol(query.Address)
	if err != nil {
		libhttputil.ResponseFailure(
			c, http.StatusInternalServerError, err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"symbol": symbol,
		},
	)
}

type tokenDetail struct {
	Address string `json:"address" binding:"required"`
	Symbol  string `json:"symbol" binding:"required"`
}

type updateSymbolRequest []tokenDetail

func (sv *Server) updateSymbol(c *gin.Context) {
	var (
		addresses, symbol []string
		query             updateSymbolRequest
	)
	if err := c.ShouldBindJSON(&query); err != nil {
		libhttputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}
	for _, token := range query {
		addresses = append(addresses, token.Address)
		symbol = append(symbol, token.Symbol)
	}

	if err := sv.storage.UpdateTokens(addresses, symbol); err != nil {
		libhttputil.ResponseFailure(
			c, http.StatusInternalServerError, err,
		)
		return
	}
	c.JSON(
		http.StatusOK,
		nil,
	)
}

func (sv *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/trade-logs", sv.getTradeLogs)
	r.GET("/burn-fee", sv.getBurnFee)
	r.GET("/asset-volume", sv.getAssetVolume)
	r.GET("/monthly-volume", sv.getMonthlyVolume)
	r.GET("/reserve-volume", sv.getReserveVolume)
	r.GET("/wallet-fee", sv.getWalletFee)
	r.GET("/trade-summary", sv.getTradeSummary)
	r.GET("/user-volume", sv.getUserVolume)
	r.GET("/user-list", sv.getUserList)
	r.GET("/wallet-stats", sv.getWalletStats)
	r.GET("/country-stats", sv.getCountryStats)
	r.GET("/heat-map", sv.getTokenHeatMap)
	r.GET("/integration-volume", sv.getIntegrationVolume)

	// token symbol
	r.GET("/symbol", sv.getSymbol)
	r.POST("/symbol", sv.updateSymbol)

	return r
}

// Start running http server to serve trade logs data
func (sv *Server) Start() error {
	r := sv.setupRouter()
	return r.Run(sv.host)
}
