package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/withdrawalstorage"
	"github.com/KyberNetwork/reserve-stats/accounting/cex-withdrawal/http"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	huobiPostgres "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

const (
	defaultPostGresDB = common.DefaultDB
)

func newServerCli() *cli.App {
	app := libapp.NewApp()
	app.Name = "cex-trade-withdrawal-api"
	app.Usage = "server for query accounting cex-trade withdrawal"
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingCEXWithdrawalsPort)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostGresDB)...)
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

	huobiDB, err := huobiPostgres.NewDB(sugar, db)
	if err != nil {
		return err
	}

	binanceDB, err := withdrawalstorage.NewDB(sugar, db)
	if err != nil {
		return err
	}

	host := httputil.NewHTTPAddressFromContext(c)
	server, err := http.NewServer(host, huobiDB, binanceDB, sugar)
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
