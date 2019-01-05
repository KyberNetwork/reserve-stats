package workers

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"

	"github.com/KyberNetwork/tokenrate/coingecko"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
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

func (fj *FetcherJob) getInternalNetworkAddress(contractAddressClient *contracts.CachedContractAddressClient,
	proxyAddress ethereum.Address) ([]ethereum.Address, error) {
	var (
		addresses []ethereum.Address
	)
	addressFrom, err := contractAddressClient.InternalNetworkContractAddress(proxyAddress, fj.from)
	if err != nil {
		return addresses, err
	}
	addresses = append(addresses, addressFrom)
	addressTo, err := contractAddressClient.InternalNetworkContractAddress(proxyAddress, fj.to)
	if err != nil {
		return addresses, err
	}
	if addressTo != addressFrom {
		addresses = append(addresses, addressTo)
	}
	return addresses, nil
}

func (fj *FetcherJob) getPricingAddress(contractAddressClient *contracts.CachedContractAddressClient,
	internalNetworkAddress ethereum.Address) ([]ethereum.Address, error) {
	var (
		addresses []ethereum.Address
	)
	addressFrom, err := contractAddressClient.PricingContractAddress(internalNetworkAddress, fj.from)
	if err != nil {
		return addresses, err
	}
	addresses = append(addresses, addressFrom)
	addressTo, err := contractAddressClient.PricingContractAddress(internalNetworkAddress, fj.to)
	if err != nil {
		return addresses, err
	}
	if addressTo != addressFrom {
		addresses = append(addresses, addressTo)
	}
	return addresses, nil
}

func (fj *FetcherJob) getBurnerAddress(contractAddressClient *contracts.CachedContractAddressClient,
	internalNetworkAddress ethereum.Address) ([]ethereum.Address, error) {
	var (
		addresses []ethereum.Address
	)
	addressFrom, err := contractAddressClient.BurnerContractAddress(internalNetworkAddress, fj.from)
	if err != nil {
		return addresses, err
	}
	addresses = append(addresses, addressFrom)
	addressTo, err := contractAddressClient.BurnerContractAddress(internalNetworkAddress, fj.to)
	if err != nil {
		return addresses, err
	}
	if addressTo != addressFrom {
		addresses = append(addresses, addressTo)
	}
	return addresses, nil
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

	client, err := blockchain.NewEthereumClientFromFlag(fj.c)
	if err != nil {
		return nil, err
	}
	cachedContractAddressClient := contracts.NewCachedContractAddressClient(client)
	proxyAddress := contracts.ProxyContractAddress().MustGetOneFromContext(fj.c)
	addresses := []ethereum.Address{proxyAddress}

	logger.Debugw("get addresses")

	internalNetworkAddress, err := fj.getInternalNetworkAddress(cachedContractAddressClient, proxyAddress)
	if err != nil {
		logger.Debug("internal network address error:", err.Error())
		return nil, err
	}
	addresses = append(addresses, internalNetworkAddress...)

	for _, address := range internalNetworkAddress {
		pricingAddress, err := fj.getPricingAddress(cachedContractAddressClient, address)
		if err != nil {
			logger.Debug("internal network address: ", address.Hex())
			logger.Debug("get pricing address error:", err.Error())
			return nil, err
		}
		addresses = append(addresses, pricingAddress...)

		burnerAddress, err := fj.getBurnerAddress(cachedContractAddressClient, address)
		if err != nil {
			logger.Debug("internal network address: ", address.Hex())
			logger.Debug("get burner address error:", err.Error())
			return nil, err
		}
		addresses = append(addresses, burnerAddress...)
	}

	crawler, err := tradelogs.NewCrawler(logger, client, bc, coingecko.New(), addresses, startingBlocks)
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
	errCh chan error

	mutex                 *sync.Mutex
	lastCompletedJobOrder int  // Keep the order of the last completed job
	failed                bool // mark as failed, all subsequent persistent storage will be passed
	storage               storage.Interface
}

// NewPool returns a pool of workers to handle jobs concurrently
func NewPool(sugar *zap.SugaredLogger, maxWorkers int, storage storage.Interface) *Pool {
	var p = &Pool{
		sugar:                 sugar,
		jobCh:                 make(chan job),
		errCh:                 make(chan error, maxWorkers),
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
					"order", order,
					"from", from.String(),
					"to", to.String())

				tradeLogs, err := j.execute(sugar)
				if err != nil {
					logger.Errorw("fetcher job execution failed",
						"from", from.String(),
						"to", to.String(),
						"err", err)
					p.errCh <- err
					p.markAsFailed(order)
					break
				}

				logger.Infow("fetcher job executed successfully",
					"order", order,
					"from", from.String(),
					"to", to.String())
				if err = p.serialSaveTradeLogs(order, tradeLogs); err != nil {
					p.errCh <- err
					break
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

func (p *Pool) markAsFailed(order int) {
	var (
		logger = p.sugar.With(
			"func", "tradelogs/workers/Pool.markAsFailed",
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

// serialSaveTradeLogs waits until the job with order right before it completed and saving the logs to database.
func (p *Pool) serialSaveTradeLogs(order int, logs []common.TradeLog) error {
	var (
		logger = p.sugar.With(
			"func", "tradelogs/workers/Pool.serialSaveTradeLogs",
			"order", order,
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
			if err = p.storage.SaveTradeLogs(logs); err != nil {
				logger.Errorw("save trade logs into db failed",
					"err", err)
				p.mutex.Unlock()
				p.markAsFailed(order)
				return err
			}

			p.lastCompletedJobOrder++
			logger.Infow("save trade logs into db success")
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
	close(p.errCh)
}

// ErrCh returns error reporting channel of workers pool.
func (p *Pool) ErrCh() chan error {
	return p.errCh
}
