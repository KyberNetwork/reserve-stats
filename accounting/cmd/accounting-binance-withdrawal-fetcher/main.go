package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/fetcher"
	withdrawstorage "github.com/KyberNetwork/reserve-stats/accounting/binance/storage/withdrawalstorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	retryDelayFlag    = "retry-delay"
	attemptFlag       = "attempt"
	batchSizeFlag     = "batch-size"
	defaultRetryDelay = 2 * time.Minute
	defaultAttempt    = 4
	defaultBatchSize  = 100
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting binance trades fetcher"
	app.Usage = "Fetch and store trades history from binance"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.DurationFlag{
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
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexWithdrawalsDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		fromTime, toTime time.Time
		err              error
		accounts         []common.Account
	)

	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}

	defer flusher()

	sugar.Info("initiate fetcher")

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	binanceStorage, err := withdrawstorage.NewDB(sugar, db)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := binanceStorage.Close(); cErr != nil {
			sugar.Errorf("Close database error", "error", cErr)
		}
	}()

	sugar.Infow("from timestamp", "fromTime", fromTime)

	toTime, err = timeutil.ToTimeMillisFromContext(c)
	if err != nil {
		return err
	}
	if toTime.IsZero() {
		toTime = time.Now()
	}

	retryDelay := c.Duration(retryDelayFlag)
	attempt := c.Int(attemptFlag)
	batchSize := c.Int(batchSizeFlag)

	accounts, err = binance.AccountsFromContext(c)
	if err != nil {
		return err
	}
	for _, account := range accounts {

		fromTime, err = timeutil.FromTimeMillisFromContext(c)
		if err != nil {
			return err
		}
		if fromTime.IsZero() {
			sugar.Info("from time is not provided, get latest timestamp from database")
			fromTime, err = binanceStorage.GetLastStoredTimestamp(account.Name)
			if err != nil {
				return err
			}
		}

		binanceClient, err := binance.NewBinance(account.APIKey, account.SecretKey, sugar)
		if err != nil {
			return err
		}

		binanceFetcher := fetcher.NewFetcher(sugar, binanceClient, retryDelay, attempt, batchSize, nil, "", nil)

		withdrawHistory, err := binanceFetcher.GetWithdrawHistory(fromTime, toTime)
		if err != nil {
			return err
		}

		if err := binanceStorage.UpdateWithdrawHistory(withdrawHistory, account.Name); err != nil {
			return err
		}
	}
	return binanceStorage.Close()
}
