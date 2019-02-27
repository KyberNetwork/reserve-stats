package main

import (
	"log"
	"math"
	"math/big"
	"os"
	"time"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	storage "github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage"
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

	defaultDB = "reserve_stats"
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
	)
	app.Flags = append(app.Flags, broadcast.NewCliFlags()...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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

func run(c *cli.Context) error {
	var (
		err error
	)

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}

	defer logger.Sync()

	sugar := logger.Sugar()

	// tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
	// if err != nil {
	// 	return err
	// }
	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
	if err != nil {
		return err
	}

	postgreStorage, err := storage.NewTradeLogDB(sugar, db, tokenAmountFormatter)
	if err != nil {
		sugar.Infow("error", err)
		return err
	}

	maxWorkers := c.Int(maxWorkersFlag)
	maxBlocks := c.Int(maxBlocksFlag)
	attempts := c.Int(attemptsFlag) // exit if failed to fetch logs after attempts times
	planner, err := newCrawlerPlanner(sugar, c, postgreStorage)
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
		p := workers.NewPool(sugar, requiredWorkers, postgreStorage)
		sugar.Debugw("number of fetcher jobs",
			"from_block", fromBlock.String(),
			"to_block", toBlock.String(),
			"workers", requiredWorkers,
			"max_blocks", maxBlocks)

		go func(fromBlock, toBlock, maxBlocks int64) {
			var jobOrder = p.GetLastCompleteJobOrder()
			for i := int64(fromBlock); i < toBlock; i = i + maxBlocks {
				jobOrder++
				p.Run(workers.NewFetcherJob(c, jobOrder, big.NewInt(i), big.NewInt(mathutil.MinInt64(i+maxBlocks, toBlock)), attempts))
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
