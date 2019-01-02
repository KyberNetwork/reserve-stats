package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	libredis "github.com/KyberNetwork/reserve-stats/lib/redis"
	"github.com/KyberNetwork/reserve-stats/users/cacher"
	"github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

const (
	defaultDB = "users"
)

func main() {
	app := libapp.NewApp()
	app.Name = "app for caching user information"
	app.Usage = "cache user info hourly"
	app.Action = run
	app.Version = "0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultDB)...)
	app.Flags = append(app.Flags, libredis.NewCliFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("Run user public cacher")

	postgresDB, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	userDB, err := storage.NewDB(
		sugar,
		postgresDB,
	)
	sugar.Debugw("Initiated postgres client", "client", userDB)

	influxDBClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	redisCacheClient, err := libredis.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Debugw("Initiated redis cached", "cache", redisCacheClient)

	redisCacher := cacher.NewRedisCacher(sugar, userDB, influxDBClient, redisCacheClient)

	return redisCacher.CacheUserInfo()
}