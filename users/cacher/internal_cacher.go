package cacher

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
)

const (
	uidPrefix = "uid"
)

// InternalRedisCacher is instance for redis cache
type InternalRedisCacher struct {
	sugar       *zap.SugaredLogger
	redisClient *redis.Client
	expiration  time.Duration
	db          *sqlx.DB
}

// NewInternalRedisCacher returns a new redis cacher instance
func NewInternalRedisCacher(sugar *zap.SugaredLogger,
	redisClient *redis.Client, expiration time.Duration,
	db *sqlx.DB) *InternalRedisCacher {
	return &InternalRedisCacher{
		sugar:       sugar,
		redisClient: redisClient,
		expiration:  expiration,
		db:          db,
	}
}

// Cache24hVolume cache user volume daily by uid
func (irc *InternalRedisCacher) Cache24hVolume() error {
	if err := irc.cache24hVolumeByUID(); err != nil {
		return err
	}
	return nil
}

// User24hVolume ...
type User24hVolume struct {
	UserID int64   `db:"user_id"`
	Volume float64 `db:"volume"`
}

func (irc *InternalRedisCacher) cache24hVolumeByUID() error {
	var (
		logger = irc.sugar.With("func", caller.GetCurrentFunctionName())
		res    []User24hVolume
	)

	// read total trade 24h
	query := `
		SELECT 
			users.id as user_id,
			SUM(split.eth_amount * tradelogs.eth_usd_rate) as volume
		FROM split
		JOIN tradelogs on tradelogs.id = split.trade_id
		JOIN users on tradelogs.user_address_id = users.id
		WHERE tradelogs.timestamp > now() - interval '24' hour 
		GROUP BY user.id
		`

	logger.Debugw("query", "query 24h fiat volume", query)
	if err := irc.db.Select(&res, query); err != nil {
		logger.Errorw("failed to get user volume 24 hours", "error", err)
		return err
	}

	pipe := irc.redisClient.Pipeline()
	for _, user := range res {
		uid := user.UserID

		// save to cache with configured expiration duration
		if err := irc.pushToPipeline(pipe, fmt.Sprintf("%s:%d", uidPrefix, uid), user.Volume, irc.expiration); err != nil {
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
