package fetcher

import (
	"sync"
	"time"

	"github.com/KyberNetwork/tokenrate"
	ethereum "github.com/ethereum/go-ethereum"
	ethereumCommon "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/client"
	rrstorage "github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	lbdCommon "github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/crawler"
)

//defaultStartingTime for reserve-rate fetcher is on 31-01-2018
const defaultStartingTime uint64 = 1517356800000

//Fetcher is the struct taking care of fetching reserve-rates for accounting
type Fetcher struct {
	sugar                 *zap.SugaredLogger
	storage               rrstorage.Interface
	crawler               *crawler.ReserveRatesCrawler
	ethUSDRateFetcher     tokenrate.ETHUSDRateProvider
	lastBlockResolver     *lastblockdaily.LastBlockResolver
	sleepTime             time.Duration
	retryDelayTime        time.Duration
	retryAttempts         int
	lastCompletedJobOrder uint64
	mutex                 *sync.Mutex
	failed                bool
	addressClient         client.Interface
}

//Option set the init behaviour of Fetcher
type Option func(fc *Fetcher)

//NewFetcher return a fetcher with options
func NewFetcher(sugar *zap.SugaredLogger,
	storage rrstorage.Interface,
	crawler *crawler.ReserveRatesCrawler,
	lastBlockResolver *lastblockdaily.LastBlockResolver,
	ethusdRate tokenrate.ETHUSDRateProvider,
	retryDelay, sleepTime time.Duration,
	retryAttempts int,
	addressClient client.Interface,
	options ...Option) (*Fetcher, error) {

	fetcher := &Fetcher{
		sugar:             sugar,
		storage:           storage,
		crawler:           crawler,
		lastBlockResolver: lastBlockResolver,
		retryDelayTime:    retryDelay,
		sleepTime:         sleepTime,
		retryAttempts:     retryAttempts,
		ethUSDRateFetcher: ethusdRate,
		mutex:             &sync.Mutex{},
		failed:            false,
		addressClient:     addressClient,
	}
	for _, opt := range options {
		opt(fetcher)
	}

	return fetcher, nil
}

func (fc *Fetcher) fetch(fromTime, toTime time.Time, addresses []ethereumCommon.Address) error {
	var (
		lastBlockErrCh = make(chan error)
		rateErrChn     = make(chan error)
		lastBlockBlCh  = make(chan lbdCommon.BlockInfo)
		wg             = &sync.WaitGroup{}
		logger         = fc.sugar.With("func", "accounting/reserve-rate/fetcher/Fetcher.fetch",
			"from", fromTime.String(),
			"to", toTime.String())
		jobOrder = fc.getLastCompletedJobOrder()
	)
	logger.Debugw("start fetching...", "last completed job order", jobOrder)
	go fc.lastBlockResolver.Run(fromTime, toTime, lastBlockBlCh, lastBlockErrCh)
	for {
		select {
		case err := <-lastBlockErrCh:
			if err == ethereum.NotFound {
				logger.Info("reached the end date")
				wg.Wait()
				logger.Info("all fetcher jobs are completed")
				return nil
			}
			return err
		case err := <-rateErrChn:
			if err != nil {
				fc.markAsFailed()
				return err
			}
		case blockInfo := <-lastBlockBlCh:
			wg.Add(1)
			jobOrder++
			go func(errCh chan error, blockInfo lbdCommon.BlockInfo, attempts int, jobOrder uint64) {
				defer wg.Done()
				logger.Debugw("A job has started", "job order", jobOrder, "block", blockInfo.Block)

				rates, rateErr := retryFetchTokenRate(attempts, fc.sugar, fc.crawler, blockInfo.Block, fc.retryDelayTime, addresses)
				if rateErr != nil {
					fc.markAsFailed()
					errCh <- rateErr
				}
				//TODO: parallel this
				ethUSDRate, err := retryFetchETHUSDRate(attempts, fc.sugar, fc.ethUSDRateFetcher, blockInfo.Timestamp, fc.retryDelayTime)
				if err != nil {
					fc.markAsFailed()
					errCh <- err
				}
				if err = fc.serialDataStore(blockInfo, rates, ethUSDRate, jobOrder); err != nil {
					fc.markAsFailed()
					errCh <- err
				}
				logger.Debugw("A job has fetched successfully", "job order", jobOrder, "block", blockInfo.Block)
			}(rateErrChn, blockInfo, fc.retryAttempts, jobOrder)
		}
	}
}

func retryFetchTokenRate(maxAttempt int,
	sugar *zap.SugaredLogger,
	rsvRateCrawler *crawler.ReserveRatesCrawler,
	block uint64,
	retryInterval time.Duration,
	addresses []ethereumCommon.Address) (map[string]map[string]float64, error) {
	var (
		result = make(map[string]map[string]float64)
		err    error
		logger = sugar.With("func", "accounting/reserve-rate/fetcher/retryFetchTokenRate", "block", block)
	)

	for i := 0; i < maxAttempt; i++ {
		rates, err := rsvRateCrawler.GetReserveRatesWithAddresses(addresses, block)
		if err != nil {
			logger.Debugw("failed to fetch reserve rate", "attempt", i, "error", err)
			time.Sleep(retryInterval)
			continue
		}
		for reserve := range rates {
			result[reserve] = make(map[string]float64)
			for pair := range rates[reserve] {
				result[reserve][pair] = rates[reserve][pair].BuyReserveRate
			}
		}
		return result, nil
	}

	return nil, err
}

func retryFetchETHUSDRate(maxAttempt int,
	sugar *zap.SugaredLogger,
	fetcher tokenrate.ETHUSDRateProvider,
	timestamp time.Time,
	retryInterval time.Duration) (float64, error) {
	var (
		result float64
		err    error
		logger = sugar.With("func", "accounting/reserve-rate/fetcher/retryFetchETHUSDRate", "time", timestamp.String())
	)

	for i := 0; i < maxAttempt; i++ {
		result, err = fetcher.USDRate(timestamp)
		if err == nil {
			//ETHUSD rate is defined as "how much ETH I have to sell to get one USD"
			//coingecko return how much USD will I have to sell to get one ETH
			//Hence we do the inverse here. rate expected to be !=0, if result is 0 it will panic
			return 1 / result, nil
		}
		logger.Debugw("failed to fetch ETH-USD rate", "attempt", i, "error", err)
		time.Sleep(retryInterval)
	}

	return 0, err
}
