package main

import (
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/http"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs HTTP Api"
	app.Usage = "Serve trade logs data"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
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

		influxStorage, err := storage.NewInfluxStorage(
			sugar,
			"trade_logs",
			influxClient,
			coreCachedClient,
		)
		if err != nil {
			return err
		}

		api := http.NewServer(influxStorage, httputil.NewHTTPAddressFromContext(c), sugar, coreCachedClient)
		err = api.Start()
		if err != nil {
			return err
		}

		return nil
	}

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.TradeLogsPort)...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, core.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
