package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	huobiFetcher "github.com/KyberNetwork/reserve-stats/accounting/huobi/fetcher"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	retryDelayFlag    = "retry-delay"
	maxAttemptFlag    = "max-attempts"
	defaultMaxAttempt = 3
	defaultRetryDelay = time.Second
)

func main() {
	app := libapp.NewApp()
	app.Name = "Huobi Fetcher"
	app.Usage = "Huobi Fetcher for trade logs"
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
	)
	app.Flags = append(app.Flags, huobi.NewCliFlags()...)
	app.Flags = append(app.Flags, timeutil.NewMilliTimeRangeCliFlags()...)
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
	from, err := timeutil.FromTimeMillisFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot get from time: %v", err)
	}
	to, err := timeutil.ToTimeMillisFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot get to time: %v", err)
	}
	retryDelay := c.Duration(retryDelayFlag)
	maxAttempts := c.Int(maxAttemptFlag)

	fetcher := huobiFetcher.NewFetcher(sugar, huobiClient, retryDelay, maxAttempts)
	data, err := fetcher.GetTradeHistory(from, to)
	sugar.Debugw("fetched done", "error", err, "data", data)
	return nil
}
