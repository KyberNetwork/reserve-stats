package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/KyberNetwork/tokenrate"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger, rateProvider tokenrate.ETHUSDRateProvider, 
	storage storage.Interface, host string,
	influxStorage *storage.InfluxStorage) *Server {
	r := gin.Default()
	return &Server{
		sugar:        sugar,
		rateProvider: newCachedRateProvider(sugar, rateProvider, time.Hour),
		storage:      storage,
		r:            r,
		host:         host,
		influxStorage: influxStorage,
	}
}

//Server struct to represent a http server service
type Server struct {
	sugar        *zap.SugaredLogger
	r            *gin.Engine
	host         string
	rateProvider tokenrate.ETHUSDRateProvider
	storage      storage.Interface
	influxStorage       *storage.InfluxStorage
}

//getTransactionLimit returns cap limit of a user.
func (s *Server) getTransactionLimit(c *gin.Context) {
	address := c.Query("address")
	kyced, err := s.storage.IsKYCed(address)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("failed to check KYC status: %s", err.Error()),
			},
		)
		return
	}

	rate, err := s.rateProvider.USDRate(time.Now())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("failed to get usd rate: %s", err.Error()),
			},
		)
		return
	}

	// maximum of ETH in wei
	uc := common.NewUserCap(kyced)
	txLimit := blockchain.EthToWei(uc.TxLimit / rate)
	rich, err := s.influxStorage.IsExceedDailyLimit(address)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)	
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"cap": txLimit,
			"rich":  rich,
			"kyced": kyced,
		},
	)
}

//createOrUpdate update info of an user
func (s *Server) createOrUpdate(c *gin.Context) {
	var userData common.UserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	if err := s.storage.CreateOrUpdate(userData); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"email": userData.Email})
}

func (s *Server) register() {
	s.r.GET("/users", s.getTransactionLimit)
	s.r.POST("/users", s.createOrUpdate)
}

//Run start server and serve
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
