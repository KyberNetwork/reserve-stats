package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/users/http"
	"github.com/KyberNetwork/reserve-stats/users/storage"
	"github.com/KyberNetwork/tokenrate/coingecko"
	"github.com/urfave/cli"
)

const defaultDB = "users"

func main() {
	app := libapp.NewApp()
	app.Name = "User stat module"
	app.Usage = "Store and return user stat information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.UsersPort)...)
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
	sugar.Info("Run user module")

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	userDB, err := storage.NewDB(
		sugar,
		db,
	)
	if err != nil {
		return err
	}
	defer userDB.Close()

	// Store trade logs into influx DB
	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	influxStorage, err := storage.NewInfluxStorage(
		sugar,
		common.DatabaseName,
		influxClient,
	)
	if err != nil {
		return err
	}

	server := http.NewServer(sugar, coingecko.New(), userDB,
		httputil.NewHTTPAddressFromContext(c), influxStorage)
	return server.Run()
}
