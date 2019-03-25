package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/http"
	rrpostgres "github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/urfave/cli"
)

const (
	defaultPostGresDB = common.DefaultDB
)

func newServerCli() *cli.App {
	app := libapp.NewApp()
	app.Name = "reserve-rates-api"
	app.Usage = "server for query accounting reserve rate API"
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.ReserveRatesPort)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostGresDB)...)
	app.Action = run
	return app
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

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	ratesStorage, err := rrpostgres.NewDB(sugar, db)
	if err != nil {
		return err
	}
	defer func(err *error) {
		cErr := ratesStorage.Close()
		if err == nil {
			*err = cErr
		} else {
			sugar.Error("DB closing failed", "error", cErr)
		}
	}(&err)
	hostStr := httputil.NewHTTPAddressFromContext(c)
	server, err := http.NewServer(hostStr, ratesStorage, sugar)
	if err != nil {
		return err
	}
	err = server.Run()
	return err
}

//reserverates --addresses=0xABCDEF,0xDEFGHI --block 100
func main() {
	app := newServerCli()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
