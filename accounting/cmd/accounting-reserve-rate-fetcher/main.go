package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	rrpostgres "github.com/KyberNetwork/reserve-stats/accounting/reserve-rate-storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	rsvRateCommon "github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/crawler"

	"github.com/KyberNetwork/tokenrate"
	"github.com/KyberNetwork/tokenrate/coingecko"
	"github.com/ethereum/go-ethereum"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	addressesFlag = "addresses"

	attemptsFlag    = "attempts"
	defaultAttempts = 3

	retryDelayFlag        = "retry-delay"
	defaultRetryDelayTime = 5 * time.Minute
	defaultPostGresDB     = "reserve_stats"
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
		blockchain.NewEthereumNodeFlags(),
	)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostGresDB)...)
	app.Flags = append(app.Flags, timeutil.NewTimeRangeCliFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func retryFetchTokenRate(maxAttempt int,
	sugar *zap.SugaredLogger,
	rsvRateCrawler *crawler.ReserveRatesCrawler,
	block uint64,
	retryInterval time.Duration) (map[string]map[string]rsvRateCommon.ReserveRateEntry, error) {
	var (
		result map[string]map[string]rsvRateCommon.ReserveRateEntry
		err    error
		logger = sugar.With("function", "main/retryFetchTokenRate", "block", block)
	)

	for i := 0; i < maxAttempt; i++ {
		result, err = rsvRateCrawler.GetReserveRates(block)
		if err == nil {
			return result, nil
		}
		logger.Debugw("failed to fetch reserve rate", "attempt", i, "error", err)
		time.Sleep(retryInterval)
	}

	return nil, err
}

func retryFetchETHUSDRate(maxAttempt int,
	sugar *zap.SugaredLogger,
	fetcher tokenrate.ETHUSDRateProvider,
	timestamp time.Time,
	retryInterval time.Duration) (float64, error) {
	var (
		result float64
		err    error
		logger = sugar.With("function", "main/retryFetchTokenRate", "time", timestamp.String())
	)

	for i := 0; i < maxAttempt; i++ {
		result, err = fetcher.USDRate(timestamp)
		if err == nil {
			return result, nil
		}
		logger.Debugw("failed to fetch ETH-USD rate", "attempt", i, "error", err)
		time.Sleep(retryInterval)
	}

	return result, err
}
func run(c *cli.Context) error {
	var (
		err            error
		lastBlockErrCh = make(chan error)
		rateErrChn     = make(chan error)
		lastBlockBlCh  = make(chan lastblockdaily.BlockInfo)
		wg             = &sync.WaitGroup{}
	)

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	blockTimeResolver, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	if err != nil {
		return err
	}

	fromDate, err := timeutil.FromTimeFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot get from date, err: %v", err)
	}

	toDate, err := timeutil.ToTimeFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot get to date, err: %v", err)
	}

	attempts := c.Int(attemptsFlag)
	retryDelayTime := c.Duration(retryDelayFlag)

	addrs := c.StringSlice(addressesFlag)
	//TODO: remove this for production server
	if len(addrs) == 0 {
		addr := contracts.InternalReserveAddress().MustGetOneFromContext(c)
		addrs = append(addrs, addr.Hex())
		sugar.Infow("using internal reserve address as user does not input any", "address", addr.Hex())
	}

	symbolResolver, err := blockchain.NewTokenSymbolFromContext(c)
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

	lastBlockResolver := lastblockdaily.NewLastBlockResolver(ethClient, blockTimeResolver, fromDate, toDate, sugar)
	go lastBlockResolver.Run(lastBlockBlCh, lastBlockErrCh)

	for {
		select {
		case err := <-lastBlockErrCh:
			if err == ethereum.NotFound {
				sugar.Info("reached the end date")
				wg.Wait()
				sugar.Info("all fetcher jobs are completed")
				return nil
			}
			return err
		case err := <-rateErrChn:
			if err != nil {
				return err
			}
		case blockInfo := <-lastBlockBlCh:
			wg.Add(1)

			go func(errCh chan error, blockInfo lastblockdaily.BlockInfo, attempts int) {
				defer wg.Done()
				rates, rateErr := retryFetchTokenRate(attempts, sugar, ratesCrawler, blockInfo.Block, retryDelayTime)
				if rateErr != nil {
					errCh <- rateErr
				}
				//TODO: parallel this
				ethUSDRate, err := retryFetchETHUSDRate(attempts, sugar, cgk, blockInfo.Timestamp, retryDelayTime)
				if err != nil {
					errCh <- err
				}
				if err = ratesStorage.UpdateRatesRecords(blockInfo, rates); err != nil {
					errCh <- err
				}
				if err = ratesStorage.UpdateETHUSDPrice(blockInfo, ethUSDRate); err != nil {
					errCh <- err
				}
			}(rateErrChn, blockInfo, attempts)
		}
	}
}
