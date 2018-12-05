package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/appnames"
	"github.com/KyberNetwork/reserve-stats/lib/core"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/userprofile"
	"github.com/KyberNetwork/reserve-stats/tradelogs/http"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const (
	defaultDB = "users"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs HTTP Api"
	app.Usage = "Serve trade logs data"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
		if err := libapp.Validate(c); err != nil {
			return err
		}

		logger, err := libapp.NewLogger(c)
		if err != nil {
			return err
		}
		defer logger.Sync()

		sugar := logger.Sugar()
		coreClient, err := core.NewClientFromContext(sugar, c)
		if err != nil {
			return err
		}

		coreCachedClient := core.NewCachedClient(coreClient)
		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return err
		}

		userPostgres, err := libapp.NewDBFromContext(c)
		if err != nil {
			return err
		}

		kycChecker := storage.NewUserKYCChecker(sugar, userPostgres)

		influxStorage, err := storage.NewInfluxStorage(
			sugar,
			"trade_logs",
			influxClient,
			coreCachedClient,
			kycChecker,
		)
		if err != nil {
			return err
		}

		var options []http.ServerOption
		addrToAppName, err := appnames.NewClientFromContext(sugar, c)
		if err != nil {
			return err
		}
		if addrToAppName != nil {
			options = append(options, http.WithApplicationNames(addrToAppName))
		}
		addrToAppName, err = appnames.NewClientFromContext(sugar, c)
		if err != nil {
			return err
		}

		userClient, err := userprofile.NewClientFromContext(sugar, c)
		if err != nil {
			return err
		}

		cachedUserClient := userprofile.NewCachedClientFromContext(userClient, c)
		if cachedUserClient != nil {
			options = append(options, http.WithUserProfile(cachedUserClient))
		}
		api := http.NewServer(influxStorage, httputil.NewHTTPAddressFromContext(c),
			sugar, coreCachedClient,
			options...)
		err = api.Start()
		if err != nil {
			return err
		}

		if err = api.Start(); err != nil {
			return err
		}

		return nil
	}

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.TradeLogsPort)...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, core.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags, appnames.NewCliFlags()...)
	app.Flags = append(app.Flags, userprofile.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
