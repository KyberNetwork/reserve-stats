package main

import (
	"log"
	"os"
	"sync"
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/crawler"
	"github.com/ethereum/go-ethereum"
	"github.com/urfave/cli"
)

const (
	addressesFlag = "addresses"

	maxWorkerFlag    = "max-workers"
	defaultMaxWorker = 4

	attemptsFlag    = "attempts"
	defaultAttempts = 3

	delayFlag        = "delay"
	defaultDelayTime = time.Minute

	durationFlag         = "duration"
	shardDurationFlag    = "shard-duration"
	defaultShardDuration = time.Hour * 24
)

func main() {
	app := libapp.NewApp()
	app.Name = "reserve-rates-crawler"
	app.Usage = "get the rates of all configured reserves at a certain block"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringSliceFlag{
			Name:   addressesFlag,
			EnvVar: "RESERVE_ADDRESSES",
			Usage:  "list of reserve contract addresses. Example: --addresses={\"0x1111\",\"0x222\"}",
		},
		cli.IntFlag{
			Name:   maxWorkerFlag,
			Usage:  "The maximum number of worker to fetch rates",
			EnvVar: "MAX_WORKERS",
			Value:  defaultMaxWorker,
		},
		cli.IntFlag{
			Name:   attemptsFlag,
			Usage:  "The number of attempt to query rates from blockchain",
			EnvVar: "ATTEMPTS",
			Value:  defaultAttempts,
		},
		cli.DurationFlag{
			Name:   delayFlag,
			Usage:  "The duration to put worker pools into sleep after each batch requets",
			EnvVar: "DELAY",
			Value:  defaultDelayTime,
		},
		cli.DurationFlag{
			Name:   durationFlag,
			Usage:  "The duration of a reserve rates before considered expired",
			EnvVar: "DURATION",
		},
		cli.DurationFlag{
			Name:   shardDurationFlag,
			Usage:  "The shard duration of a reserve rates",
			EnvVar: "SHARD_DURATION",
			Value:  defaultShardDuration,
		},
		blockchain.NewEthereumNodeFlags(),
	)
	app.Flags = append(app.Flags, timeutil.NewTimeRangeCliFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		err            error
		lastBlockErrCh = make(chan error)
		rateErrChn     = make(chan error)
		lastBlockBlCh  = make(chan int64)
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
		return err
	}

	toDate, err := timeutil.ToTimeFromContext(c)
	if err != nil {
		return err
	}

	attempts := c.Int(attemptsFlag)
	delayTime := c.Duration(delayFlag)

	addrs := c.StringSlice(addressesFlag)
	if len(addrs) == 0 {
		addr := contracts.InternalReserveAddress().MustGetOneFromContext(c)
		addrs = append(addrs, addr.Hex())
		sugar.Infow("using internal reserve address as user does not input any", "address", addr.Hex())
	}

	symbolResolver, err := blockchain.NewTokenSymbolFromContext(c)
	ratesCrawler, err := crawler.NewReserveRatesCrawler(sugar, addrs, ethClient, symbolResolver)
	lastBlockResolver := lastblockdaily.NewLastBlockResolver(ethClient, blockTimeResolver, fromDate, toDate, sugar)
	go lastBlockResolver.FetchLastBlock(lastBlockErrCh, lastBlockBlCh)

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
		case block := <-lastBlockBlCh:
			wg.Add(1)

			go func(errCh chan error, block int64, attempts int) {
				defer wg.Done()
				var toBreak bool
				for {
					rates, err := ratesCrawler.GetReserveRates(uint64(block))
					if err != nil {
						sugar.Debugw("failed to fetch rates", "error", err, "attempt left", attempts)
						if attempts == 0 {
							errCh <- err
							toBreak = true
						} else {
							attempts--
							time.Sleep(delayTime)
						}
					}
					if (err == nil) || (toBreak) {
						sugar.Debugw("rate result", "block", block, "rates", rates)
						break
					}
				}
			}(rateErrChn, block, attempts)
		}
	}
}
