package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage/postgres"
	"github.com/KyberNetwork/reserve-stats/accounting/wallet-erc20/http"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting Wallet Erc20 Transaction API"
	app.Usage = "Accounting Server for Wallet ERC20 Transaction API calls"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultTransactionsDB)...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingWalletErc20Port)...)
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

	rts, err := postgres.NewStorage(sugar, db)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := db.Close(); cErr != nil {
			sugar.Errorf("failed to close database: err=%s", cErr.Error())
		}
	}()

	s := http.NewServer(sugar, httputil.NewHTTPAddressFromContext(c), rts)

	if err = s.Run(); err != nil {
		return err
	}
	return db.Close()
}
