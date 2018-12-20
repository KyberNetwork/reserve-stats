package app

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-redis/redis"
	"github.com/urfave/cli"
)

const (
	redisEndpointFlag = "redis-endpoint"
	redisPasswordFlag = "redis-password"
	redisDBFlag       = "redis-db"
)

// NewRedisFlags returns a list of flag for Redis configuration
func NewRedisFlags(defaultDB int) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   redisEndpointFlag,
			Usage:  "redis connection endpoint",
			EnvVar: "REDIS_ENDPOINT",
			Value:  "",
		},
		cli.IntFlag{
			Name:   redisDBFlag,
			Usage:  "Database for redis cache. Default to 0",
			EnvVar: "REDIS_USER_PROFILE_DB",
			Value:  defaultDB,
		},
		cli.StringFlag{
			Name:   redisPasswordFlag,
			Usage:  "redis connection password",
			EnvVar: "REDIS_PASSWORD",
			Value:  "",
		},
	}
}

//NewRedisClientFromContext creates redis client from flag agruments
func NewRedisClientFromContext(c *cli.Context) (*redis.Client, error) {
	redisURL := c.String(redisEndpointFlag)
	if err := validation.Validate(redisURL, validation.Required, is.URL); err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: c.String(redisPasswordFlag),
		DB:       c.Int(redisDBFlag),
	})
	_, err := redisClient.Ping().Result()
	return redisClient, err
}
