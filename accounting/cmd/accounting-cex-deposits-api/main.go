package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/depositstorage"
	"github.com/KyberNetwork/reserve-stats/accounting/cex-deposit/http"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

func newServerCli() *cli.App {
	app := libapp.NewApp()
	app.Name = "cex-deposit-api"
	app.Usage = "server for query accounting cex deposit"
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingCEXDepositPort)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexWithdrawalsDB)...)
	app.Action = run
	return app
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flusher()

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := db.Close(); cErr != nil {
			sugar.Errorf("failed to close database: err=%s", cErr.Error())
		}
	}()

	binanceDB, err := depositstorage.NewDB(sugar, db)
	if err != nil {
		return err
	}

	host := httputil.NewHTTPAddressFromContext(c)
	server, err := http.NewServer(host, binanceDB, sugar)
	if err != nil {
		return err
	}
	if err = server.Run(); err != nil {
		return err
	}
	return db.Close()
}

func main() {
	app := newServerCli()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
