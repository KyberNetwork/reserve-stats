package userprofile

import (
	"fmt"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	userprofileURLFlag        = "user-profile-url"
	userprofileSigningKeyFlag = "user-profile-signing-key"
	maxUserCacheFlag          = "max-user-profile-cache"
	redisCacheFlag            = "use-redis-cache"
	maxUserCacheDefault       = 1000
	redisUserProfileDBDefault = 0
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags() []cli.Flag {
	userprofileFlags := []cli.Flag{
		cli.StringFlag{
			Name:   userprofileURLFlag,
			Usage:  "user profile API URL",
			EnvVar: "USER_PROFILE_URL",
		},
		cli.BoolFlag{
			Name:   redisCacheFlag,
			Usage:  "this flag is set if redis cache is preferred fro user profile",
			EnvVar: "USE_REDIS_CACHE",
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
	}
	userprofileFlags = append(userprofileFlags, libapp.NewRedisFlags(redisUserProfileDBDefault)...)
	return userprofileFlags
}

// newClientFromContext returns new core client from cli flags.
func newClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
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

//newInMemCachedFromContext create the inmem cache client from flag agruments
func newInMemCachedFromContext(client *Client, c *cli.Context) Interface {
	maxCacheSize := c.Int64(maxUserCacheFlag)
	return newCachedClient(client, maxCacheSize)
}

// NewUserProfileCachedClientFromContext return a cached user client from context
func NewUserProfileCachedClientFromContext(c *cli.Context, sugar *zap.SugaredLogger) (Interface, error) {
	userClient, err := newClientFromContext(sugar, c)
	if err != nil {
		return nil, err
	}

	if c.Bool(redisCacheFlag) {
		sugar.Infow("use redis cache for user profile")
		redisClient, err := libapp.NewRedisClientFromContext(c)
		if err != nil {
			return nil, err
		}
		return newRedisCachedClient(userClient, redisClient), nil
	}

	sugar.Infow("use default in-mem cache for user profile ")
	return newInMemCachedFromContext(userClient, c), nil
}
