package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	http "github.com/KyberNetwork/reserve-stats/accounting/cex-withdrawal-api"
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
	app.Name = "cex-withdrawal-api"
	app.Usage = "server for query accounting cex withdrawal"
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.ReserveRatesPort)...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostGresDB)...)
	app.Action = run
	return app
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	sugar := logger.Sugar()
	defer logger.Sync()

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	huobiDb, err := huobiPostgres.NewDB(sugar, db)
	if err != nil {
		return err
	}
	defer func(err *error) {
		var cErr error
		cErr = huobiDb.Close()
		if err == nil {
			err = &cErr
		} else {
			sugar.Error("DB closing failed", "error", cErr)
		}

	}(&err)
	hostStr := httputil.NewHTTPAddressFromContext(c)
	server, err := http.NewServer(hostStr, huobiDb, nil, sugar)
	if err != nil {
		return err
	}
	err = server.Run()
	return err
}

func main() {
	app := newServerCli()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
