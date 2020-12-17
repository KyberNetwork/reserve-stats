package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/fetcher"
	depositstorage "github.com/KyberNetwork/reserve-stats/accounting/binance/storage/depositstorage"
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
	app.Name = "Accounting binance deposit fetcher"
	app.Usage = "Fetch and store deposit history from binance"
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
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexDepositsDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		fromTime, toTime time.Time
		err              error
		accounts         []common.Account
		errGroup         errgroup.Group
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

	binanceStorage, err := depositstorage.NewDB(sugar, db)
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
	fromTime, err = timeutil.ToTimeMillisFromContext(c)
	if err != nil {
		return err
	}

	retryDelay := c.Duration(retryDelayFlag)
	attempt := c.Int(attemptFlag)
	batchSize := c.Int(batchSizeFlag)

	accounts, err = binance.AccountsFromContext(c)
	if err != nil {
		return err
	}
	//
	assets, err := getAssetsList(accounts, sugar)
	if err != nil {
		return err
	}
	for _, account := range accounts {

		binanceClient, err := binance.NewBinance(account.APIKey, account.SecretKey, sugar)
		if err != nil {
			return err
		}

		binanceFetcher := fetcher.NewFetcher(sugar, binanceClient, retryDelay, attempt, batchSize, nil, "", nil, binanceStorage)

		errGroup.Go(
			func(accountName string) func() error {
				return func() error {
					_, err = binanceFetcher.GetDepositHistory(assets, fromTime, toTime, account.Name)
					return err
				}
			}(account.Name),
		)
	}
	if err := errGroup.Wait(); err != nil {
		return err
	}
	return binanceStorage.Close()
}

func getAssetsList(accounts []common.Account, sugar *zap.SugaredLogger) ([]string, error) {
	var (
		assets []string
		err    error
	)
	// init a binance client
	binanceClient, err := binance.NewBinance(accounts[0].APIKey, accounts[0].SecretKey, sugar)
	if err != nil {
		return nil, err
	}

	// get account info
	accountInfo, err := binanceClient.GetAccountInfo()
	if err != nil {
		sugar.Errorw("failed dto get assets list", "error", err)
		return assets, err
	}

	for _, balance := range accountInfo.Balances {
		assets = append(assets, balance.Asset)
	}
	return assets, nil
}
