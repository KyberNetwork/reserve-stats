package cacher

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	"github.com/KyberNetwork/reserve-stats/users/common"
)

const (
	influxDB   = "trade_logs"
	richPrefix = "rich"
)

//RedisCacher is instance for redis cache
type RedisCacher struct {
	sugar          *zap.SugaredLogger
	influxDBClient client.Client
	redisClient    *redis.Client
	expiration     time.Duration
	userCapConf    *common.UserCapConfiguration
}

//NewRedisCacher returns a new redis cacher instance
func NewRedisCacher(sugar *zap.SugaredLogger, influxDBClient client.Client,
	redisClient *redis.Client, expiration time.Duration, userCapConf *common.UserCapConfiguration) *RedisCacher {
	return &RedisCacher{
		sugar:          sugar,
		influxDBClient: influxDBClient,
		redisClient:    redisClient,
		expiration:     expiration,
		userCapConf:    userCapConf,
	}
}

//CacheUserInfo save user info to redis cache
func (rc *RedisCacher) CacheUserInfo() error {
	if err := rc.cacheRichUser(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCacher) cacheRichUser() error {
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

	if len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 || len(res[0].Series[0].Values[0]) < 2 {
		logger.Debugw("influx db is empty", "result", res)
		return nil
	}

	pipe := rc.redisClient.Pipeline()
	for _, serie := range res[0].Series {
		userAddress := serie.Tags[logSchema.UserAddr.String()]

		// check rich
		userTradeAmount, err := influxdb.GetFloat64FromInterface(serie.Values[0][1])
		if err != nil {
			logger.Errorw("values second should be a float", "value", serie.Values[0][1])
			return nil
		}

		if !rc.userCapConf.IsRich(true, userTradeAmount) {
			// if user is not rich then it is already cached before
			continue
		}

		// save to cache with configured expiration duration
		if err := rc.pushToPipeline(pipe, fmt.Sprintf("%s:%s", richPrefix, userAddress), rc.expiration); err != nil {
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
