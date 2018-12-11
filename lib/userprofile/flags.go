package userprofile

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-redis/redis"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	userprofileURLFlag        = "user-profile-url"
	userprofileSigningKeyFlag = "user-profile-signing-key"
	maxUserCacheFlag          = "max-user-profile-cache"
	maxUserCacheDefault       = 1000
	redisEndpointFlag         = "redis-endpoint"
	redisUserProfileDBFlag    = "redis-user-profile-db"
	redisUserProfileDBDefault = 0
	redisPasswordFlag         = "redis-password"
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   userprofileURLFlag,
			Usage:  "user profile API URL",
			EnvVar: "USER_PROFILE_URL",
		},
		cli.Int64Flag{
			Name:   maxUserCacheFlag,
			Usage:  "max Cache Size for user profile client. Default to 1000",
			EnvVar: "MAX_USER_CACHE",
			Value:  maxUserCacheDefault,
		},
		cli.StringFlag{
			Name:   userprofileSigningKeyFlag,
			Usage:  "user profile Signing Key",
			EnvVar: "USER_PROFILE_SIGNING_KEY",
		},
		cli.StringFlag{
			Name:   redisEndpointFlag,
			Usage:  "redis connection endpoint, if  this is not set the default mem cache will be use instead of redis",
			EnvVar: "REDIS_ENDPOINT",
			Value:  "",
		},
		cli.IntFlag{
			Name:   redisUserProfileDBFlag,
			Usage:  "Database for redis user profile cache. Default to 0",
			EnvVar: "REDIS_USER_PROFILE_DB",
			Value:  redisUserProfileDBDefault,
		},
		cli.StringFlag{
			Name:   redisPasswordFlag,
			Usage:  "redis connection password",
			EnvVar: "REDIS_PASSWORD",
			Value:  "",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	userURL := c.String(userprofileURLFlag)
	if userURL == "" {
		return nil, nil
	}
	err := validation.Validate(userURL,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("user profile url: %s", err.Error())
	}
	signingKey := c.String(userprofileSigningKeyFlag)
	err = validation.Validate(signingKey,
		validation.Required,
	)
	if err != nil {
		return nil, fmt.Errorf("user signing key: %s", err.Error())
	}
	return NewClient(sugar, userURL, signingKey)
}

//NewRedisClientFromContext creates redis client from flag agruments
func NewRedisClientFromContext(c *cli.Context) (*redis.Client, error) {
	redisURL := c.String(redisEndpointFlag)
	if redisURL == "" {
		return nil, nil
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: c.String(redisPasswordFlag),
		DB:       c.Int(redisUserProfileDBFlag),
	})
	_, err := redisClient.Ping().Result()
	return redisClient, err
}

//NewInmemCachedFromContext create the inmem cache client from flag agruments
func NewInmemCachedFromContext(client *Client, c *cli.Context) Interface {
	maxCacheSize := c.Int64(maxUserCacheFlag)
	return NewCachedClient(client, maxCacheSize)
}
