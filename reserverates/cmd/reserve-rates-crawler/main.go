package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	influxRateStorage "github.com/KyberNetwork/reserve-stats/reserverates/storage/influx"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage/postgres"
	"github.com/KyberNetwork/reserve-stats/reserverates/workers"
)

const (
	addressesFlag = "addresses"

	fromBlockFlag = "from-block"
	toBlockFlag   = "to-block"

	maxWorkerFlag    = "max-workers"
	defaultMaxWorker = 4

	attemptsFlag    = "attempts"
	defaultAttempts = 3

	delayFlag        = "delay"
	defaultDelayTime = time.Minute

	durationFlag         = "duration"
	shardDurationFlag    = "shard-duration"
	defaultShardDuration = time.Hour * 24

	dbEngineFlag      = "db-engine"
	defaultPostgresDB = "reserve_rates"
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
		cli.StringFlag{
			Name:   fromBlockFlag,
			Usage:  "Fetch rates from block",
			EnvVar: "FROM_BLOCK",
		},
		cli.StringFlag{
			Name:   toBlockFlag,
			Usage:  "Fetch rates to block",
			EnvVar: "TO_BLOCK",
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
		cli.StringFlag{
			Name:   dbEngineFlag,
			Usage:  "db engine flag",
			EnvVar: "DB_ENGINE",
			Value:  "postgres",
		},
		blockchain.NewEthereumNodeFlags(),
	)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultPostgresDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func mustNewClient(url string) *ethclient.Client {
	cc := &http.Client{Transport: roundTripperExt{c: &http.Client{}}}
	r, err := rpc.DialHTTPWithClient(url, cc)
	if err != nil {
		zap.S().Panicw("init custom ethclient failed", "err", err)
	}
	client := ethclient.NewClient(r)
	return client
}

func run(c *cli.Context) error {
	var (
		err                error
		fromBlock, toBlock *big.Int
		daemon             bool
	)

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	nodeURL := blockchain.NodeURLFromFlag(c)
	ethClient := mustNewClient(nodeURL)

	blockTimeResolver, err := blockchain.NewBlockTimeResolver(sugar, ethClient)
	if err != nil {
		return err
	}

	var rateStorage storage.ReserveRatesStorage
	if c.String(dbEngineFlag) == "influxdb" {
		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return err
		}
		var options []influxRateStorage.RateStorageOption
		duration := c.Duration(durationFlag)
		shardDuration := c.Duration(shardDurationFlag)
		if duration != 0 && shardDuration != 0 {
			options = append(options, influxRateStorage.RateStorageOptionWithRetentionPolicy(duration, shardDuration))
		}

		if rateStorage, err = influxRateStorage.NewRateInfluxDBStorage(sugar, influxClient, common.DatabaseName, blockTimeResolver, options...); err != nil {
			return err
		}
	} else {
		db, err := libapp.NewDBFromContext(c)
		if err != nil {
			return err
		}
		if rateStorage, err = postgres.NewPostgresStorage(db, sugar, blockTimeResolver); err != nil {
			return err
		}
	}

	if c.String(fromBlockFlag) == "" {
		sugar.Info("no from block flag provided, checking last stored block")
	} else {
		fromBlock, err = libapp.ParseBigIntFlag(c, fromBlockFlag)
		if err != nil {
			return err
		}
	}

	if c.String(toBlockFlag) == "" {
		daemon = true
		sugar.Info("running in daemon mode")
	} else {
		toBlock, err = libapp.ParseBigIntFlag(c, toBlockFlag)
		if err != nil {
			return err
		}
	}

	maxWorkers := c.Int(maxWorkerFlag)
	attempts := c.Int(attemptsFlag)
	delayTime := c.Duration(delayFlag)

	addrs := c.StringSlice(addressesFlag)
	if len(addrs) == 0 {
		addr := contracts.InternalReserveAddress().MustGetOneFromContext(c)
		addrs = append(addrs, addr.Hex())
		sugar.Infow("using internal reserve address as user does not input any", "address", addr.Hex())
	}

	var ethAddrs []ethereum.Address
	for _, addr := range addrs {
		if !ethereum.IsHexAddress(addr) {
			return fmt.Errorf("non etherum address input %s", addr)
		}
		ethAddrs = append(ethAddrs, ethereum.HexToAddress(addr))
	}
	for {
		currentHeader, fErr := ethClient.HeaderByNumber(context.Background(), nil)
		if fErr != nil {
			return fErr
		}

		if fromBlock == nil {
			lastBlock, fErr := rateStorage.LastBlock()
			if fErr != nil {
				return fErr
			}

			if lastBlock != 0 {
				fromBlock = big.NewInt(0).SetInt64(lastBlock)
				sugar.Infow("using last stored block number", "from_block", fromBlock)
			} else {
				sugar.Infow("using latest known block", "from_block", currentHeader.Number.String())
				fromBlock = currentHeader.Number
			}
		}

		if toBlock == nil {
			toBlock = big.NewInt(0).Add(currentHeader.Number, big.NewInt(1))
			sugar.Infow("fetching reserve rates up to latest known block number", "to_block", toBlock.String())
		}

		pool := workers.NewPool(sugar, maxWorkers, rateStorage)
		doneCh := make(chan struct{})

		go func(fromBlock, toBlock int64) {
			var jobOrder = pool.GetLastCompleteJobOrder()

			for block := fromBlock; block < toBlock; block++ {
				jobOrder++
				pool.Run(workers.NewFetcherJob(c, jobOrder, uint64(block), ethAddrs, attempts))
			}

			for pool.GetLastCompleteJobOrder() < jobOrder {
				time.Sleep(time.Second)
			}

			doneCh <- struct{}{}
		}(fromBlock.Int64(), toBlock.Int64())

		for {
			var toBreak = false

			select {
			case <-doneCh:
				sugar.Info("all jobs are successfully executed, waiting for the workers pool to shut down")
				pool.Shutdown()
			case fErr := <-pool.ErrCh():
				if fErr != nil {
					sugar.Errorw("job failed to execute", "error", fErr)
					log.Fatal(fErr)
				} else {
					sugar.Info("workers pool is successfully shut down")
					toBreak = true
				}
			}

			if toBreak {
				break
			}

		}

		if daemon {
			sugar.Infow("waiting before fetching new rates",
				"last_from_block", fromBlock.String(),
				"last_to_block", toBlock.String(),
				"sleep", delayTime.String())
			fromBlock, toBlock = nil, nil
			time.Sleep(delayTime)
		} else {
			break
		}
	}

	return nil
}
