package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/server"
	"github.com/KyberNetwork/reserve-stats/accounting/listed-tokens/storage"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

func main() {
	app := libapp.NewApp()
	app.Name = "accounting-reserve-tokens-api"
	app.Usage = "Accounting listed token api"
	app.Action = run

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingReserveTokensPort)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultListedTokenDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
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
	listedTokenStorage, err := storage.NewDB(sugar, db)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := listedTokenStorage.Close(); err != nil {
			sugar.Errorw("Close database error", "error", cErr)
		}
	}()

	s := server.NewServer(logger, httputil.NewHTTPAddressFromContext(c), listedTokenStorage)
	if err = s.Run(); err != nil {
		return err
	}
	return listedTokenStorage.Close()
}
