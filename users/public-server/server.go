package server

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
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
	blacklist    *common.Blacklist
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
	userCapConf *common.UserCapConfiguration,
	blacklist *common.Blacklist) *Server {
	r := gin.Default()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	sugar := logger.Sugar()
	return &Server{
		sugar:        sugar,
		r:            r,
		host:         host,
		rateProvider: trlib.NewCachedRateProvider(sugar, rateProvider, trlib.WithExpires(time.Hour)),
		redisClient:  storage,
		userCapConf:  userCapConf,
		blacklist:    blacklist,
	}
}

func (s *Server) getAddressVolumeByKey(prefix, userAddress string) (float64, error) {
	data := s.redisClient.Get(fmt.Sprintf("%s:%s", prefix, userAddress))
	if data.Err() != nil {
		if data.Err() == redis.Nil {
			return 0, nil
		}
		return 0, data.Err()
	}
	return data.Float64()
}

func (s *Server) getUsers(c *gin.Context) {
	var (
		logger  = s.sugar.With("func", caller.GetCurrentFunctionName())
		query   userQuery
		rich    bool
		userCap *big.Int
		err     error
		volume  float64
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	logger.Info("query", "user query", query)
	//return 0 for banned address
	if s.blacklist.IsBanned(ethereum.HexToAddress(query.Address)) {
		c.JSON(http.StatusOK, common.UserResponse{
			Cap:  big.NewInt(0),
			Rich: false,
		})
		return
	}
	volume, err = s.getAddressVolumeByKey(richPrefix, query.Address)
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, err)
		return
	}

	rate, err := s.rateProvider.USDRate(time.Now())
	if err != nil {
		httputil.ResponseFailure(c, http.StatusInternalServerError, errors.Wrap(err, "failed  to get usd rate"))
		return
	}

	userCap = blockchain.EthToWei(s.userCapConf.UserCap(false).TxLimit / rate)
	rich = s.userCapConf.IsRich(false, volume)

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
