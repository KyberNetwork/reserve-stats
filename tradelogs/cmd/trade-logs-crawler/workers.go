package main

import (
	"errors"
	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/tokenrate/coingecko"
	"math/big"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

// TODO: refactor all below to a packages in tradeslogs/

const dbName = "trade_logs"

type job interface {
	// TODO: wrap execute in a retry function with following signature.
	// retry the given fn function for attempts time with sleep duration between before returns an error.
	// retry(fn function(*zap.SugaredLogger), attempts int)([]common.TradeLog, error)
	execute(sugar *zap.SugaredLogger) ([]common.TradeLog, error)
	info() (order int, from, to *big.Int)
}

func newFetcherJob(c *cli.Context, order int, from, to *big.Int) *fetcherJob {
	return &fetcherJob{
		c:     c,
		order: order,
		from:  from,
		to:    to,
	}
}

type fetcherJob struct {
	c     *cli.Context
	order int
	from  *big.Int
	to    *big.Int
}

func (fj *fetcherJob) execute(sugar *zap.SugaredLogger) ([]common.TradeLog, error) {
	logger := sugar.With(
		"from", fj.from.String(),
		"to", fj.to.String())
	logger.Debugw("fetching trade logs")

	bc, err := broadcast.NewClientFromContext(sugar, fj.c)
	if err != nil {
		return nil, err
	}

	coreClient, err := core.NewClientFromContext(sugar, fj.c)
	if err != nil {
		return nil, err
	}

	influxClient, err := influxdb.NewClientFromContext(fj.c)
	if err != nil {
		return nil, err
	}

	is, err := storage.NewInfluxStorage(
		sugar,
		dbName,
		influxClient,
		core.NewCachedClient(coreClient),
	)
	if err != nil {
		return nil, err
	}

	client, err := app.NewEthereumClientFromFlag(fj.c)
	if err != nil {
		return nil, err
	}

	crawler, err := tradelogs.NewCrawler(sugar, client, bc, is, coingecko.New())
	if err != nil {
		return nil, err
	}

	tradeLogs, err := crawler.GetTradeLogs(fj.from, fj.to, time.Second*5)
	if err != nil {
		return nil, err
	}

	return nil, errors.New("I always return an error")

	return tradeLogs, nil
}

func (fj *fetcherJob) info() (int, *big.Int, *big.Int) {
	return fj.order, fj.from, fj.to
}

type pool struct {
	sugar *zap.SugaredLogger
	wg    sync.WaitGroup

	jobCh chan job
	errCh chan error

	// TODO: add a lastCompletedJob field and only inserted to database if the job order = lastCompleted + 1
	// that means only the fetching job are running concurrently, the database writing is still consecutively
	// this should be protected with a sync mutex or using atomic.StoreInt.
	//lastCompleted int
}

func newPool(sugar *zap.SugaredLogger, maxWorkers int) *pool {
	var p = &pool{
		sugar: sugar,
		jobCh: make(chan job),
		errCh: make(chan error, maxWorkers),
	}

	p.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func(sugar *zap.SugaredLogger, workerID int) {
			logger := sugar.With("worker_id", workerID)
			logger.Infow("starting worker",
				"func", "tradelogs/cmd/trade-logs-crawler/newPool",
				"max_workers", maxWorkers,
			)
			for j := range p.jobCh {
				_, from, to := j.info()
				logger.Infow("executing fetcher job",
					"from", from.String(),
					"to", to.String())

				// TODO: compare the job id to lastCompleted and save to database if job order = lastCompleted + 1,
				// otherwise blocks worker
				if _, err := j.execute(logger); err != nil {
					logger.Errorw("fetcher job execution failed",
						"from", from.String(),
						"to", to.String(),
						"err", err)
					p.errCh <- err
					break
				}
				logger.Infow("fetcher job executed successfully",
					"from", from.String(),
					"to", to.String())

			}
			logger.Infow("worker stopped",
				"func", "tradelogs/cmd/trade-logs-crawler/newPool",
				"max_workers", maxWorkers,
			)
			p.wg.Done()
		}(sugar, i)
	}

	return p
}

func (p *pool) run(j job) {
	_, from, to := j.info()
	p.sugar.Infow("putting new job to queue",
		"func", "tradelogs/cmd/trade-logs-crawler/run",
		"from", from.String(),
		"to", to.String())
	p.jobCh <- j
}

func (p *pool) shutdown() {
	p.sugar.Infow("workers pool shutting down",
		"func", "tradelogs/cmd/trade-logs-crawler/shutdown")
	close(p.jobCh)
	p.wg.Wait()
	close(p.errCh)
}
