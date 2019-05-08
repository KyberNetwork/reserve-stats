package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/chainalysis"
	"github.com/KyberNetwork/reserve-stats/tradelogs/client"
)

const (
	tradeLogAccessKeyIDFlag     = "trade-log-access-key-id"
	tradeLogSecretAccessKeyFlag = "trade-log-secret-access-key"

	retryTimesFlag        = "retry-times"
	defaultRetryTimesFlag = 3

	delayTimeRetry = 5 * time.Second
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs chainalysis"
	app.Version = "0.0.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   tradeLogAccessKeyIDFlag,
			Usage:  "access key id for read trade log",
			EnvVar: "READ_TRADE_LOG_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   tradeLogSecretAccessKeyFlag,
			Usage:  "secret access key for read trade log",
			EnvVar: "READ_TRADE_LOG_SECRET_ACCESS_KEY",
		},
		cli.IntFlag{
			Name:   retryTimesFlag,
			Usage:  "number times retry when func get error",
			EnvVar: "RETRY_TIMES",
			Value:  defaultRetryTimesFlag,
		},
	}

	app.Flags = append(app.Flags, timeutil.NewMilliTimeRangeCliFlags()...)
	app.Flags = append(app.Flags, client.NewTradeLogCliFlags()...)
	app.Flags = append(app.Flags, chainalysis.NewChainAlysisCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		err error

		retryTimes = c.Int(retryTimesFlag)
	)
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	fromTime, err := timeutil.FromTimestampMillisFromContext(c)
	if err == timeutil.ErrEmptyFlag {
		sugar.Infof("no from time is provided, using default from time on trade log api")
	} else if err != nil {
		return err
	}

	toTime, err := timeutil.ToTimestampMillisFromContext(c)
	if err == timeutil.ErrEmptyFlag {
		sugar.Infof("no to time provided, using default to time on trade log api")
	} else if err != nil {
		return err
	}

	opt := client.WithAuth(c.String(tradeLogAccessKeyIDFlag), c.String(tradeLogSecretAccessKeyFlag))
	tradeLogCli, err := client.NewClientFromContext(sugar, c, opt)
	if err != nil {
		return err
	}

	tradeLogs, err := tradeLogCli.GetTradeLogs(fromTime, toTime)
	if err != nil {
		return err
	}

	chainAlysisCli, err := chainalysis.NewClientFromContext(sugar, c)
	if err != nil {
		return err
	}

	for i := 1; i <= retryTimes; i++ {
		err = chainAlysisCli.PushETHSentTransferEvent(tradeLogs)
		if err != nil {
			sugar.Infof("times retry: %d", i)
			time.Sleep(delayTimeRetry)
			continue
		}
		return nil
	}
	return err
}
