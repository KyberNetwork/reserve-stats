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
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/tokenrate"
)

const (
	richPrefix  = "rich"
	kycedPrefix = "kyced"
)

//Server is server to serve api
type Server struct {
	sugar        *zap.SugaredLogger
	r            *gin.Engine
	host         string
	rateProvider tokenrate.ETHUSDRateProvider
	redisClient  *redis.Client
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
		rateProvider: httputil.NewCachedRateProvider(sugar, rateProvider, time.Hour),
		redisClient:  storage,
	}
}

func (s *Server) getUserFromRedis(userAddress string) (common.UserResponse, error) {
	var (
		user common.UserResponse
	)
	// try to get from rich users
	data := s.redisClient.Get(fmt.Sprintf("%s:%s", richPrefix, userAddress))
	if data.Err() == nil {
		userBytes, err := data.Bytes()
		if err != nil {
			return user, err
		}
		err = json.Unmarshal(userBytes, &user)
		return user, err
	}
	if data.Err() != redis.Nil {
		return user, data.Err()
	}

	// if not exist as rich user, try to get from kyced users
	data = s.redisClient.Get(fmt.Sprintf("%s:%s", kycedPrefix, userAddress))
	if data.Err() == nil {
		userBytes, err := data.Bytes()
		if err != nil {
			return user, err
		}
		err = json.Unmarshal(userBytes, &user)
		return user, err
	}
	if data.Err() != redis.Nil {
		return user, data.Err()
	}

	return user, nil
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

	userResponse, err := s.getUserFromRedis(query.Address)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed  to get usd rate: %s", err.Error())},
		)
		return
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
