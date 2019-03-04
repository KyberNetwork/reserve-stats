package cacher

import (
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

const (
	influxDB    = "trade_logs"
	richPrefix  = "rich"
	kycedPrefix = "kyced"
)

var kycedAddresses map[string]int

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
func (rc *RedisCacher) CacheUserInfo(expireTime time.Duration) error {
	if err := rc.cacheAllKycedUsers(expireTime); err != nil {
		return err
	}
	if err := rc.cacheRichUser(expireTime); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCacher) cacheAllKycedUsers(expireTime time.Duration) error {
	var (
		logger    = rc.sugar.With("func", "user/cacher/cacheAllKycedUsers")
		err       error
		addresses []string
	)
	kycedAddresses = make(map[string]int)
	// read all address from addresses table in postgres
	if addresses, err = rc.postgresDB.GetAllAddresses(); err != nil {
		logger.Errorw("error from query postgres db", "error", err.Error())
		return err
	}
	logger.Debugw("addresses from postgres", "addresses", addresses)

	pipe := rc.redisClient.Pipeline()
	for _, address := range addresses {
		// push to redis
		addressHex := ethereum.HexToAddress(address).Hex()
		if err := rc.pushToPipeline(pipe, fmt.Sprintf("%s:%s", kycedPrefix, addressHex), expireTime); err != nil {
			if dErr := pipe.Discard(); dErr != nil {
				err = fmt.Errorf("%s - %s", dErr.Error(), err.Error())
			}
			return err
		}
		// push to mem
		kycedAddresses[address] = 1
	}
	_, err = pipe.Exec()
	return err
}

func (rc *RedisCacher) cacheRichUser(expireTime time.Duration) error {
	var (
		logger = rc.sugar.With("func", "user/cacher/cacheRichUser")
	)

	// read total trade 24h
	query := fmt.Sprintf(`SELECT SUM(amount) as daily_fiat_amount FROM 
	(SELECT %s*%s as amount FROM trades WHERE time >= (now()-24h)) GROUP BY user_addr`, logSchema.EthAmount.String(), logSchema.EthUSDRate.String())

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

	pipe := rc.redisClient.Pipeline()
	for _, serie := range res[0].Series {
		userAddress := serie.Tags[logSchema.UserAddr.String()]
		// check kyced
		_, kyced := kycedAddresses[userAddress]

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

		// save to cache with 1 hour
		if err := rc.pushToPipeline(pipe, fmt.Sprintf("%s:%s", richPrefix, userAddress), expireTime); err != nil {
			if dErr := pipe.Discard(); dErr != nil {
				err = fmt.Errorf("%s - %s", dErr.Error(), err.Error())
			}
			return err
		}
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}

	return err
}

func (rc *RedisCacher) pushToPipeline(pipeline redis.Pipeliner, key string, expireTime time.Duration) error {
	if err := pipeline.Set(key, 1, expireTime).Err(); err != nil {
		rc.sugar.Debugw("set cache to redis error", "error", err)
		return err
	}

	rc.sugar.Debugw("save data to cache succes", "key", key)
	return nil
}
