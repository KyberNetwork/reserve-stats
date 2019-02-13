package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/users/common"
	userhttp "github.com/KyberNetwork/reserve-stats/users/http"
	"github.com/KyberNetwork/tokenrate"
)

//Server is server to serve api
type Server struct {
	sugar        *zap.SugaredLogger
	r            *gin.Engine
	host         string
	rateProvider tokenrate.ETHUSDRateProvider
	storage      *redis.Client
}

//UserQuery is query for user info
type userQuery struct {
	//Address is user address
	Address string `form:"address" binding:"required,isAddress"`
}

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger, host string, rateProvider tokenrate.ETHUSDRateProvider, storage *redis.Client) *Server {
	r := gin.Default()
	return &Server{
		sugar:        sugar,
		r:            r,
		host:         host,
		rateProvider: userhttp.NewCachedRateProvider(sugar, rateProvider, time.Hour),
		storage:      storage,
	}
}

func (s *Server) getUsers(c *gin.Context) {
	var (
		logger = s.sugar.With(
			"func", "users-public-stats/Server.getUser",
		)
		query        userQuery
		userResponse common.UserResponse
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
	}

	logger.Info("query", "user query", query)

	user := s.storage.Get(query.Address)
	if user.Err() == nil {
		userBytes, err := user.Bytes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = json.Unmarshal(userBytes, &userResponse)
		logger.Debugw("user get from redis", "user", userResponse)
	}

	rate, err := s.rateProvider.USDRate(time.Now())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed  to get usd rate: %s", err.Error())},
		)
		return
	}
	if userResponse.KYCed {
		kycedCap := common.NewUserCap(true)
		userResponse.Cap = blockchain.EthToWei(kycedCap.TxLimit / rate)
	} else {
		nonKycedCap := common.NewUserCap(false)
		userResponse.Cap = blockchain.EthToWei(nonKycedCap.TxLimit / rate)
	}

	c.JSON(
		http.StatusOK,
		userResponse,
	)
}

func (s *Server) register() {
	s.r.GET("/users", s.getUsers)
}

//Run start the server
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
