package userprofile

import (
	"fmt"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis"
)

const (
	nameField                  = "name"
	idField                    = "id"
	expireDurationForNilResult = 24 * time.Hour
)

// RedisCachedClient is the wrapper of User Profile Client with caching ability.
type RedisCachedClient struct {
	*Client
	redisClient *redis.Client
}

// NewRedisCachedClient creates a new User Profile cached client instance.
func NewRedisCachedClient(client *Client, redisClient *redis.Client) *RedisCachedClient {
	return &RedisCachedClient{
		Client:      client,
		redisClient: redisClient,
	}
}

// LookUpCache will lookup the Userprofile from cache
// Redis cache will return error if the addr is not exist.
func (cc *RedisCachedClient) LookUpCache(addr ethereum.Address) (UserProfile, error) {
	ps, err := cc.redisClient.HMGet(addr.Hex(), nameField, idField).Result()
	if err != nil {
		return UserProfile{}, err
	}
	if len(ps) < 2 {
		return UserProfile{}, fmt.Errorf("result cached in redis returned wrong len")
	}
	name, ok := ps[0].(string)
	if !ok {
		return UserProfile{}, fmt.Errorf("cannot assert name field (%v) to string type", ps[0])
	}
	pID, ok := ps[1].(string)
	if !ok {
		return UserProfile{}, fmt.Errorf("cannot assert pID field(%v) to string type", ps[1])
	}
	pIDint, err := strconv.ParseInt(pID, 10, 64)
	if err != nil {
		return UserProfile{}, err
	}
	return UserProfile{
		UserName:  name,
		ProfileID: int64(pIDint),
	}, nil
}

// LookUpUserProfile will look for the UserProfile of input addr in cache first
// If this fail then it will query from endpoint
func (cc *RedisCachedClient) LookUpUserProfile(addr ethereum.Address) (UserProfile, error) {
	logger := cc.sugar.With(
		"func", "lib/core/RedisCachedClient.Token",
		"address", addr.Hex(),
	)
	p, err := cc.LookUpCache(addr)
	if err == nil {
		logger.Debugw("cache hit")
		return p, nil
	}

	logger.Debugw("cache missed", "redis error", err)

	p, err = cc.Client.LookUpUserProfile(addr)
	if err != nil {
		return p, err
	}

	//cache the result into redis and return result
	err = cc.cacheRedis(addr, p)
	if err != nil {
		logger.Debugw("cache to redis failed", "redis error", err)
	}
	return p, err
}

func (cc *RedisCachedClient) cacheRedis(addr ethereum.Address, p UserProfile) error {
	//empty result
	cc.sugar.Debugf("cached user profile from API into redis", "address", addr.Hex(), "user profile", p)
	resultMap := map[string]interface{}{
		nameField: p.UserName,
		idField:   p.ProfileID,
	}
	if _, err := cc.redisClient.HMSet(addr.Hex(), resultMap).Result(); err != nil {
		return err
	}
	if (p == UserProfile{}) {
		return cc.redisClient.Expire(addr.Hex(), expireDurationForNilResult).Err()
	}
	return nil
}
