package main

import (
	"log"
	"os"

	huobiFetcher "github.com/KyberNetwork/reserve-stats/accounting/accouting-huobi-fetcher"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Huobi Fetcher"
	app.Usage = "Huobi Fetcher  Reserve Addresses Manager"
	app.Action = run
	app.Version = "0.0.1"
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
	from, err := timeutil.FromTimeFromContext(c)
	if err != nil {
		return err
	}
	to, err := timeutil.ToTimeFromContext(c)
	if err != nil {
		return err
	}

	fetcher := huobiFetcher.NewFetcher(sugar, huobiClient)
	data, err := fetcher.GetTradeHistory(from, to)
	sugar.Debugw("fetched done", "error", err, "data", data)
	return nil
}
