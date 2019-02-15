package server

import (
	"fmt"
	"net/http"
	"strconv"
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

func fromInterfaceToUser(data []interface{}) (common.UserResponse, error) {
	if len(data) != 2 {
		return common.UserResponse{}, fmt.Errorf("data len should be 2: %d", len(data))
	}
	kyced, ok := data[0].(string)
	if !ok {
		return common.UserResponse{}, fmt.Errorf("kyced value shoud be string: %v", data[0])
	}
	kycInt, err := strconv.ParseInt(kyced, 10, 64)
	if err != nil {
		return common.UserResponse{}, err
	}
	rich, ok := data[1].(string)
	if !ok {
		return common.UserResponse{}, fmt.Errorf("rich should be string: %v", data[1])
	}
	richInt, err := strconv.ParseInt(rich, 10, 64)
	if err != nil {
		return common.UserResponse{}, err
	}
	return common.UserResponse{
		KYCed: kycInt == 1,
		Rich:  richInt == 1,
	}, nil
}

func (s *Server) getUserFromRedis(userAddress string) (common.UserResponse, error) {
	// try to get from rich users
	data, err := s.redisClient.HMGet(fmt.Sprintf("%s:%s", richPrefix, userAddress), kycedPrefix, richPrefix).Result()
	if err == nil && data[0] != nil {
		return fromInterfaceToUser(data)
	}

	// if not exist as rich user, try to get from kyced users
	data, err = s.redisClient.HMGet(fmt.Sprintf("%s:%s", kycedPrefix, userAddress), kycedPrefix, richPrefix).Result()
	if err == nil && data[0] != nil {
		return fromInterfaceToUser(data)
	}

	return common.UserResponse{}, nil
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
			gin.H{"error": fmt.Sprintf("failed  to get user: %s", err.Error())},
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
