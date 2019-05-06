package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	libredis "github.com/KyberNetwork/reserve-stats/lib/redis"
	"github.com/KyberNetwork/reserve-stats/users/cacher"
	"github.com/KyberNetwork/reserve-stats/users/common"
)

const (
	expireTimeFlag    = "expire-time"
	defaultExpireTime = 3600 // 1 hour
)

func main() {
	app := libapp.NewApp()
	app.Name = "app for caching user information"
	app.Usage = "cache user info hourly"
	app.Action = run
	app.Version = "0.1"

	app.Flags = append(app.Flags, libredis.NewCliFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, common.NewUserCapCliFlags()...)
	app.Flags = append(app.Flags,
		cli.IntFlag{
			Name:   expireTimeFlag,
			Usage:  "Time to expire redis cache, count by second",
			EnvVar: "EXPIRE_TIME",
			Value:  defaultExpireTime,
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	sugar.Info("Run user public cacher")

	influxDBClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	redisCacheClient, err := libredis.NewClientFromContext(c)
	if err != nil {
		return err
	}

	expireTimeSecond := c.Int64(expireTimeFlag)
	expireTime := time.Duration(expireTimeSecond) * time.Second
	userCapConf := common.NewUserCapConfigurationFromContext(c)
	sugar.Debugw("Initiated redis cached", "cache", redisCacheClient)

	redisCacher := cacher.NewRedisCacher(sugar, influxDBClient, redisCacheClient, expireTime, userCapConf)

	return redisCacher.CacheUserInfo()
}
