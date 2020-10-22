package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	huobiFetcher "github.com/KyberNetwork/reserve-stats/accounting/huobi/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"

	"github.com/urfave/cli"
)

const (
	retryDelayFlag    = "retry-delay"
	maxAttemptFlag    = "max-attempts"
	defaultMaxAttempt = 3
	defaultRetryDelay = time.Second
	fromIDFlag        = "from-id"
	defaultFromID     = 0
)

func main() {
	app := libapp.NewApp()
	app.Name = "Huobi Fetcher"
	app.Usage = "Huobi Fetcher for withdrawal records. It will fetch from input ID to the latest withdrawal"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags,
		cli.IntFlag{
			Name:   maxAttemptFlag,
			Usage:  "The maximum number of attempts to retry fetching data",
			EnvVar: "MAX_ATTEMPTS",
			Value:  defaultMaxAttempt,
		},
		cli.DurationFlag{
			Name:   retryDelayFlag,
			Usage:  "The duration to put fetcher job to sleep after each fail attempt",
			EnvVar: "RETRY_DELAY",
			Value:  defaultRetryDelay,
		},
		cli.Uint64Flag{
			Name:   fromIDFlag,
			Usage:  "The ID from which to query withdrawal history from. Default is 0 (fetch from beginning)",
			EnvVar: "FROM_ID",
			Value:  defaultFromID,
		},
	)
	app.Flags = append(app.Flags, huobi.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexWithdrawalsDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}

	defer flusher()

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot create db from flags: %v", err)
	}

	hdb, err := postgres.NewDB(sugar, db)
	if err != nil {
		return fmt.Errorf("cannot create huobi database instance: %v", err)

	}

	defer func() {
		if cErr := hdb.Close(); cErr != nil {
			sugar.Errorf("Close database error", "error", cErr)
		}
	}()

	fromID := c.Uint64(fromIDFlag)
	sugar.Infow("get withdraw history from", "ID", fromID+1)

	retryDelay := c.Duration(retryDelayFlag)
	maxAttempts := c.Int(maxAttemptFlag)
	options, err := huobi.ClientOptionFromContext(c)
	if err != nil {
		return err
	}
	accounts, err := huobi.AccountsFromContext(c)
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if fromID == 0 {
			sugar.Info("From id is not provided, get last id stored in db")
			fromID, err = hdb.GetLastIDStored(account.Name)
			if err != nil {
				return err
			}
		}
		huobiClient, err := huobi.NewClient(account.APIKey, account.SecretKey, sugar, options...)
		if err != nil {
			return err
		}
		fetcher := huobiFetcher.NewFetcher(sugar, huobiClient, retryDelay, maxAttempts)
		data, err := fetcher.GetWithdrawHistory(fromID + 1)
		if err != nil {
			return err
		}
		var records []huobi.WithdrawHistory

		for _, record := range data {
			records = append(records, record...)
		}

		if err = hdb.UpdateWithdrawHistory(records, account.Name); err != nil {
			return err
		}
	}
	return hdb.Close()
}
