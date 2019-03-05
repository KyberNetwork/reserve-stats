package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	fetcher "github.com/KyberNetwork/reserve-stats/accounting/binance-fetcher"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	fromFlag = "from"
	toFlag   = "to"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting binance trades fetcher"
	app.Usage = "Fetch and store trades history from binance"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.Uint64Flag{
			Name:   fromFlag,
			Usage:  "From timestamp(millisecond) to get trade history from",
			EnvVar: "FROM",
		},
		cli.Uint64Flag{
			Name:   toFlag,
			Usage:  "To timestamp(millisecond) to get trade history to",
			EnvVar: "TO",
		},
	)

	app.Flags = append(app.Flags, binance.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		fromTime, toTime time.Time
	)

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}

	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("initiate fetcher")

	binanceClient, err := binance.NewClientFromContext(c, sugar)
	if err != nil {
		return err
	}

	if c.Uint64(fromFlag) != 0 {
		fromTime = timeutil.TimestampMsToTime(c.Uint64(fromFlag))
	}

	if c.Uint64(toFlag) != 0 {
		toTime = timeutil.TimestampMsToTime(c.Uint64(toFlag))
	}

	binanceFetcher := fetcher.NewFetcher(sugar, binanceClient)

	err = binanceFetcher.GetTradeHistory(fromTime, toTime)
	if err != nil {
		return err
	}

	return nil
}
