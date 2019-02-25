package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"              // import custom validator functions
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger, rateProvider tokenrate.ETHUSDRateProvider,
	storage storage.Interface, host string,
	influxStorage *storage.InfluxStorage) *Server {
	r := gin.Default()
	return &Server{
		sugar:         sugar,
		rateProvider:  httputil.NewCachedRateProvider(sugar, rateProvider, time.Hour),
		storage:       storage,
		r:             r,
		host:          host,
		influxStorage: influxStorage,
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

//getTransactionLimit returns cap limit of a user.
func (s *Server) getTransactionLimit(c *gin.Context) {
	var logger = s.sugar.With(
		"func", "users/http/Server.getTransactionLimit",
	)

	address := c.Query("address")
	if !ethereum.IsHexAddress(address) {
		httputil.ResponseFailure(
			c,
			http.StatusBadRequest,
			fmt.Errorf("provided address is not valid: %s", address),
		)
		return
	}

	logger = logger.With("address", address)

	kyced, err := s.storage.IsKYCedAtTime(address, time.Now())
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			fmt.Errorf("failed to check kyc status: %s", err.Error()),
		)
		return
	}

	rate, err := s.rateProvider.USDRate(time.Now())
	if err != nil {
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			fmt.Errorf("failed to get usd rate: %s", err.Error()),
		)
		return
	}

	// maximum of ETH in wei
	uc := common.NewUserCap(kyced)
	txLimit := blockchain.EthToWei(uc.TxLimit / rate)
	rich, err := s.influxStorage.IsExceedDailyLimit(address, uc.DailyLimit)
	if err != nil {
		var errMsg = "could not retrieve user volume"
		logger.Errorw(errMsg, "err", err.Error())
		httputil.ResponseFailure(
			c,
			http.StatusInternalServerError,
			errors.New(errMsg),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"cap":   txLimit,
			"rich":  rich,
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

func (s *Server) register() {
	s.r.GET("/users", s.getTransactionLimit)
	s.r.GET("/kyced", s.isKyced)
	s.r.POST("/users", s.createOrUpdate)
}

//Run start server and serve
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
