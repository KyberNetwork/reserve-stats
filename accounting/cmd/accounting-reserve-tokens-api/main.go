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
	app.Name = "Accounting listed token api"
	app.Usage = ""
	app.Action = run

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingListedTokenPort)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}

	defer flusher()

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}
	listedTokenStorage, err := storage.NewDB(sugar, db, common.ListedTokenTable)
	if err != nil {
		return err
	}

	server := server.NewServer(sugar, httputil.NewHTTPAddressFromContext(c), listedTokenStorage)
	return server.Run()
}
