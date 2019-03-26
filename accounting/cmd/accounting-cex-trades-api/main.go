package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	tradestorage "github.com/KyberNetwork/reserve-stats/accounting/binance-storage/trade-storage"
	"github.com/KyberNetwork/reserve-stats/accounting/cex/http"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	huobistorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting Reserve Addresses"
	app.Usage = "Accounting Reserve Addresses Manager"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultDB)...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingCEXTradesPort)...)
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

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	hs, err := huobistorage.NewDB(sugar, db)
	if err != nil {
		return err
	}

	// TODO: make table name optional
	bs, err := tradestorage.NewDB(sugar, db, "binance_trades")
	if err != nil {
		return err
	}

	defer func() {
		if cErr := db.Close(); cErr != nil {
			sugar.Errorf("failed to close database: err=%s", cErr.Error())
		}
	}()

	s := http.NewServer(sugar, httputil.NewHTTPAddressFromContext(c), hs, bs)

	return s.Run()
}
