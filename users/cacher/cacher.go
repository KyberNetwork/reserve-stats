package cacher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

const (
	influxDB   = "trade_logs"
	expireTime = time.Hour
)

//RedisCacher is instance for redis cache
type RedisCacher struct {
	sugar          *zap.SugaredLogger
	postgresDB     *storage.UserDB
	influxDBClient client.Client
	redisClient    *redis.Client
	kycedCap       *common.UserCap
	nonKycedCap    *common.UserCap
}

//NewRedisCacher returns a new redis cacher instance
func NewRedisCacher(sugar *zap.SugaredLogger, postgresDB *storage.UserDB,
	influxDBClient client.Client, redisClient *redis.Client) *RedisCacher {
	return &RedisCacher{
		sugar:          sugar,
		postgresDB:     postgresDB,
		influxDBClient: influxDBClient,
		redisClient:    redisClient,
		kycedCap:       common.NewUserCap(true),
		nonKycedCap:    common.NewUserCap(false),
	}
}

//CacheUserInfo save user info to redis cache
func (rc *RedisCacher) CacheUserInfo() error {
	if err := rc.cacheAllKycedUsers(); err != nil {
		return err
	}
	if err := rc.cacheRichUser(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCacher) cacheAllKycedUsers() error {
	var (
		logger    = rc.sugar.With("func", "user/cacher/cacheAllKycedUsers")
		addresses []string
		err       error
	)
	// read all address from addresses table in postgres
	if addresses, err = rc.postgresDB.GetAllAddresses(); err != nil {
		logger.Errorw("error from query postgres db", "error", err.Error())
		return err
	}
	logger.Debugw("addresses from postgres", "addresses", addresses)

	for _, address := range addresses {
		user := common.UserResponse{
			KYCed: true,
		}
		rc.saveToCache(address, user, 0)
	}
	return nil
}

func (rc *RedisCacher) cacheRichUser() error {
	var (
		logger = rc.sugar.With("func", "user/cacher/cacheRichUser")
	)
	// read total trade 24h
	query := fmt.Sprintf(`SELECT SUM(amount) as daily_fiat_amount FROM 
	(SELECT %s*%s as amount FROM trades WHERE time >= (now()-24h))
	GROUP BY user_addr`, logSchema.EthAmount.String(), logSchema.EthUSDRate.String())

	logger.Debugw("query", "query 24h trades", query)

	res, err := influxdb.QueryDB(rc.influxDBClient, query, influxDB)
	if err != nil {
		logger.Errorw("error from query", "err", err)
		return err
	}

	// loop all user, check kyced
	if len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 || len(res[0].Series[0].Values[0]) < 2 {
		logger.Debugw("influx db is empty", "result", res)
		return nil
	}

	for _, serie := range res[0].Series {
		userAddress := serie.Tags[logSchema.UserAddr.String()]
		// check kyced
		kyced, err := rc.isKyced(userAddress)
		if err != nil {
			return err
		}

		// check rich
		userTradeAmount, err := influxdb.GetFloat64FromInterface(serie.Values[0][1])
		if err != nil {
			logger.Errorw("values second should be a float", "value", serie.Values[0][1])
			return nil
		}

		if (kyced && userTradeAmount < rc.kycedCap.DailyLimit) ||
			(!kyced && userTradeAmount < rc.nonKycedCap.DailyLimit) {
			// if user is not rich then it is already cached before
			continue
		}
		user := common.UserResponse{
			Rich:  true,
			KYCed: kyced,
		}

		// save to cache with 1 hour
		rc.saveToCache(userAddress, user, expireTime)
	}
	return nil
}

func (rc *RedisCacher) saveToCache(key string, value common.UserResponse, expireTime time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		rc.sugar.Debugw("Cannot marshal value", "error", err)
		return err
	}
	if err := rc.redisClient.Set(key, data, expireTime).Err(); err != nil {
		rc.sugar.Debugw("set cache to redis error", "error", err)
		return err
	}
	rc.sugar.Debugw("save data to cache succes", "key", key, "value", value)
	return nil
}

func (rc *RedisCacher) isKyced(userAddress string) (bool, error) {
	if err := rc.redisClient.Get(userAddress).Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		}
		rc.sugar.Debugw("get data from redis failed", "address", userAddress, "error", err.Error())
		return false, err
	}
	return true, nil
}
