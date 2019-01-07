package main

import (
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/tokenrate/coingecko"

	"github.com/urfave/cli"
)

const (
	kyberNetworkTokenID = "kyber-network"
	usdCurrencyID       = "usd"
	dbName              = "token_rate"
)

func main() {
	app := libapp.NewApp()
	app.Name = "KNC USD rate Fetcher"
	app.Usage = "Fetch KNC-USD Rate from provider"
	app.Version = "0.0.1"
	app.Action = run
	app.Flags = append(app.Flags, timeutil.NewTimeRangeCliFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	cgk := coingecko.New()

	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	tokenRate, err := tokenrate.NewRateFetcher(sugar, dbName, influxClient, cgk)
	if err != nil {
		return err
	}

	from, err := timeutil.MustGetFromTimeFromContext(c)
	if err != nil {
		return err
	}
	to := timeutil.GetToTimeFromContext(c)

	return tokenRate.FetchRatesInRanges(from, to, kyberNetworkTokenID, usdCurrencyID)
}
