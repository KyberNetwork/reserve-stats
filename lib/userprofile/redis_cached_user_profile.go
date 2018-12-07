package userprofile

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/go-redis/redis"
)

const (
	separator = ":"
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
	result, err := cc.redisClient.Get(addr.Hex()).Result()
	if err != nil {
		return UserProfile{}, err
	}
	// empty
	if result == "" {
		return UserProfile{
			UserName:  "",
			ProfileID: 0,
		}, nil
	}
	data := strings.Split(result, separator)
	if len(data) < 2 {
		return UserProfile{}, fmt.Errorf("Redis cached wrong result string, expect format name:id, get %s", result)
	}
	pID, err := strconv.ParseInt(data[1], 10, 64)
	if err != nil {
		return UserProfile{}, err
	}
	return UserProfile{
		UserName:  data[0],
		ProfileID: pID,
	}, err
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
	cc.sugar.Debugf("get user profile from API", "user profile", p)
	if (p == UserProfile{}) {
		return cc.redisClient.Set(addr.Hex(), "", 24*time.Hour).Err()
	}

	resultStr := fmt.Sprintf("%s:%s", p.UserName, strconv.FormatInt(p.ProfileID, 10))
	return cc.redisClient.Set(addr.Hex(), resultStr, 0).Err()
}
