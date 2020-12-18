package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	fetcher "github.com/KyberNetwork/reserve-stats/accounting/binance/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/tradestorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

const (
	retryDelayFlag    = "retry-delay"
	attemptFlag       = "attempt"
	batchSizeFlag     = "batch-size"
	symbolsFlag       = "symbols"
	defaultRetryDelay = 2 * time.Minute
	defaultAttempt    = 4
	defaultBatchSize  = 20
)

var sugar *zap.SugaredLogger

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
		cli.StringSliceFlag{
			Name:   symbolsFlag,
			Usage:  "symbol to get trade history for, if not provide then get from binance",
			EnvVar: "SYMBOLS",
		},
	)

	app.Flags = append(app.Flags, binance.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexTradesDB)...)

	if err := app.Run(os.Args); err != nil {
		sugar.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		flusher  func()
		err      error
		accounts []common.Account
		errGroup errgroup.Group
	)
	sugar, flusher, err = libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}

	defer flusher()

	sugar.Info("initiate fetcher")

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

	binanceClient, err := binance.NewBinance("", "", sugar) // this is public client to get exchange info
	if err != nil {
		return err
	}

	var tokenPairs []binance.Symbol
	exchangeInfo, err := binanceClient.GetExchangeInfo()
	if err != nil {
		return err
	}
	tokenPairs = exchangeInfo.Symbols

	retryDelay := c.Duration(retryDelayFlag)
	attempt := c.Int(attemptFlag)
	batchSize := c.Int(batchSizeFlag)
	options, err := binance.ClientOptionFromContext(c)
	if err != nil {
		return err
	}
	accounts, err = binance.AccountsFromContext(c)
	if err != nil {
		return err
	}
	// notETHTrades := make(map[*binance.Symbol][]binance.TradeHistory)
	for _, account := range accounts {
		fromIDs := make(map[string]uint64)
		for _, pair := range tokenPairs {
			sugar.Info("from id is not provided, get latest from id stored in database")
			from, err := binanceStorage.GetLastStoredID(pair.Symbol, account.Name)
			if err != nil {
				return err
			}
			fromIDs[pair.Symbol] = from
		}

		binanceClient, err := binance.NewBinance(account.APIKey, account.SecretKey, sugar, options...)
		if err != nil {
			return err
		}

		binanceFetcher := fetcher.NewFetcher(sugar, binanceClient, retryDelay, attempt, batchSize, binanceStorage, account.Name, nil)
		errGroup.Go(
			func(accountName string) func() error {
				return func() error {
					return binanceFetcher.GetTradeHistory(fromIDs, tokenPairs, accountName)
				}
			}(account.Name))
	}
	if err := errGroup.Wait(); err != nil {
		return err
	}
	return nil
}
