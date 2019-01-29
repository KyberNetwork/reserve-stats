package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/redis"
)

func main() {
	app := app.NewApp()
	app.Name = "User stat public service"
	app.Usage = "Return user stat information from cache"
	app.Action = run
	app.Version = "0.1"

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.UsersPublicPort)...)
	app.Flags = append(app.Flags, redis.NewCliFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := app.Validate(c); err != nil {
		return err
	}

	logger, err := app.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("Run user stats public service")

	cachedClient, err := redis.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Debugw("initiate redis client", "client", cachedClient)

	return nil
}
