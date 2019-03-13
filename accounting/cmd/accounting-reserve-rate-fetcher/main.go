package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/fetcher"
	rrpostgres "github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/crawler"

	"github.com/KyberNetwork/tokenrate/coingecko"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

const (
	addressesFlag = "addresses"

	attemptsFlag    = "attempts"
	defaultAttempts = 3

	sleepTimeFlag    = "sleep-time"
	defaultSleepTime = 12 * time.Hour

	retryDelayFlag        = "retry-delay"
	defaultRetryDelayTime = 5 * time.Minute
	defaultPostGresDB     = common.DefaultDB
)

func main() {
	app := libapp.NewApp()
	app.Name = "accounting-reserve-rates-fetcher"
	app.Usage = "get the rates of all configured reserves at the last block of a day"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringSliceFlag{
			Name:   addressesFlag,
			EnvVar: "ADDRESSES",
			Usage:  "list of reserve contract addresses. Example: --addresses={\"0x1111\",\"0x222\"}",
		},
		cli.IntFlag{
			Name:   attemptsFlag,
			Usage:  "The number of attempt to query rates from blockchain",
			EnvVar: "ATTEMPTS",
			Value:  defaultAttempts,
		},
		cli.DurationFlag{
			Name:   retryDelayFlag,
			Usage:  "The duration to put worker pools into sleep after each batch request",
			EnvVar: "RETRY_DELAY",
			Value:  defaultRetryDelayTime,
		},
		cli.DurationFlag{
			Name:   sleepTimeFlag,
			Usage:  "The duration for the process to sleep after latest fetch in daemon mode",
			EnvVar: "SLEEP_TIME",
			Value:  defaultSleepTime,
		},
		blockchain.NewEthereumNodeFlags(),
	)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostGresDB)...)
	app.Flags = append(app.Flags, timeutil.NewTimeRangeCliFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		options        []fetcher.Option
		attempts       = c.Int(attemptsFlag)
		retryDelayTime = c.Duration(retryDelayFlag)
		sleepTime      = c.Duration(sleepTimeFlag)
	)

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	blockTimeResolver, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	if err != nil {
		return err
	}

	addrs := c.StringSlice(addressesFlag)
	if len(addrs) == 0 {
		return fmt.Errorf("empty reserve address")
	}

	for _, addr := range addrs {
		if !ethereum.IsHexAddress(addr) {
			return fmt.Errorf("malformed address: %s", addr)
		}
	}

	symbolResolver, err := blockchain.NewTokenInfoGetterFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot create symbol Resolver, err: %v", err)
	}

	ratesCrawler, err := crawler.NewReserveRatesCrawler(sugar, addrs, ethClient, symbolResolver)
	if err != nil {
		return fmt.Errorf("cannot rate crawler, err: %v", err)
	}

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	ratesStorage, err := rrpostgres.NewDB(sugar, db)
	if err != nil {
		return err
	}
	defer ratesStorage.Close()
	cgk := coingecko.New()

	lastBlockResolver := lastblockdaily.NewLastBlockResolver(ethClient, blockTimeResolver, sugar)

	fromDate, err := timeutil.FromTimeFromContext(c)
	if err != nil {
		if err != timeutil.ErrEmptyFlag {
			return err
		}
		sugar.Info("no fromDate provided. Checking permanent storage for the last stored rate")
	} else {
		options = append(options, fetcher.WithFromTime(fromDate))

	}

	toDate, err := timeutil.ToTimeFromContext(c)
	if err != nil {
		if err != timeutil.ErrEmptyFlag {
			return err
		}
		sugar.Info("no toDate provided. Running in daemon mode")
	} else {
		options = append(options, fetcher.WithToTime(toDate))
	}

	rrFetcher, err := fetcher.NewFetcher(sugar, ratesStorage, ratesCrawler, lastBlockResolver, cgk, retryDelayTime, sleepTime, attempts, options...)
	if err != nil {
		return err
	}
	return rrFetcher.Run()
}
