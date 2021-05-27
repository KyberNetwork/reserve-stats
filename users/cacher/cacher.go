package cacher

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
)

const (
	richPrefix = "rich"
)

// RedisCacher is instance for redis cache
type RedisCacher struct {
	sugar       *zap.SugaredLogger
	redisClient *redis.Client
	expiration  time.Duration
	db          *sqlx.DB
}

// NewRedisCacher returns a new redis cacher instance
func NewRedisCacher(sugar *zap.SugaredLogger, redisClient *redis.Client, db *sqlx.DB, expiration time.Duration) *RedisCacher {
	return &RedisCacher{
		sugar:       sugar,
		redisClient: redisClient,
		expiration:  expiration,
		db:          db,
		//userCapConf:    userCapConf,
	}
}

// CacheUserInfo save user info to redis cache
func (rc *RedisCacher) CacheUserInfo() error {
	if err := rc.cacheRichUser(); err != nil {
		return err
	}
	return nil
}

type UserDailyVolume struct {
	UserAddress     string  `db:"user_addr"`
	DailyFiatAmount float64 `db:"daily_fiat_amount"`
}

func (rc *RedisCacher) cacheRichUser() error {
	var (
		logger = rc.sugar.With("func", caller.GetCurrentFunctionName())
		result []UserDailyVolume
	)

	// read total trade 24h
	query := `SELECT 
		users.address as user_addr,
		sum(eth_amount*eth_usd_rate) as daily_fiat_amount
		FROM tradelogs WHERE timestamp > now() - interval 24 hour 
		JOIN users ON tradelogs.user_address_id = users.id
		GROUP BY users.address`

	logger.Debugw("query", "query 24h trades", query)

	err := rc.db.Select(&result, query)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		logger.Errorw("error from query", "err", err)
		return err
	}

	pipe := rc.redisClient.Pipeline()
	for _, r := range result {
		userAddress := r.UserAddress

		// check rich
		userTradeAmount := r.DailyFiatAmount

		// save to cache with configured expiration duration
		if err := rc.pushToPipeline(pipe, fmt.Sprintf("%s:%s", richPrefix, userAddress), userTradeAmount, rc.expiration); err != nil {
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

func (rc *RedisCacher) pushToPipeline(pipeline redis.Pipeliner, key string, value float64, expireTime time.Duration) error {
	if err := pipeline.Set(key, value, expireTime).Err(); err != nil {
		rc.sugar.Debugw("set cache to redis error", "error", err)
		return err
	}

	rc.sugar.Debugw("save data to cache succes", "key", key)
	return nil
}
