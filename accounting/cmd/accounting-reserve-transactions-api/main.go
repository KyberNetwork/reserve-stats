package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/http"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting Reserve Transaction API"
	app.Usage = "Accounting Reserve Server for Transaction API calls"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultTransactionsDB)...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingTransactionsPort)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer libapp.NewFlusher(logger)()
	sugar := logger.Sugar()

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

	s := http.NewServer(logger, httputil.NewHTTPAddressFromContext(c), rts)

	if err = s.Run(); err != nil {
		return err
	}
	return db.Close()
}
