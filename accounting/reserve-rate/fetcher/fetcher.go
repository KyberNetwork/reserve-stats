package fetcher

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/KyberNetwork/tokenrate"

	rrstorage "github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/crawler"

	"github.com/ethereum/go-ethereum"
	"go.uber.org/zap"
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
	fromTime              time.Time
	toTime                time.Time
	sleepTime             time.Duration
	retryDelayTime        time.Duration
	retryAttempts         int
	lastCompletedJobOrder uint64
	mutex                 *sync.Mutex
	failed                bool
}

//Option set the init behaviour of Fetcher
type Option func(fc *Fetcher)

//WithFromTime set Fromtime for fetcher with a blockInfo of input timestamp and unknowm blockNumber
//without this init function Fetcher will
func WithFromTime(from time.Time) Option {
	return func(fc *Fetcher) {
		fc.fromTime = from
	}
}

//WithToTime set toblock for fetcher with a blockInfo of inputt timestamp and unknown blockNumber.
//Without this init function, Fetcher is put into daemon mode
func WithToTime(to time.Time) Option {
	return func(fc *Fetcher) {
		fc.toTime = to
	}
}

//NewFetcher return a fetcher with options
func NewFetcher(sugar *zap.SugaredLogger,
	storage rrstorage.Interface,
	crawler *crawler.ReserveRatesCrawler,
	lastBlockResolver *lastblockdaily.LastBlockResolver,
	ethusdRate tokenrate.ETHUSDRateProvider,
	retryDelay, sleepTime time.Duration,
	retryAttempts int,
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
	}
	for _, opt := range options {
		opt(fetcher)
	}
	//get last crawled blockInfo if fetcher is init without  from tim
	if fetcher.fromTime.IsZero() {
		sugar.Debugw("empty from time, trying to get from block from db...")
		fromBlockInfo, err := fetcher.storage.GetLastResolvedBlockInfo()
		if err == sql.ErrNoRows {
			fetcher.fromTime = timeutil.TimestampMsToTime(defaultStartingTime)
			sugar.Debugw("There is no row from DB, running from default from time", "from time", fetcher.fromTime.String())
		} else if err != nil {
			return nil, fmt.Errorf("cannot get last resolved block info from db, err: %v", err)
		} else if err == nil {
			fetcher.lastBlockResolver.LastResolved = fromBlockInfo
		}
	}

	return fetcher, nil
}

func (fc *Fetcher) fetch(fromTime, toTime time.Time) error {
	var (
		lastBlockErrCh = make(chan error)
		rateErrChn     = make(chan error)
		lastBlockBlCh  = make(chan lastblockdaily.BlockInfo)
		wg             = &sync.WaitGroup{}
		logger         = fc.sugar.With("function", "accounting/reserve-rate/fetcher/fetcher.Fetch()",
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
			go func(errCh chan error, blockInfo lastblockdaily.BlockInfo, attempts int, jobOrder uint64) {
				defer wg.Done()
				logger.Debugw("A job has started", "job order", jobOrder, "block", blockInfo.Block)

				rates, rateErr := retryFetchTokenRate(attempts, fc.sugar, fc.crawler, blockInfo.Block, fc.retryDelayTime)
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

//Run start the fetcher
func (fc *Fetcher) Run() error {
	var (
		toTime     = fc.toTime
		fromTime   = fc.fromTime
		logger     = fc.sugar.With("function", "accounting/reserve-rate/fetcher.Run")
		daemonMode = false
	)
	if fc.toTime.IsZero() {
		toTime = time.Now()
		logger.Info("no end time specified, fetcher run in daemon mode")
		daemonMode = true
	}
	for {
		if err := fc.fetch(fromTime, toTime); err != nil {
			return err
		}
		if !daemonMode {
			break
		}
		time.Sleep(fc.sleepTime)
		fromTime = toTime
		toTime = time.Now()
	}
	return nil
}

func retryFetchTokenRate(maxAttempt int,
	sugar *zap.SugaredLogger,
	rsvRateCrawler *crawler.ReserveRatesCrawler,
	block uint64,
	retryInterval time.Duration) (map[string]map[string]float64, error) {
	var (
		result = make(map[string]map[string]float64)
		err    error
		logger = sugar.With("function", "accounting/reserve-rate/fetcher/fetcher.retryFetchTokenRate", "block", block)
	)

	for i := 0; i < maxAttempt; i++ {
		// TODO: GetReserveRates might fetch too more data than we need here?
		rates, err := rsvRateCrawler.GetReserveRates(block)
		if err != nil {
			logger.Debugw("failed to fetch reserve rate", "attempt", i, "error", err)
			time.Sleep(retryInterval)
			continue
		}
		for k1 := range rates {
			result[k1] = make(map[string]float64)
			for k2 := range rates[k1] {
				// TODO: is it really buy reserve rate?
				result[k1][k2] = rates[k1][k2].BuyReserveRate
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
		logger = sugar.With("function", "accounting/reserve-rate/fetcher/fetcher.retryFetchETHUSDRate", "time", timestamp.String())
	)

	for i := 0; i < maxAttempt; i++ {
		result, err = fetcher.USDRate(timestamp)
		if err == nil {
			return result, nil
		}
		logger.Debugw("failed to fetch ETH-USD rate", "attempt", i, "error", err)
		time.Sleep(retryInterval)
	}

	return result, err
}
