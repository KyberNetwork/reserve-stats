package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/http"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/storage/postgresql"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
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
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AccountingReserveAddressPort)...)
	app.Flags = append(app.Flags, etherscan.NewCliFlags()...)
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
	defer logger.Sync()

	sugar := logger.Sugar()

	etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
	if err != nil {
		return err
	}

	resolv := blockchain.NewEtherscanContractTimestampResolver(sugar, etherscanClient)

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	st, err := postgresql.NewStorage(sugar, db, resolv)
	if err != nil {
		return err
	}

	s := http.NewServer(sugar, httputil.NewHTTPAddressFromContext(c), st)

	return s.Run()
}
