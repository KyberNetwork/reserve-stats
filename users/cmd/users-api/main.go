package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/tokenrate/coingecko"
	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	usercommon "github.com/KyberNetwork/reserve-stats/users/common"
	"github.com/KyberNetwork/reserve-stats/users/http"
	"github.com/KyberNetwork/reserve-stats/users/storage"
)

func main() {
	app := libapp.NewApp()
	app.Name = "User stat module"
	app.Usage = "Store and return user stat information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(usercommon.DefaultDB)...)
	app.Flags = append(app.Flags, usercommon.NewUserCapCliFlags()...)
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

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

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

	userCapConf := usercommon.NewUserCapConfigurationFromContext(c)

	server := http.NewServer(
		sugar,
		coingecko.New(),
		userDB,
		httputil.NewHTTPAddressFromContext(c),
		influxStorage,
		userCapConf)
	return server.Run()
}
