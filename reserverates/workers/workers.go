package workers

import (
	"errors"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/crawler"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

type job interface {
	execute(sugar *zap.SugaredLogger) (map[string]map[string]common.ReserveRateEntry, error)
	info() (order int, block uint64)
}

// FetcherJob represent a job to crawl rates at given block
type FetcherJob struct {
	c        *cli.Context
	order    int
	block    uint64
	attempts int
	addrs    []string
}

// NewFetcherJob return an instance of FetcherJob
func NewFetcherJob(c *cli.Context, order int, block uint64, addrs []string, attempts int) *FetcherJob {
	return &FetcherJob{
		c:        c,
		order:    order,
		block:    block,
		attempts: attempts,
		addrs:    addrs,
	}
}

// retry the given fn function for attempts time with sleep duration between before returns an error.
func retry(fn func(*zap.SugaredLogger) (map[string]map[string]common.ReserveRateEntry, error), attempts int, logger *zap.SugaredLogger) (map[string]map[string]common.ReserveRateEntry, error) {
	var (
		result map[string]map[string]common.ReserveRateEntry
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

	return nil, err
}

func (fj *FetcherJob) fetch(sugar *zap.SugaredLogger) (map[string]map[string]common.ReserveRateEntry, error) {

	client, err := app.NewEthereumClientFromFlag(fj.c)
	if err != nil {
		return nil, err
	}

	coreClient, err := core.NewClientFromContext(sugar, fj.c)
	if err != nil {
		return nil, err
	}

	ratesCrawler, err := crawler.NewReserveRatesCrawler(sugar, fj.addrs, client, coreClient)
	if err != nil {
		return nil, err
	}

	rates, err := ratesCrawler.GetReserveRates(fj.block)
	if err != nil {
		return nil, err
	}

	return rates, nil
}

func (fj *FetcherJob) execute(sugar *zap.SugaredLogger) (map[string]map[string]common.ReserveRateEntry, error) {
	return retry(fj.fetch, fj.attempts, sugar)
}

func (fj *FetcherJob) info() (int, uint64) {
	return fj.order, fj.block
}

// Pool represents a group of workers
type Pool struct {
	sugar *zap.SugaredLogger

	wg sync.WaitGroup

	jobCh chan job
	errCh chan error

	mutex                 *sync.Mutex
	lastCompletedJobOrder int  // Keep the order of the last completed job
	failed                bool // mark as failed, all subsequent persistent storage will be passed

	rateStorage storage.ReserveRatesStorage
}

// NewPool returns a pool of workers
func NewPool(sugar *zap.SugaredLogger, maxWorkers int, rateStorage storage.ReserveRatesStorage) *Pool {
	var pool = &Pool{
		sugar:                 sugar,
		jobCh:                 make(chan job),
		errCh:                 make(chan error, maxWorkers),
		mutex:                 &sync.Mutex{},
		lastCompletedJobOrder: 0,
		rateStorage:           rateStorage,
	}

	pool.wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go func(sugar *zap.SugaredLogger, workerID int) {
			logger := sugar.With("worker_id", workerID)
			logger.Infow("starting worker",
				"func", "reserverates/workers/NewPool",
				"max_workers", maxWorkers)

			for j := range pool.jobCh {
				order, block := j.info()
				rates, err := j.execute(sugar)
				if err != nil {
					logger.Errorw("fetcher job execution failed",
						"block", block,
						"err", err)
					pool.errCh <- err
					pool.markAsFailed(order)
					break
				}

				logger.Infow("fetcher job executed successfully",
					"block", block)
				if err = pool.serialSaveTradeLogs(order, block, rates); err != nil {
					pool.errCh <- err
					break
				}
			}

			logger.Infow("worker stopped",
				"func", "reserverates/workers/NewPool",
				"max_workers", maxWorkers)
			pool.wg.Done()
		}(sugar, i)
	}

	return pool
}

func (p *Pool) markAsFailed(order int) {
	var (
		logger = p.sugar.With(
			"func", "reserverates/workers/Pool.markAsFailed",
			"order", order,
		)
	)
	for {
		p.mutex.Lock()
		if order == p.lastCompletedJobOrder+1 {
			logger.Warn("mark as failed")
			p.failed = true
			p.mutex.Unlock()
			return
		}

		p.mutex.Unlock()
		time.Sleep(time.Second)
	}
}

func (p *Pool) serialSaveTradeLogs(
	order int,
	blockNumber uint64,
	rates map[string]map[string]common.ReserveRateEntry) error {
	var (
		logger = p.sugar.With(
			"func", "reserverates/workers/Pool.serialSaveTradeLogs",
			"order", order,
			"block_number", blockNumber,
		)
		err error
	)

	for {
		p.mutex.Lock()

		if p.failed {
			p.mutex.Unlock()
			return errors.New("pool has been marked as failed")
		}

		if order == p.lastCompletedJobOrder+1 {
			if err = p.rateStorage.UpdateRatesRecords(blockNumber, rates); err != nil {
				logger.Error("saving rates to persistent storage",
					"err", err)
				p.mutex.Unlock()
				p.markAsFailed(order)
				return err
			}

			p.lastCompletedJobOrder++
			logger.Infow("save rates to storage success")
			p.mutex.Unlock()
			return nil
		}

		logger.Debugw("waiting for previous job to be completed",
			"last_completed", p.lastCompletedJobOrder)
		p.mutex.Unlock()
		time.Sleep(time.Second)
	}
}

// GetLastCompleteJobOrder return the order of the latest completed job
func (p *Pool) GetLastCompleteJobOrder() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.lastCompletedJobOrder
}

// Run puts new job to queue
func (p *Pool) Run(j job) {
	order, block := j.info()
	p.sugar.Infow("putting new job to queue",
		"func", "reserverates/workers/Run",
		"order", order,
		"block", block)
	p.jobCh <- j
}

// Shutdown stops the workers pool
func (p *Pool) Shutdown() {
	p.sugar.Infow("workers pool shutting down",
		"func", "reserverates/workers/Shutdown")
	close(p.jobCh)
	p.wg.Wait()
	close(p.errCh)
}

// ErrCh returns error reporting channel of workers pool.
func (p *Pool) ErrCh() chan error {
	return p.errCh
}
