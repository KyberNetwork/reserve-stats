package cacher

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/tradelog"
)

const (
	uidPrefix = "uid"
)

//InternalRedisCacher is instance for redis cache
type InternalRedisCacher struct {
	sugar          *zap.SugaredLogger
	influxDBClient client.Client
	redisClient    *redis.Client
	expiration     time.Duration
}

//NewInternalRedisCacher returns a new redis cacher instance
func NewInternalRedisCacher(sugar *zap.SugaredLogger, influxDBClient client.Client,
	redisClient *redis.Client, expiration time.Duration) *InternalRedisCacher {
	return &InternalRedisCacher{
		sugar:          sugar,
		influxDBClient: influxDBClient,
		redisClient:    redisClient,
		expiration:     expiration,
	}
}

//Cache24hVolume cache user volume daily by uid
func (irc *InternalRedisCacher) Cache24hVolume() error {
	if err := irc.cache24hVolumeByUID(); err != nil {
		return err
	}
	return nil
}

func (irc *InternalRedisCacher) cache24hVolumeByUID() error {
	var (
		logger = irc.sugar.With("func", caller.GetCurrentFunctionName())
	)

	// read total trade 24h
	query := fmt.Sprintf(`SELECT SUM(amount) FROM (SELECT %s*%s as amount FROM trades WHERE time >= now()-24h AND uid != '') WHERE time >= now()-24h  GROUP BY uid`,
		logSchema.EthAmount.String(),
		logSchema.EthUSDRate.String())

	logger.Debugw("query", "query 24h fiat volume", query)

	res, err := influxdb.QueryDB(irc.influxDBClient, query, influxDB)
	if err != nil {
		logger.Errorw("error from query", "err", err)
		return err
	}

	if len(res) == 0 || len(res[0].Series) == 0 || len(res[0].Series[0].Values) == 0 || len(res[0].Series[0].Values[0]) < 2 {
		logger.Debugw("influx db is empty", "result", res)
		return nil
	}

	pipe := irc.redisClient.Pipeline()
	for _, serie := range res[0].Series {
		uid := serie.Tags[logSchema.UID.String()]

		// check rich
		userTradeAmount, err := influxdb.GetFloat64FromInterface(serie.Values[0][1])
		if err != nil {
			logger.Errorw("values second should be a float", "value", serie.Values[0][1])
			return nil
		}

		// save to cache with configured expiration duration
		if err := irc.pushToPipeline(pipe, fmt.Sprintf("%s:%s", uidPrefix, uid), userTradeAmount, irc.expiration); err != nil {
			if dErr := pipe.Discard(); dErr != nil {
				err = fmt.Errorf("%s - %s", dErr.Error(), err.Error())
			}
			return err
		}
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}

	return nil
}

func (irc *InternalRedisCacher) pushToPipeline(pipeline redis.Pipeliner, key string, value float64, expireTime time.Duration) error {
	var (
		logger = irc.sugar.With("func", caller.GetCurrentFunctionName())
	)
	if err := pipeline.Set(key, value, expireTime).Err(); err != nil {
		logger.Debugw("set cache to redis error", "error", err)
		return err
	}

	logger.Debugw("save data to cache succes", "key", key)
	return nil
}
