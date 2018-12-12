package app

import (
	"github.com/go-redis/redis"
	"github.com/urfave/cli"
)

const (
	redisEndpointFlag = "redis-endpoint"
	redisPasswordFlag = "redis-password"
	redisDBFlag       = "redis-db"
)

func NewRedisFlags(defaultDB int) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   redisEndpointFlag,
			Usage:  "redis connection endpoint, if  this is not set the default mem cache will be use instead of redis",
			EnvVar: "REDIS_ENDPOINT",
			Value:  "",
		},
		cli.IntFlag{
			Name:   redisDBFlag,
			Usage:  "Database for redis user profile cache. Default to 0",
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
	if redisURL == "" {
		return nil, nil
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: c.String(redisPasswordFlag),
		DB:       c.Int(redisDBFlag),
	})
	_, err := redisClient.Ping().Result()
	return redisClient, err
}
