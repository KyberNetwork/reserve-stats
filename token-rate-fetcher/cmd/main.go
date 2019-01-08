package main

import (
	"os"
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	tokenrate "github.com/KyberNetwork/reserve-stats/token-rate-fetcher"
	"github.com/KyberNetwork/reserve-stats/token-rate-fetcher/storage"
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

	influxStorage, err := storage.NewInfluxStorage(influxClient, dbName, sugar)
	if err != nil {
		return err
	}

	tokenRate, err := tokenrate.NewRateFetcher(sugar, influxStorage, cgk)
	if err != nil {
		return err
	}

	from, err := timeutil.MustGetFromTimeFromContext(c)
	if err != nil {
		sugar.Info("No from time is provided, seeking for the first data point in DB...")
		from, err = influxStorage.GetFirstTimePoint(cgk.Name(), kyberNetworkTokenID, usdCurrencyID)
		if err != nil {
			return err
		}
	}
	to, err := timeutil.GetToTimeFromContextWithDeamon(c)
	if err == timeutil.ErrEmptyFlag {
		sugar.Info("No to time is provide, running in daemon mode...")
		for {
			to = time.Now()
			if err := tokenRate.FetchRatesInRanges(from, to, kyberNetworkTokenID, usdCurrencyID); err != nil {
				return err
			}
			from = to
			time.Sleep(12 * time.Hour)
		}
	}
	if err != nil {
		return err
	}

	return tokenRate.FetchRatesInRanges(from, to, kyberNetworkTokenID, usdCurrencyID)
}
