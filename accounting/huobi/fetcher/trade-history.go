package fetcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"golang.org/x/sync/errgroup"
)

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

func (fc *Fetcher) getTradeHistoryWithSymbol(symbol string, from, to time.Time) ([]huobi.TradeHistory, error) {
	var (
		startTime = from
		endTime   = to
		result    []huobi.TradeHistory
	)
	for {
		tradeHistoriesResponse, err := retry(fc.client.GetTradeHistory, symbol, startTime, endTime, fc.attempt, fc.retryDelay)
		if err != nil {
			return result, err
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

	return result, nil
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
		errGroup    errgroup.Group
	)

	symbols, err := fc.client.GetSymbolsPair()
	if err != nil {
		return result, err
	}
	for _, sym := range symbols {
		errGroup.Go(
			func(symbol string) func() error {
				return func() error {
					singleResult, err := fc.getTradeHistoryWithSymbol(symbol, from, to)
					if err != nil {
						return err
					}
					if len(singleResult) > 0 {
						fetchResult.Store(symbol, singleResult)
					}
					logger.Debugw("Fetching done", "symbol", symbol, "error", err, "time", time.Now())
					return nil
				}
			}(sym.SymBol),
		)
	}
	errGroup.Wait()

	if err := errGroup.Wait(); err != nil {
		return result, nil
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
