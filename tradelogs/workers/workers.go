package workers

import (
	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/tokenrate/coingecko"
	"math/big"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

type executeJob func(*zap.SugaredLogger) ([]common.TradeLog, error)

type job interface {
	execute(sugar *zap.SugaredLogger) ([]common.TradeLog, error)
	info() (order int, from, to *big.Int)
}

// NewFetcherJob return an instance of fetcherJob
func NewFetcherJob(c *cli.Context, order int, from, to *big.Int, attempts int) *FetcherJob {
	return &FetcherJob{
		c:        c,
		order:    order,
		from:     from,
		to:       to,
		attempts: attempts,
	}
}

// FetcherJob represent a job to crawl trade logs from block to block
type FetcherJob struct {
	c        *cli.Context
	order    int
	from     *big.Int
	to       *big.Int
	attempts int
}

// retry the given fn function for attempts time with sleep duration between before returns an error.
func retry(fn executeJob, attempts int, logger *zap.SugaredLogger) ([]common.TradeLog, error) {
	var (
		result []common.TradeLog
		err    error
	)

	for i := 0; i < attempts; i++ {
		result, err = fn(logger)
		if err == nil {
			return result, nil
		}

		logger.Debugw("failed to execute job", "attempt", i)
		time.Sleep(time.Second)
	}

	return result, err
}

func (fj *FetcherJob) fetch(sugar *zap.SugaredLogger) ([]common.TradeLog, error) {
	logger := sugar.With(
		"from", fj.from.String(),
		"to", fj.to.String())
	logger.Debugw("fetching trade logs")

	bc, err := broadcast.NewClientFromContext(sugar, fj.c)
	if err != nil {
		return nil, err
	}

	client, err := app.NewEthereumClientFromFlag(fj.c)
	if err != nil {
		return nil, err
	}

	crawler, err := tradelogs.NewCrawler(sugar, client, bc, coingecko.New())
	if err != nil {
		return nil, err
	}

	tradeLogs, err := crawler.GetTradeLogs(fj.from, fj.to, time.Second*5)
	if err != nil {
		return nil, err
	}

	return tradeLogs, nil
}

func (fj *FetcherJob) execute(sugar *zap.SugaredLogger) ([]common.TradeLog, error) {
	return retry(fj.fetch, fj.attempts, sugar)
}

func (fj *FetcherJob) info() (int, *big.Int, *big.Int) {
	return fj.order, fj.from, fj.to
}

// Pool represents a group of workers which is capable of handle many jobs
// at a time
type Pool struct {
	sugar *zap.SugaredLogger
	wg    sync.WaitGroup

	jobCh chan job
	ErrCh chan error

	mutex                 *sync.Mutex
	lastCompletedJobOrder int // Keep the order of the last completed job

	storage storage.Interface
}

// NewPool returns a pool of workers to handle jobs concurrently
func NewPool(sugar *zap.SugaredLogger, maxWorkers int, storage storage.Interface) *Pool {
	var p = &Pool{
		sugar:                 sugar,
		jobCh:                 make(chan job),
		ErrCh:                 make(chan error, maxWorkers),
		mutex:                 &sync.Mutex{},
		storage:               storage,
		lastCompletedJobOrder: 0,
	}

	p.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func(sugar *zap.SugaredLogger, workerID int) {
			logger := sugar.With("worker_id", workerID)
			logger.Infow("starting worker",
				"func", "tradelogs/workers/NewPool",
				"max_workers", maxWorkers,
			)
			for j := range p.jobCh {
				order, from, to := j.info()
				logger.Infow("executing fetcher job",
					"from", from.String(),
					"to", to.String())

				tradeLogs, err := j.execute(sugar)
				if err != nil {
					logger.Errorw("fetcher job execution failed",
						"from", from.String(),
						"to", to.String(),
						"err", err)
					p.ErrCh <- err
					break
				}
				logger.Infow("fetcher job executed successfully",
					"from", from.String(),
					"to", to.String())

				// Compare the job order to lastCompletedOrder and save to database if job order = lastCompletedOrder + 1,
				// otherwise blocks worker
				saveSuccess := false
				for {
					var err error

					p.mutex.Lock()
					if order == p.lastCompletedJobOrder+1 {
						err = p.storage.SaveTradeLogs(tradeLogs)
						if err == nil {
							saveSuccess = true
							p.lastCompletedJobOrder++
						}
					}
					p.mutex.Unlock()

					if err != nil {
						logger.Errorw("save trade logs into db failed",
							"from", from.String(),
							"to", to.String(),
							"err", err)
						p.ErrCh <- err
						break
					} else {
						if saveSuccess {
							logger.Infow("save trade logs into db success",
								"from", from.String(),
								"to", to.String())
							break
						} else {
							time.Sleep(time.Second)
						}
					}
				}
			}
			logger.Infow("worker stopped",
				"func", "tradelogs/workers/NewPool",
				"max_workers", maxWorkers,
			)
			p.wg.Done()
		}(sugar, i)
	}

	return p
}

// GetLastCompleteJobOrder return the order of the latest completed job
func (p *Pool) GetLastCompleteJobOrder() int {
	p.mutex.Lock()
	result := p.lastCompletedJobOrder
	p.mutex.Unlock()

	return result
}

// Run puts new job to queue
func (p *Pool) Run(j job) {
	_, from, to := j.info()
	p.sugar.Infow("putting new job to queue",
		"func", "tradelogs/workers/Run",
		"from", from.String(),
		"to", to.String())
	p.jobCh <- j
}

// Shutdown stops the workers pool
func (p *Pool) Shutdown() {
	p.sugar.Infow("workers pool shutting down",
		"func", "tradelogs/workers/Shutdown")
	close(p.jobCh)
	p.wg.Wait()
	close(p.ErrCh)
}
