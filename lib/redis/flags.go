package redis

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-redis/redis"
	"github.com/urfave/cli"
)

const (
	redisEndpointFlag = "redis-endpoint"
	redisDBFlag       = "redis-db"
	redisPasswordFlag = "redis-password"

	redisEndpointDefaultValue = "localhost:6379"
	redisDBDefaultValue       = 0
)

//NewCliFlags return flags for redis service
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   redisEndpointFlag,
			Usage:  "Endpoint to redis",
			EnvVar: "REDIS_ENDPOINT",
			Value:  redisEndpointDefaultValue,
		},
		cli.StringFlag{
			Name:   redisPasswordFlag,
			Usage:  "Redis password",
			EnvVar: "REDIS_PASSWORD",
		},
		cli.IntFlag{
			Name:   redisDBFlag,
			Usage:  "Redis database",
			EnvVar: "REDIS_DB",
			Value:  redisDBDefaultValue,
		},
	}
}

//NewClientFromContext return new redis client
func NewClientFromContext(c *cli.Context) (*redis.Client, error) {
	redisEndpoint := c.String(redisEndpointFlag)
	redisPassword := c.String(redisPasswordFlag)

	if err := validation.Validate(redisEndpoint, validation.Required, is.URL); err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: redisPassword,
	})
	_, err := redisClient.Ping().Result()
	return redisClient, err
}
