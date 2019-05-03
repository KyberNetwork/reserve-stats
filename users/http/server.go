package http

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"              // import custom validator functions
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	trlib "github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger,
	rateProvider tokenrate.ETHUSDRateProvider,
	storage storage.Interface,
	host string,
	influxStorage *storage.InfluxStorage,
	userCapConf *common.UserCapConfiguration,
) *Server {
	r := gin.Default()
	return &Server{
		sugar:         sugar,
		rateProvider:  trlib.NewCachedRateProvider(sugar, rateProvider, trlib.WithTimeout(time.Hour)),
		storage:       storage,
		r:             r,
		host:          host,
		influxStorage: influxStorage,
		userCapConf:   userCapConf,
	}
}

//Server struct to represent a http server service
type Server struct {
	sugar         *zap.SugaredLogger
	r             *gin.Engine
	host          string
	rateProvider  tokenrate.ETHUSDRateProvider
	storage       storage.Interface
	influxStorage *storage.InfluxStorage
	userCapConf   *common.UserCapConfiguration
}

type userQuery struct {
	UserAddr  string `form:"address" binding:"isAddress"`
	TimeStamp uint64 `form:"time" binding:"required"`
}

func (s *Server) isKyced(c *gin.Context) {
	var query userQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}
	kyced, err := s.storage.IsKYCedAtTime(query.UserAddr, timeutil.TimestampMsToTime(query.TimeStamp))
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			fmt.Errorf("failed to check kyc status: %s", err.Error()),
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"kyced": kyced,
		},
	)
}

//createOrUpdate update info of an user
func (s *Server) createOrUpdate(c *gin.Context) {
	var userData common.UserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if err := s.storage.CreateOrUpdate(userData); err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": userData.Email})
}

type userStatsQuery struct {
	UID   string `form:"uid" binding:"required"`
	KYCed bool   `form:"kyced"`
}

// stats returns cap of the user with given uid.
func (s *Server) userStats(c *gin.Context) {
	var (
		logger = s.sugar.With("func", "users/http/Server.userStats")
		input  userStatsQuery

		userCap *big.Int
		rich    bool
	)
	if err := c.ShouldBindQuery(&input); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	logger = logger.With(
		"uid", input.UID,
		"kyced", input.KYCed,
	)

	logger.Debugw("querying stats for user")

	volume, err := s.influxStorage.Last24hVolume(input.UID)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}

	rate, err := s.rateProvider.USDRate(time.Now())
	if err != nil {
		logger.Errorw("failed to get usd rate", "err", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to get usd rate: %s", err.Error())},
		)
		return
	}

	userCap = blockchain.EthToWei(s.userCapConf.UserCap(input.KYCed).TxLimit / rate)
	rich = s.userCapConf.IsRich(input.KYCed, volume)

	logger.Infow("got last 24h volume of user",
		"volume", volume,
		"cap", userCap,
		"rich", rich,
	)

	c.JSON(http.StatusOK, gin.H{
		"cap":   userCap,
		"kyced": input.KYCed,
		"rich":  rich,
	})
}

func (s *Server) register() {
	s.r.POST("/users", s.createOrUpdate)
	s.r.GET("/kyced", s.isKyced)

	s.r.GET("/users", s.userStats)
}

//Run start server and serve
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
