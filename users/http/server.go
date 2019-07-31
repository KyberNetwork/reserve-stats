package http

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
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
	maxBatchSize int,
) *Server {
	r := gin.Default()
	return &Server{
		sugar:        sugar,
		rateProvider: trlib.NewCachedRateProvider(sugar, rateProvider, trlib.WithTimeout(time.Hour)),
		r:            r,
		host:         host,
		redisClient:  redisClient,
		userCapConf:  userCapConf,
		maxBatchSize: maxBatchSize,
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
	maxBatchSize int
}

type userStatsQuery struct {
	UID   string `form:"uid" binding:"required"`
	KYCed bool   `form:"kyced"`
}

type userStatsBatchQuery struct {
	UIDs  string `form:"uids" binding:"required"`
	KYCed string `form:"kyced"  binding:"required"`
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

func (s *Server) getUserVolumeByUIDs(uids []string) ([]float64, error) {
	var err error
	pipeline := s.redisClient.Pipeline()
	var cmds []*redis.StringCmd
	for _, uid := range uids {
		cmds = append(cmds, pipeline.Get(fmt.Sprintf("%s:%s", uidPrefix, uid)))
	}
	if _, err := pipeline.Exec(); err != nil {
		if err != redis.Nil {
			return nil, errors.Wrap(err, "failed to exec pipeline")
		}
	}

	var result []float64
	for _, cmd := range cmds {
		var volume float64
		switch cmd.Err() {
		case nil:
			if volume, err = cmd.Float64(); err != nil {
				return nil, errors.Wrap(err, "failed to convert result to float64")
			}
		case redis.Nil:
			volume = 0
		default:
			return nil, errors.Wrap(cmd.Err(), "failed to exec singer cmd")
		}
		result = append(result, volume)
	}
	return result, nil
}

func (s *Server) convertQueryParams(query userStatsBatchQuery) ([]string, []bool, error) {
	var uidArr []string
	uidArr = append(uidArr, strings.Split(query.UIDs, ",")...)
	var kycedArr []bool
	for _, kycedString := range strings.Split(query.KYCed, ",") {
		kyced, err := strconv.ParseBool(kycedString)
		if err != nil {
			return nil, nil, err
		}
		kycedArr = append(kycedArr, kyced)
	}
	if len(uidArr) >= s.maxBatchSize {
		return nil, nil, errors.Errorf("batch size is too big (current size %v, max size=%v)", len(uidArr), s.maxBatchSize)
	}
	if len(uidArr) != len(kycedArr) {
		return nil, nil, errors.New("len uids and kyced are not match")
	}
	return uidArr, kycedArr, nil
}

// stats-batch returns cap of the user with given uids, max size = 1k
func (s *Server) userStatsBatch(c *gin.Context) {
	var (
		logger = s.sugar.With("func", "users/http/Server.userStats")
		input  userStatsBatchQuery

		userCap *big.Int
		rich    bool
	)
	if err := c.ShouldBindQuery(&input); err != nil {
		httputil.ResponseFailure(c, http.StatusBadRequest, err)
		return
	}

	logger = logger.With(
		"uid", input.UIDs,
		"kyced", input.KYCed,
	)
	logger.Debugw("querying stats batch for user")

	uidArr, kycedArr, err := s.convertQueryParams(input)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("failed to parse params: %s", err.Error())},
		)
		return
	}

	volume, err := s.getUserVolumeByUIDs(uidArr)
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
	//output
	var jsonOutput []gin.H
	for i := range uidArr {
		userCap = blockchain.EthToWei(s.userCapConf.UserCap(kycedArr[i]).TxLimit / rate)
		volumeInWei := blockchain.EthToWei(volume[i] / rate)
		rich = s.userCapConf.IsRich(kycedArr[i], volume[i])
		logger.Infow("got last 24h volume of user",
			"volume", volumeInWei,
			"cap", userCap,
			"rich", rich,
		)

		jsonOutput = append(jsonOutput, gin.H{
			"cap":    userCap,
			"kyced":  kycedArr[i],
			"rich":   rich,
			"volume": volumeInWei,
		})
	}
	c.JSON(http.StatusOK, jsonOutput)
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
	s.r.GET("/users-batch", s.userStatsBatch)
}

//Run start server and serve
func (s *Server) Run() error {
	s.register()
	return s.r.Run(s.host)
}
