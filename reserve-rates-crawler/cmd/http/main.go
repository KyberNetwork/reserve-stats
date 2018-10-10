package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/http"
	influxRateStorage "github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/storage/influx"

	"github.com/urfave/cli"
)

const (
	defaultPort   = 8002
	portEnvPrefix = "RESERVE_RATE"
)

func newServerCli() *cli.App {
	app := libapp.NewApp()
	app.Name = "reserve-rates-server"
	app.Usage = "server for query rate API"
	app.Flags = append(app.Flags, httputil.NewHTTPFlags(portEnvPrefix, defaultPort)...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Action = func(c *cli.Context) error {
		logger, err := libapp.NewLogger(c)
		if err != nil {
			return err
		}
		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return err
		}
		rateStorage, err := influxRateStorage.NewRateInfluxDBStorage(influxClient)
		if err != nil {
			return err
		}
		hostStr := httputil.NewHostFromContext(c)
		server, err := http.NewServer(hostStr, rateStorage, logger.Sugar())
		return server.Run()
	}
	return app
}

//reserve-rates-crawler --addresses=0xABCDEF,0xDEFGHI --block 100
func main() {
	app := newServerCli()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
