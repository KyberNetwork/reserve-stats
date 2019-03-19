package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	fetcher "github.com/KyberNetwork/reserve-stats/accounting/binance-fetcher"
	withdrawstorage "github.com/KyberNetwork/reserve-stats/accounting/binance-storage/withdraw-storage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	retryDelayFlag       = "retry-delay"
	attemptFlag          = "attempt"
	batchSizeFlag        = "batch-size"
	defaultRetryDelay    = 2 // minute
	defaultAttempt       = 4
	defaultBatchSize     = 100
	binanceWithdrawTable = "binance_withdraws"
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
	)

	app.Flags = append(app.Flags, binance.NewCliFlags()...)
	app.Flags = append(app.Flags, timeutil.NewMilliTimeRangeCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		fromTime, toTime time.Time
		err              error
	)

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

	fromTime, err = timeutil.FromTimeMillisFromContext(c)
	if err != nil {
		return err
	}

	toTime, err = timeutil.ToTimeMillisFromContext(c)
	if err != nil {
		return err
	}
	if toTime.IsZero() {
		toTime = time.Now()
	}

	retryDelay := c.Int(retryDelayFlag)
	attempt := c.Int(attemptFlag)
	batchSize := c.Int(batchSizeFlag)
	binanceFetcher := fetcher.NewFetcher(sugar, binanceClient, retryDelay, attempt, batchSize)

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	binanceStorage, err := withdrawstorage.NewDB(sugar, db, binanceWithdrawTable, "applyTime", "text")
	if err != nil {
		return err
	}

	withdrawHistory, err := binanceFetcher.GetWithdrawHistory(fromTime, toTime)
	if err != nil {
		return err
	}

	return binanceStorage.UpdateWithdrawHistory(withdrawHistory)
}
