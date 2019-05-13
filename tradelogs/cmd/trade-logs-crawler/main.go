package main

import (
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
	"github.com/KyberNetwork/reserve-stats/lib/etherscan"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	userkyced "github.com/KyberNetwork/reserve-stats/lib/userkyced"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
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

	defaultDB = "users"
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
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, broadcast.NewCliFlags()...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, cq.NewCQFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)
	app.Flags = append(app.Flags, userkyced.NewCliFlags()...)
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

func run(c *cli.Context) error {
	var (
		err        error
		kycChecker storage.KycChecker
	)

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	if err = manageCQFromContext(c, influxClient, sugar); err != nil {
		return err
	}
	userKycedClient, err := userkyced.NewClientFromContext(sugar, c)
	switch err {
	case userkyced.ErrNoClientURL:
		sugar.Info("User kyced checker URL is not provided. Use default Postgres instead")
		db, err := libapp.NewDBFromContext(c)
		if err != nil {
			return err
		}
		kycChecker = storage.NewUserKYCChecker(sugar, db)
	case nil:
		sugar.Info("User kyced checker URL provided. check KYCed status from userKYCed client")
		kycChecker = userKycedClient
	default:
		return err
	}

	tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
	if err != nil {
		return err
	}

	influxStorage, err := storage.NewInfluxStorage(
		sugar,
		common.DatabaseName,
		influxClient,
		tokenAmountFormatter,
		kycChecker,
	)
	if err != nil {
		return err
	}

	maxWorkers := c.Int(maxWorkersFlag)
	maxBlocks := c.Int(maxBlocksFlag)
	attempts := c.Int(attemptsFlag) // exit if failed to fetch logs after attempts times
	planner, err := newCrawlerPlanner(sugar, c, influxStorage)
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
		p := workers.NewPool(sugar, requiredWorkers, influxStorage)
		sugar.Debugw("number of fetcher jobs",
			"from_block", fromBlock.String(),
			"to_block", toBlock.String(),
			"workers", requiredWorkers,
			"max_blocks", maxBlocks)

		go func(fromBlock, toBlock, maxBlocks int64) {
			var jobOrder = p.GetLastCompleteJobOrder()
			for i := fromBlock; i < toBlock; i += maxBlocks {
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
