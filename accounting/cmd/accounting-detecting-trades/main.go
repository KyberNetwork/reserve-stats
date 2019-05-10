package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/storage/postgres"
	tradedetect "github.com/KyberNetwork/reserve-stats/accounting/trades-detect"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting Trade Detection Service"
	app.Usage = "Accounting Trade Detection Service"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultTransactionsDB)...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		result = make(map[int64]bool)
	)
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
	txs, err := rts.GetERC20WithTradeNull()
	if err != nil {
		return err
	}
	sugar.Debugw("number of tx with is trade is null", "length", len(txs))

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}
	for id, tx := range txs {
		isTrade, err := tradedetect.DetectTradeTransaction(tx.Hash, ethClient)
		if err != nil {
			return err
		}
		sugar.Debugw("trade detect", "tx hash", tx.Hash, "is trade", isTrade)
		result[id] = isTrade
	}
	return rts.UpdateERC20IsTrade(result)
}
