package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/http"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/storage"
	"github.com/urfave/cli"
)

const (
	defaultDB = "price_analytics"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Price analytics data"
	app.Usage = "store price analytic data"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.PriceAnalytic)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	sugar.Info("Run price analytic module")

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	priceDB, err := storage.NewPriceStorage(
		sugar,
		db,
	)
	if err != nil {
		return err
	}
	defer priceDB.Close()

	server := http.NewHTTPServer(sugar, httputil.NewHTTPAddressFromContext(c),
		priceDB)

	return server.Run()
}
