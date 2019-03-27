package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	fetcher "github.com/KyberNetwork/reserve-stats/accounting/binance/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/tradestorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

const (
	fromIDFlag        = "from-id"
	retryDelayFlag    = "retry-delay"
	attemptFlag       = "attempt"
	batchSizeFlag     = "batch-size"
	defaultRetryDelay = 2 // minute
	defaultAttempt    = 4
	defaultBatchSize  = 100
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting binance trades fetcher"
	app.Usage = "Fetch and store trades history from binance"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.IntFlag{
			Name:   retryDelayFlag,
			Usage:  "delay time when do a retry",
			EnvVar: "RETRY_DELAY",
			Value:  defaultRetryDelay,
		},
		cli.IntFlag{
			Name:   attemptFlag,
			Usage:  "number of time doing retry",
			EnvVar: "ATTEMPT",
			Value:  defaultAttempt,
		},
		cli.IntFlag{
			Name:   batchSizeFlag,
			Usage:  "batch to request to binance",
			EnvVar: "BATCH_SIZE",
			Value:  defaultBatchSize,
		},
		cli.Uint64Flag{
			Name:   fromIDFlag,
			Usage:  "id to get trade history from",
			EnvVar: "FROM_ID",
		},
	)

	app.Flags = append(app.Flags, binance.NewCliFlags()...)
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

	sugar.Info("initiate fetcher")

	binanceClient, err := binance.NewClientFromContext(c, sugar)
	if err != nil {
		return err
	}

	storage, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	binanceStorage, err := tradestorage.NewDB(sugar, storage)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := binanceStorage.Close(); cErr != nil {
			sugar.Errorf("Close database error", "error", cErr)
		}
	}()

	fromID := c.Uint64(fromIDFlag)
	if fromID == 0 {
		// get last stored ID in storage
		fromID, err = binanceStorage.GetLastStoredID()
		if err != nil {
			return err
		}
	}

	sugar.Infow("fetch trade from id", "id", fromID)

	retryDelay := c.Int(retryDelayFlag)
	attempt := c.Int(attemptFlag)
	batchSize := c.Int(batchSizeFlag)
	binanceFetcher := fetcher.NewFetcher(sugar, binanceClient, retryDelay, attempt, batchSize)

	tradeHistories, err := binanceFetcher.GetTradeHistory(fromID)
	if err != nil {
		return err
	}
	sugar.Debugw("trade histories", "result", tradeHistories)

	if err := binanceStorage.UpdateTradeHistory(tradeHistories); err != nil {
		return err
	}
	return binanceStorage.Close()
}
