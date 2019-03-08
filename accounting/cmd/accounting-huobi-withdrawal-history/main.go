package main

import (
	"log"
	"os"
	"time"

	huobiFetcher "github.com/KyberNetwork/reserve-stats/accounting/huobi/fetcher"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
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

	huobiClient, err := huobi.NewClientFromContext(c, sugar)
	if err != nil {
		return err
	}
	fromID := c.Uint64(fromIDFlag)
	retryDelay := c.Duration(retryDelayFlag)
	maxAttempts := c.Int(maxAttemptFlag)

	fetcher := huobiFetcher.NewFetcher(sugar, huobiClient, retryDelay, maxAttempts)
	data, err := fetcher.GetWithdrawHistory(fromID)
	sugar.Debugw("fetched done", "error", err, "data", data)
	return nil
}
