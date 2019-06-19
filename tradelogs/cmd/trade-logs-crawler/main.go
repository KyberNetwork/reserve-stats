package main

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage"
	"github.com/KyberNetwork/reserve-stats/tradelogs/workers"
)

const (
	fromBlockFlag    = "from-block"
	defaultFromBlock = 5069586
	toBlockFlag      = "to-block"

	maxWorkersFlag    = "max-workers"
	defaultMaxWorkers = 5

	maxBlocksFlag    = "max-blocks"
	defaultMaxBlocks = 100

	attemptsFlag    = "attempts"
	defaultAttempts = 5

	delayFlag        = "delay"
	defaultDelayTime = time.Minute

	blockConfirmationsFlag    = "wait-for-confirmations"
	defaultBlockConfirmations = 7

	dbEngineFlag    = "db-engine"
	defaultDbEngine = "influx"
	influxDbEngine  = "influx"
	postgreDbEngine = "postgre"

	defaultDb = "reserve_stats"
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
			Name:   maxWorkersFlag,
			Usage:  "The maximum number of worker to fetch trade logs",
			EnvVar: "MAX_WORKER",
			Value:  defaultMaxWorkers,
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
			Usage:  "The duration to put worker pools into sleep after each batch requests",
			EnvVar: "DELAY",
			Value:  defaultDelayTime,
		},
		cli.Int64Flag{
			Name:   blockConfirmationsFlag,
			Usage:  "The number of block confirmations to latest known block",
			EnvVar: "WAIT_FOR_CONFIRMATIONS",
			Value:  defaultBlockConfirmations,
		},
		cli.StringFlag{
			Name:   dbEngineFlag,
			Usage:  "db engine to write trade logs, pls select influx or postgre",
			EnvVar: "DB_ENGINE",
			Value:  defaultDbEngine,
		},
	)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDb)...)
	app.Flags = append(app.Flags, broadcast.NewCliFlags()...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, cq.NewCQFlags()...)
	app.Flags = append(app.Flags, etherscan.NewCliFlags()...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func manageCQFromContext(c *cli.Context, influxClient client.Client, sugar *zap.SugaredLogger) error {
	//Deploy CQ	before get/store trade logs
	cqs, err := tradelogcq.CreateAssetVolumeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	reserveVolumeCqs, err := tradelogcq.CreateReserveVolumeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, reserveVolumeCqs...)
	userVolumeCqs, err := tradelogcq.CreateUserVolumeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, userVolumeCqs...)
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
	summaryCqs, err := tradelogcq.CreateSummaryCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, summaryCqs...)
	walletStatsCqs, err := tradelogcq.CreateWalletStatsCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, walletStatsCqs...)
	countryStatsCqs, err := tradelogcq.CreateCountryCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, countryStatsCqs...)
	integrationVolumeCqs, err := tradelogcq.CreateIntegrationVolumeCq(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, integrationVolumeCqs...)

	return cq.ManageCQs(c, cqs, influxClient, sugar)
}

// requiredWorkers returns number of workers to start. If the number of jobs is smaller than max workers,
// only start the number of required workers instead of max workers.
func requiredWorkers(fromBlock, toBlock *big.Int, maxBlocks, maxWorkers int) int {
	jobs := int(math.Ceil(float64(toBlock.Int64()-fromBlock.Int64()) / float64(maxBlocks)))
	if jobs < maxWorkers {
		return jobs
	}
	return maxWorkers
}

