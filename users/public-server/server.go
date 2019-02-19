package server

import (
	"fmt"
	"math/big"
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
	kycedCap     *common.UserCap
	nonKycedCap  *common.UserCap
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
		kycedCap:     common.NewUserCap(true),
		nonKycedCap:  common.NewUserCap(false),
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
		query       userQuery
		kyced, rich bool
		userCap     *big.Int
		err         error
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	logger.Info("query", "user query", query)

	kyced, err = s.getUserByKey(kycedPrefix, query.Address)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": err.Error()},
		)
		return
	}

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

	if kyced {
		userCap = blockchain.EthToWei(s.kycedCap.TxLimit / rate)
	} else {
		userCap = blockchain.EthToWei(s.nonKycedCap.TxLimit / rate)
	}

	c.JSON(
		http.StatusOK,
		common.UserResponse{
			Cap:   userCap,
			KYCed: kyced,
			Rich:  rich,
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
