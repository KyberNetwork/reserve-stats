package server

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	trlib "github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/tokenrate"
)

const (
	richPrefix = "rich"
)

//Server is server to serve api
type Server struct {
	sugar        *zap.SugaredLogger
	r            *gin.Engine
	host         string
	rateProvider tokenrate.ETHUSDRateProvider
	redisClient  *redis.Client
	userCapConf  *common.UserCapConfiguration
}

//UserQuery is query for user info
type userQuery struct {
	//Address is user address
	Address string `form:"address" binding:"required,isAddress"`
}

//NewServer return new server instance
func NewServer(
	logger *zap.Logger,
	host string,
	rateProvider tokenrate.ETHUSDRateProvider,
	storage *redis.Client,
	userCapConf *common.UserCapConfiguration) *Server {
	r := gin.Default()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	sugar := logger.Sugar()
	return &Server{
		sugar:        sugar,
		r:            r,
		host:         host,
		rateProvider: trlib.NewCachedRateProvider(sugar, rateProvider, trlib.WithTimeout(time.Hour)),
		redisClient:  storage,
		userCapConf:  userCapConf,
	}
}

func (s *Server) getUserByKey(prefix, userAddress string) (bool, error) {
	data := s.redisClient.Get(fmt.Sprintf("%s:%s", prefix, userAddress))
	if data.Err() != nil {
		if data.Err() == redis.Nil {
			return false, nil
		}
		return false, data.Err()
	}
	return true, nil
}

func (s *Server) getUsers(c *gin.Context) {
	var (
		logger = s.sugar.With(
			"func", "users-public-stats/Server.getUser",
		)
		query   userQuery
		rich    bool
		userCap *big.Int
		err     error
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	logger.Info("query", "user query", query)

	rich, err = s.getUserByKey(richPrefix, query.Address)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
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

	userCap = blockchain.EthToWei(s.userCapConf.UserCap(false).TxLimit / rate)

	c.JSON(
		http.StatusOK,
		common.UserResponse{
			Cap:  userCap,
			Rich: rich,
		},
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
