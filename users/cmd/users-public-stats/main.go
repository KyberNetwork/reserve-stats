package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	_ "github.com/KyberNetwork/reserve-stats/lib/httputil/validators" // import custom validator functions
	libredis "github.com/KyberNetwork/reserve-stats/lib/redis"
	server "github.com/KyberNetwork/reserve-stats/users/public-server"
	"github.com/KyberNetwork/tokenrate/coingecko"
)

func main() {
	app := libapp.NewApp()
	app.Name = "User stat public service"
	app.Usage = "Return user stat information from cache"
	app.Action = run
	app.Version = "0.1"

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.UsersPublicPort)...)
	app.Flags = append(app.Flags, libredis.NewCliFlags()...)
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

	redisClient, err := libredis.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Debugw("initiate redis client", "client", redisClient)

	publicServer := server.NewServer(sugar, httputil.NewHTTPAddressFromContext(c), coingecko.New(), redisClient)

	return publicServer.Run()
}
