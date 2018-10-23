package main

import (
	"context"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const (
	fromBlockFlag    = "from-block"
	defaultFromblock = 5069586
	toBlockFlag      = "to-block"

	envVarPrefix = "TRADE_LOGS_CRAWLER_"
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
	)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, core.NewCliFlags()...)
	app.Flags = append(app.Flags, broadcast.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewEthereumNodeFlags())

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

func run(c *cli.Context) error {
	const (
		// TODO: make this a configuration flag with default value: 5
		maxWorker = 20
		// TODO: make this a configuration flag with default value: 100
		maxBlocks = 1000
	)
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

	// TODO: consider exit if successive failed to fetch logs for 5 times
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
			client, fErr := libapp.NewEthereumClientFromFlag(c)
			currentHeader, fErr := client.HeaderByNumber(context.Background(), nil)
			if fErr != nil {
				return fErr
			}
			toBlock = currentHeader.Number
			sugar.Infow("fetching trade logs up to latest known block number", "to_block", toBlock.String())
		}

		// TODO: if jobs < maxWorkers, jobs = n, only start n workers
		jobs := int(math.Ceil(float64(toBlock.Int64()-fromBlock.Int64()) / maxBlocks))
		p := newPool(sugar, maxWorker)
		sugar.Debugw("number of fetcher jobs", "jobs", jobs, "max_blocks", maxBlocks)

		go func(fromBlock, toBlock, maxBlocks int64) {
			var order = 0
			for i := int64(fromBlock); i < toBlock; i = i + maxBlocks {
				p.run(newFetcherJob(c, order, big.NewInt(i), big.NewInt(i+maxBlocks)))
				order++
			}
			doneCh <- struct{}{}
		}(fromBlock.Int64(), toBlock.Int64(), maxBlocks)

		for {
			var toBreak = false
			select {
			case <-doneCh:
				sugar.Info("all jobs are queued, waiting for completion")
				p.shutdown()
			case fErr := <-p.errCh:
				if fErr != nil {
					sugar.Errorw("job failed to execute", "error", fErr)
				} else {
					sugar.Info("all jobs are successfully executed")
				}
				toBreak = true
			}

			if toBreak {
				break
			}
		}

		if daemon {
			// TODO: make this a flag
			sleep := time.Minute
			sugar.Infow("waiting before fetching new trade logs",
				"last_from_block", fromBlock.String(),
				"last_to_block", toBlock.String(),
				"sleep", sleep.String())
			fromBlock, toBlock = nil, nil
			time.Sleep(sleep)
		} else {
			break
		}
	}

	return nil
}
