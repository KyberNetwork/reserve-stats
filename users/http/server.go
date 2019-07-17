package http

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"              // import custom validator functions
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	trlib "github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/users/common"
)

const (
	uidPrefix = "uid"
)

//NewServer return new server instance
func NewServer(sugar *zap.SugaredLogger,
	rateProvider tokenrate.ETHUSDRateProvider,
	host string,
	redisClient *redis.Client,
	userCapConf *common.UserCapConfiguration,
) *Server {
	r := gin.Default()
	return &Server{
		sugar:        sugar,
		rateProvider: trlib.NewCachedRateProvider(sugar, rateProvider, trlib.WithTimeout(time.Hour)),
		r:            r,
		host:         host,
		redisClient:  redisClient,
		userCapConf:  userCapConf,
	}
}

//Server struct to represent a http server service
type Server struct {
	sugar        *zap.SugaredLogger
	r            *gin.Engine
	host         string
	rateProvider tokenrate.ETHUSDRateProvider
	redisClient  *redis.Client
	userCapConf  *common.UserCapConfiguration
}

type userStatsQuery struct {
	UID   string `form:"uid" binding:"required"`
	KYCed bool   `form:"kyced"`
}

func (s *Server) getUserVolumeByUID(uid string) (float64, error) {
	data := s.redisClient.Get(fmt.Sprintf("%s:%s", uidPrefix, uid))
	if data.Err() != nil {
		if data.Err() == redis.Nil {
			return 0, nil
		}
		return 0, data.Err()
	}
	return data.Float64()
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

	volume, err := s.getUserVolumeByUID(input.UID)
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
	volumeInWei := blockchain.EthToWei(volume / rate)
	rich = s.userCapConf.IsRich(input.KYCed, volume)

	logger.Infow("got last 24h volume of user",
		"volume", volumeInWei,
		"cap", userCap,
		"rich", rich,
	)

	c.JSON(http.StatusOK, gin.H{
		"cap":    userCap,
		"kyced":  input.KYCed,
		"rich":   rich,
		"volume": volumeInWei,
	})
}

func (s *Server) register() {
	s.r.GET("/users", s.userStats)
}

//Run start server and serve
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
