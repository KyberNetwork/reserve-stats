package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	rediscache "github.com/KyberNetwork/reserve-stats/lib/redis"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/go-redis/redis"
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

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags, rediscache.NewCliFlags()...)

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

	cached, err := rediscache.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Debugw("Initiated redis cached", "cache", cached)

	return nil
}

func cachedUserInfo(postgresDB *sqlx.DB, cached *redis.Client) error {
	// cache all user address
	
	return nil
}
