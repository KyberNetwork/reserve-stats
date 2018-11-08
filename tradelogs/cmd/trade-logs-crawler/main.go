package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/workers"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	fromBlockFlag    = "from-block"
	defaultFromblock = 5069586
	toBlockFlag      = "to-block"

	maxWorkerFlag    = "max-workers"
	defaultMaxWorker = 5

	maxBlocksFlag    = "max-blocks"
	defaultMaxBlocks = 100

	attemptsFlag    = "attempts"
	defaultAttempts = 5

	delayFlag        = "delay"
	defaultDelayTime = time.Minute
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs Fetcher"
	app.Usage = "Fetch trade logs on KyberNetwork"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   fromBlockFlag,
			Usage:  "Fetch trade logs from block",
			EnvVar: "FROM_BLOCK",
		},
		cli.StringFlag{
			Name:   toBlockFlag,
			Usage:  "Fetch trade logs to block",
			EnvVar: "TO_BLOCK",
		},
		cli.IntFlag{
			Name:   maxWorkerFlag,
			Usage:  "The maximum number of worker to fetch trade logs",
			EnvVar: "MAX_WORKER",
			Value:  defaultMaxWorker,
		},
		cli.IntFlag{
			Name:   maxBlocksFlag,
			Usage:  "The maximum number of block on each query",
			EnvVar: "MAX_BLOCKS",
			Value:  defaultMaxBlocks,
		},
		cli.IntFlag{
			Name:   attemptsFlag,
			Usage:  "The number of attempt to query trade log from blockchain",
			EnvVar: "ATTEMPTS",
			Value:  defaultAttempts,
		},
		cli.DurationFlag{
			Name:   delayFlag,
			Usage:  "The duration to put worker pools into sleep after each batch requets",
			EnvVar: "DELAY",
			Value:  defaultDelayTime,
		},
	)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, core.NewCliFlags()...)
	app.Flags = append(app.Flags, broadcast.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, cq.NewCQFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseBigIntFlag(c *cli.Context, flag string) (*big.Int, error) {
	val := c.String(flag)
	if err := validation.Validate(val, validation.Required, is.Digit); err != nil {
		return nil, err
	}

	result, ok := big.NewInt(0).SetString(val, 0)
	if !ok {
		return nil, fmt.Errorf("invalid number %s", c.String(flag))
	}
	return result, nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func manageCQFromContext(c *cli.Context, influxClient client.Client, sugar *zap.SugaredLogger) error {
	//Deploy CQ	before get/ store tradelog
	cqs, err := tradelogcq.CreateAssetVolumeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	burnFeeCqs, err := tradelogcq.CreateBurnFeeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, burnFeeCqs...)
	walletFeeCqs, err := tradelogcq.CreateWalletFeeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, walletFeeCqs...)
	return cq.ManageCQs(c, cqs, influxClient, sugar)
}

func run(c *cli.Context) error {
	var (
		err       error
		fromBlock *big.Int
		toBlock   *big.Int
		daemon    bool
	)

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}

	defer logger.Sync()

	sugar := logger.Sugar()

	coreClient, err := core.NewClientFromContext(sugar, c)
	if err != nil {
		return err
	}

	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	influxStorage, err := storage.NewInfluxStorage(
		sugar,
		common.DatabaseName,
		influxClient,
		core.NewCachedClient(coreClient),
	)
	if err != nil {
		return err
	}
	if err = manageCQFromContext(c, influxClient, sugar); err != nil {
		return err
	}
	if c.String(fromBlockFlag) == "" {
		sugar.Info("no from block flag provided, checking last stored block")
	} else {
		fromBlock, err = parseBigIntFlag(c, fromBlockFlag)
		if err != nil {
			return err
		}
	}

	if c.String(toBlockFlag) == "" {
		daemon = true
		sugar.Info("running in daemon mode")
	} else {
		toBlock, err = parseBigIntFlag(c, toBlockFlag)
		if err != nil {
			return err
		}
	}

	maxWorker := c.Int(maxWorkerFlag)
	maxBlock := c.Int(maxBlocksFlag)
	attempts := c.Int(attemptsFlag) // exit if failed to fetch logs after attempts times
	delayTime := c.Duration(delayFlag)

	for {
		var doneCh = make(chan struct{})

		if fromBlock == nil {
			lastBlock, fErr := influxStorage.LastBlock()
			if fErr != nil {
				return fErr
			}
			if lastBlock != 0 {
				fromBlock = big.NewInt(0).SetInt64(lastBlock)
				sugar.Infow("using last stored block number", "from_block", fromBlock)
			} else {
				sugar.Infow("using default from block number", "from_block", defaultFromblock)
				fromBlock = big.NewInt(0).SetInt64(defaultFromblock)
			}
		}

		if toBlock == nil {
			ethClient, fErr := libapp.NewEthereumClientFromFlag(c)
			currentHeader, fErr := ethClient.HeaderByNumber(context.Background(), nil)
			if fErr != nil {
				return fErr
			}
			toBlock = currentHeader.Number
			sugar.Infow("fetching trade logs up to latest known block number", "to_block", toBlock.String())
		}

		jobs := int(math.Ceil(float64(toBlock.Int64()-fromBlock.Int64()) / float64(maxBlock)))
		if jobs < maxWorker {
			maxWorker = jobs // if jobs < maxWorkers, jobs = n, only start n workers
		}
		p := workers.NewPool(sugar, maxWorker, influxStorage)
		sugar.Debugw("number of fetcher jobs", "jobs", jobs, "max_blocks", maxBlock)

		go func(fromBlock, toBlock, maxBlocks int64) {
			var jobOrder = p.GetLastCompleteJobOrder()
			for i := int64(fromBlock); i < toBlock; i = i + maxBlocks {
				jobOrder++
				p.Run(workers.NewFetcherJob(c, jobOrder, big.NewInt(i), big.NewInt(min(i+maxBlocks, toBlock)), attempts))
			}
			for p.GetLastCompleteJobOrder() < jobOrder {
				time.Sleep(time.Second)
			}
			doneCh <- struct{}{}
		}(fromBlock.Int64(), toBlock.Int64(), int64(maxBlock))

		for {
			var toBreak = false
			select {
			case <-doneCh:
				sugar.Info("all jobs are successfully executed, waiting for the workers pool to shut down")
				p.Shutdown()
			case fErr := <-p.ErrCh():
				if fErr != nil {
					sugar.Fatalw("job failed to execute", "error", fErr)
				} else {
					sugar.Info("workers pool is successfully shut down")
				}
				toBreak = true
			}

			if toBreak {
				break
			}
		}

		if daemon {
			sugar.Infow("waiting before fetching new trade logs",
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
