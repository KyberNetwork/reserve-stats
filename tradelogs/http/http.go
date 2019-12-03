package http

import (
	"errors"
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
	symbol, err := sv.symbolResolver.Symbol(tokenAddress)
	if err != nil {
		return "", err
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

type tradeLogsByTxHashParam struct {
	TxHash string `uri:"tx_hash" binding:"required"`
}

func (sv *Server) getTradeLogsByTx(c *gin.Context) {
	var param tradeLogsByTxHashParam
	if err := c.ShouldBindUri(&param); err != nil {
		libhttputil.ResponseFailure(
			c, http.StatusInternalServerError, err,
		)
		return
	}
	if !blockchain.IsValidTxHash(param.TxHash) {
		libhttputil.ResponseFailure(
			c, http.StatusBadRequest, errors.New("invalid transaction hash"),
		)
		return
	}
	tradeLogs, err := sv.storage.LoadTradeLogsByTxHash(ethereum.HexToHash(param.TxHash))
	if err != nil {
		libhttputil.ResponseFailure(
			c, http.StatusInternalServerError, err,
		)
		return
	}

	c.JSON(
		http.StatusOK,
		tradeLogs,
	)
}

func (sv *Server) getStats(c *gin.Context) {
	/** --> query
	 **/
	var query libhttputil.TimeRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	from, to, err := query.Validate()
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	stats, err := sv.storage.GetStats(from, to)
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(
		http.StatusOK,
		stats,
	)
}

func (sv *Server) getTopTokens(c *gin.Context) {
	var query libhttputil.TimeRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	from, to, err := query.Validate()
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	topTokens, err := sv.storage.GetTopTokens(from, to)
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(
		http.StatusOK,
		topTokens,
	)
}

func (sv *Server) getTopIntegration(c *gin.Context) {
	var query libhttputil.TimeRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	from, to, err := query.Validate()
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	topIntegration, err := sv.storage.GetTopIntegrations(from, to)
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(
		http.StatusOK,
		topIntegration,
	)
}

func (sv *Server) getTopReserves(c *gin.Context) {
	var query libhttputil.TimeRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	from, to, err := query.Validate()
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	topReserves, err := sv.storage.GetTopReserves(from, to)
	if err != nil {
		libhttputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(
		http.StatusOK,
		topReserves,
	)
}

func (sv *Server) setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/trade-logs", sv.getTradeLogs)
	r.GET("/trade-logs/:tx_hash", sv.getTradeLogsByTx)
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

	// twitter api
	r.GET("/stats", sv.getStats)
	r.GET("/top-tokens", sv.getTopTokens)
	r.GET("/top-integrations", sv.getTopIntegration)
	r.GET("/top-reserves", sv.getTopReserves)

	return r
}

// Start running http server to serve trade logs data
func (sv *Server) Start() error {
	r := sv.setupRouter()
	return r.Run(sv.host)
}
