package fetcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"go.uber.org/zap"
)

//Fetcher is the struct to fetch/store Data from huobi
type Fetcher struct {
	sugar      *zap.SugaredLogger
	client     huobi.Interface
	retryDelay time.Duration
	attempt    int
}

//NewFetcher return a fetcher object
func NewFetcher(sugar *zap.SugaredLogger, client huobi.Interface, retryDelay time.Duration, attempt int) *Fetcher {
	return &Fetcher{
		sugar:      sugar,
		client:     client,
		retryDelay: retryDelay,
		attempt:    attempt,
	}
}

func retry(fn func(string, time.Time, time.Time) (huobi.TradeHistoryList, error), symbol string, startTime, endTime time.Time, attempt int, retryDelay time.Duration) (huobi.TradeHistoryList, error) {
	var (
		result huobi.TradeHistoryList
		err    error
	)
	for i := 0; i < attempt; i++ {
		result, err = fn(symbol, startTime, endTime)
		if err == nil {
			return result, nil
		}
		time.Sleep(retryDelay)
	}
	return result, err
}

func (fc *Fetcher) getTradeHistoryWithSymbol(symbol string, from, to time.Time, tradeHistory *sync.Map) error {
	var (
		startTime = from
		endTime   = to
		result    []huobi.TradeHistory
	)
	for {
		tradeHistoriesResponse, err := retry(fc.client.GetTradeHistory, symbol, startTime, endTime, fc.attempt, fc.retryDelay)
		if err != nil {
			return err
		}
		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse.Data) == 0 {
			break
		}
		result = append(result, tradeHistoriesResponse.Data...)
		lastTrade := tradeHistoriesResponse.Data[0]
		startTime = timeutil.Midnight(timeutil.TimestampMsToTime(lastTrade.CreateAt)).AddDate(0, 0, 1)
		if endTime.Before(startTime) {
			break
		}
	}
	if len(result) != 0 {
		tradeHistory.Store(symbol, result)
	}
	return nil
}

//GetTradeHistory return all trade history between from-to and
func (fc *Fetcher) GetTradeHistory(from, to time.Time) (map[string][]huobi.TradeHistory, error) {
	var (
		logger = fc.sugar.With(
			"func", "accounting/accounting-huobi-fetcher/GetTradeHistory",
			"from", from,
			"to", to,
		)
		result      = make(map[string][]huobi.TradeHistory)
		fetchResult = sync.Map{}
		assertError error
		fetchErr    error
	)

	symbols, err := fc.client.GetSymbolsPair()
	if err != nil {
		return result, err
	}
	wg := &sync.WaitGroup{}
	for _, sym := range symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			err := fc.getTradeHistoryWithSymbol(symbol, from, to, &fetchResult)
			if err != nil {
				fetchErr = err
			}
			logger.Debugw("Fetching done", "symbol", symbol, "error", err, "time", time.Now())
		}(sym.SymBol)
	}
	wg.Wait()

	if fetchErr != nil {
		return nil, fetchErr
	}
	fetchResult.Range(func(key, value interface{}) bool {
		symbol, ok := key.(string)
		if !ok {
			assertError = fmt.Errorf("cannot assert key (value: %v) of map result to symbol", key)
			return false
		}
		historyList, ok := value.([]huobi.TradeHistory)
		if !ok {
			assertError = fmt.Errorf("cannot assert value (value: %v) of map result to TradeHistoryList", historyList)
			return false
		}
		result[symbol] = historyList
		return true
	})
	return result, assertError
}
