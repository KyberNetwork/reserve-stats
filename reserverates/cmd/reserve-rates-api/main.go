package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/reserverates/http"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage/postgres"
	"github.com/urfave/cli"
)

const (
	defaultPostgresDB = "reserve_rates"
)

func newServerCli() *cli.App {
	app := libapp.NewApp()
	app.Name = "reserve-rates-api"
	app.Usage = "server for query rate API"
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.ReserveRatesPort)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostgresDB)...)
	app.Action = func(c *cli.Context) error {
		if err := libapp.Validate(c); err != nil {
			return err
		}

		sugar, flusher, err := libapp.NewSugaredLogger(c)
		if err != nil {
			return err
		}
		defer flusher()

		var rateStorage storage.ReserveRatesStorage
		db, err := libapp.NewDBFromContext(c)
		if err != nil {
			return err
		}
		if rateStorage, err = postgres.NewPostgresStorage(db, sugar, nil); err != nil {
			return err
		}

		hostStr := httputil.NewHTTPAddressFromContext(c)
		server, err := http.NewServer(hostStr, rateStorage, sugar)
		if err != nil {
			return err
		}
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
