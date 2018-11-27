package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/http"
	influxRateStorage "github.com/KyberNetwork/reserve-stats/reserverates/storage/influx"
	"github.com/urfave/cli"
)

func newServerCli() *cli.App {
	app := libapp.NewApp()
	app.Name = "reserve-rates-api"
	app.Usage = "server for query rate API"
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.ReserveRatesPort)...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Action = func(c *cli.Context) error {
		if err := libapp.Validate(c); err != nil {
			return err
		}

		logger, err := libapp.NewLogger(c)
		if err != nil {
			return err
		}
		defer logger.Sync()

		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return err
		}

		rateStorage, err := influxRateStorage.NewRateInfluxDBStorage(logger.Sugar(), influxClient, common.DatabaseName, nil)
		if err != nil {
			return err
		}

		hostStr := httputil.NewHTTPAddressFromContext(c)
		server, err := http.NewServer(hostStr, rateStorage, logger.Sugar())
		return server.Run()
	}
	return app
}

//reserverates --addresses=0xABCDEF,0xDEFGHI --block 100
func main() {
	app := newServerCli()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