func newStorageInterface(sugar *zap.SugaredLogger, c *cli.Context) (storage.Interface, error) {
	var (
		err error
	)

	tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
	if err != nil {
		return nil, err
	}

	dbEngine := c.String(dbEngineFlag)
	switch dbEngine {
	case influxDbEngine:
		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return nil, err
		}

		if err = manageCQFromContext(c, influxClient, sugar); err != nil {
			return nil, err
		}

		influxStorage, err := influxstorage.NewInfluxStorage(
			sugar,
			common.DatabaseName,
			influxClient,
			tokenAmountFormatter,
		)
		if err != nil {
			return nil, err
		}
		return influxStorage, nil
	case postgreDbEngine:
		db, err := libapp.NewDBFromContext(c)
		if err != nil {
			return nil, err
		}
		postgresStorage, err := postgrestorage.NewTradeLogDB(sugar, db, tokenAmountFormatter)
		if err != nil {
			sugar.Infow("error", err)
			return nil, err
		}
		return postgresStorage, nil
	default:
		return nil, fmt.Errorf("invalid db engine: %q", dbEngine)
	}
}

func run(c *cli.Context) error {
	var (
		err              error
		storageInterface storage.Interface
	)

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	storageInterface, err = newStorageInterface(sugar, c)
	if err != nil {
		return err
	}

	etherscanClient, err := etherscan.NewEtherscanClientFromContext(c)
	if err != nil {
		return err
	}

	maxWorkers := c.Int(maxWorkersFlag)
	maxBlocks := c.Int(maxBlocksFlag)
	attempts := c.Int(attemptsFlag) // exit if failed to fetch logs after attempts times
	planner, err := newCrawlerPlanner(sugar, c, storageInterface)
	if err != nil {
		return nil
	}

	for {
		var (
			doneCh             = make(chan struct{})
			fromBlock, toBlock *big.Int
		)
		if fromBlock, toBlock, err = planner.next(); err == errEOF {
			sugar.Info("completed!")
			break
		} else if err != nil {
			return err
		}

		requiredWorkers := requiredWorkers(fromBlock, toBlock, maxBlocks, maxWorkers)
		startingBlocks := deployment.MustGetStartingBlocksFromContext(c)
		p := workers.NewPool(sugar, requiredWorkers, storageInterface)
		sugar.Debugw("number of fetcher jobs",
			"from_block", fromBlock.String(),
			"to_block", toBlock.String(),
			"workers", requiredWorkers,
			"max_blocks", maxBlocks,
			"starting_blocks", startingBlocks)

		go func(fromBlock, toBlock, maxBlocks int64) {
			var jobOrder = p.GetLastCompleteJobOrder()
			for i := fromBlock; i < toBlock; i += maxBlocks {
				end := mathutil.MinInt64(i+maxBlocks, toBlock)
				switch {
				//if job start at block v2 and end at block v3 then split job
				case uint64(end) >= startingBlocks.V3() && uint64(i) < startingBlocks.V3():
					jobOrder++
					p.Run(workers.NewFetcherJob(c, jobOrder, big.NewInt(i), big.NewInt(int64(startingBlocks.V3())), attempts, etherscanClient))
					jobOrder++
					p.Run(workers.NewFetcherJob(c, jobOrder, big.NewInt(int64(startingBlocks.V3())), big.NewInt(end), attempts, etherscanClient))
				//if job start at block v1 and end at block v2 then split job
				case uint64(end) >= startingBlocks.V2() && uint64(i) < startingBlocks.V2():
					jobOrder++
					p.Run(workers.NewFetcherJob(c, jobOrder, big.NewInt(i), big.NewInt(int64(startingBlocks.V2())), attempts, etherscanClient))
					jobOrder++
					p.Run(workers.NewFetcherJob(c, jobOrder, big.NewInt(int64(startingBlocks.V2())), big.NewInt(end), attempts, etherscanClient))
				default:
					jobOrder++
					p.Run(workers.NewFetcherJob(c, jobOrder, big.NewInt(i), big.NewInt(end), attempts, etherscanClient))
				}
			}
			for p.GetLastCompleteJobOrder() < jobOrder {
				time.Sleep(time.Second)
			}
			doneCh <- struct{}{}
		}(fromBlock.Int64(), toBlock.Int64(), int64(maxBlocks))

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
	}
	return nil
}
