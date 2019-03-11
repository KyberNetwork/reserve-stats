package fetcher

import (
	"log"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	batchSize         = 100
	defaultRetryDelay = 2 * time.Minute
	defaultAttempt    = 4
)

//Fetcher is a fetcher for get binance data
type Fetcher struct {
	sugar      *zap.SugaredLogger
	client     *binance.Client
	retryDelay time.Duration
	attempt    int
	//TODO: storage will be add in another PR
}

//NewFetcher return a new fetcher instance
func NewFetcher(sugar *zap.SugaredLogger, client *binance.Client) *Fetcher {
	return &Fetcher{
		sugar:      sugar,
		client:     client,
		retryDelay: defaultRetryDelay,
		attempt:    defaultAttempt,
	}
}

func (f *Fetcher) getTradeHistoryWithRetry(symbol string, startTime, endTime time.Time) ([]binance.TradeHistory, error) {
	var (
		tradeHistoriesResponse []binance.TradeHistory
		err                    error
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		tradeHistoriesResponse, err = f.client.GetTradeHistory(symbol, -1, startTime, endTime)
		if err == nil {
			return tradeHistoriesResponse, nil
		}
		time.Sleep(f.retryDelay)
	}
	return tradeHistoriesResponse, err
}

func (f *Fetcher) getTradeHistoryForOneSymBol(fromTime, toTime time.Time, symbol string,
	tradeHistories *sync.Map, wg *sync.WaitGroup) error {
	var (
		logger = f.sugar.With("func", "accounting/binance-fetcher.getTradeHistoryForOneSymbol")
		result = []binance.TradeHistory{}
	)
	startTime := fromTime
	endTime := toTime
	defer wg.Done()
	for {
		tradeHistoriesResponse, err := f.getTradeHistoryWithRetry(symbol, startTime, endTime)
		if err != nil {
			logger.Debugw("get trade history error", "symbol", symbol, "error", err)
			return err
		}
		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse) == 0 {
			break
		}
		logger.Debugw("trade history for", "symbol", symbol, "history", tradeHistoriesResponse)
		result = append(result, tradeHistoriesResponse...)
		lastTrade := tradeHistoriesResponse[len(tradeHistoriesResponse)-1]
		startTime = timeutil.TimestampMsToTime(lastTrade.Time + 1)
	}
	if len(result) != 0 {
		tradeHistories.Store(symbol, result)
	}
	return nil
}

//GetTradeHistory get all trade history from trades for all token
func (f *Fetcher) GetTradeHistory(fromTime, toTime time.Time) error {
	var (
		tradeHistories sync.Map
		logger         = f.sugar.With("func", "accounting/binance-fetcher.getTradeHistory")
	)
	// get list token
	exchangeInfo, err := f.client.GetExchangeInfo()
	if err != nil {
		return err
	}
	tokenPairs := exchangeInfo.Symbols
	wg := sync.WaitGroup{}
	index := 0
	for index < len(tokenPairs) {
		for count := 0; count < batchSize && index+count < len(tokenPairs); count++ {
			pair := tokenPairs[index+count]
			log.Println(index + count)
			wg.Add(1)
			logger.Debugw("token", "pair", pair.Symbol)
			go f.getTradeHistoryForOneSymBol(fromTime, toTime, pair.Symbol, &tradeHistories, &wg)
		}
		wg.Wait()
		index += batchSize
	}

	// log here for test get trade history without persistence storage
	tradeHistories.Range(func(key, value interface{}) bool {
		logger.Info("symbol", key, "history", value)
		return true
	})
	// TODO: save to storage
	return nil
}
