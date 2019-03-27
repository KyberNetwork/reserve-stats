package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KyberNetwork/tokenrate/coingecko"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/client"
	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/fetcher"
	rrpostgres "github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	lbdpostgres "github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/storage/postgres"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/crawler"
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
	app.Flags = append(app.Flags, etherscan.NewCliFlags()...)
	app.Flags = append(app.Flags, timeutil.NewTimeRangeCliFlags()...)
	app.Flags = append(app.Flags, client.NewClientFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		attempts       = c.Int(attemptsFlag)
		retryDelayTime = c.Duration(retryDelayFlag)
		sleepTime      = c.Duration(sleepTimeFlag)
		addressClient  client.Interface
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

	symbolResolver, err := blockchain.NewTokenInfoGetterFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot create symbol Resolver, err: %v", err)
	}

	ratesCrawler, err := crawler.NewReserveRatesCrawler(sugar, ethClient, symbolResolver)
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
	defer func() {
		if cErr := ratesStorage.Close(); cErr != nil {
			sugar.Errorf("failed to close rate storage: err=%s", cErr.Error())
		}
	}()
	cgk := coingecko.New()

	lbdDB, err := lbdpostgres.NewDB(sugar, db)
	if err != nil {
		return err
	}
	lastBlockResolver := lastblockdaily.NewLastBlockResolver(ethClient, blockTimeResolver, sugar, lbdDB)

	addrs := c.StringSlice(addressesFlag)
	if len(addrs) != 0 {
		sugar.Infow("using provided addresses instead of querying from accounting-reserve-addresses service")
		etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
		if err != nil {
			return err
		}
		resolver := blockchain.NewEtherscanContractTimestampResolver(sugar, etherscanClient)
		addressClient, err = client.NewFixedAddresses(addrs, resolver)
		if err != nil {
			return err
		}
	} else {
		addressClient, err = client.NewClientFromContext(c, sugar)
		if err != nil {
			return err
		}
	}

	rrFetcher, err := fetcher.NewFetcher(sugar, ratesStorage, ratesCrawler, lastBlockResolver, cgk, retryDelayTime, sleepTime, attempts, addressClient)
	if err != nil {
		return err
	}

	fromDate, err := timeutil.FromTimeFromContext(c)
	switch err {
	case timeutil.ErrEmptyFlag:
		sugar.Info("fromDate not provide. Fetcher running in daemon mode...")
	case nil:
		sugar.Infof("fromDate provided. Fetcher run from %s ...", fromDate.String())
	default:
		return err
	}

	toDate, err := timeutil.FromTimeFromContext(c)
	switch err {
	case timeutil.ErrEmptyFlag:
		sugar.Info("toDate not provide. Fetcher running till now...")
		toDate = time.Now()
	case nil:
		sugar.Infof("toDate provided. Fetcher run to %s ...", toDate.String())
	default:
		return err
	}

	if !fromDate.IsZero() && !toDate.IsZero() {
		var ethAddrs []ethereum.Address
		addrs, err := addressClient.GetAllReserveAddress()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			ethAddrs = append(ethAddrs, addr.Address)
		}
		if err = rrFetcher.Fetch(fromDate, toDate, ethAddrs); err != nil {
			return err
		}
		return ratesStorage.Close()
	}

	if err = rrFetcher.Run(); err != nil {
		return err
	}
	return ratesStorage.Close()
}
