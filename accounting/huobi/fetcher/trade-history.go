package fetcher

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

type tradeHistoryFetcher func(string, time.Time, time.Time, ...huobi.ExtrasTradeHistoryParams) (huobi.TradeHistoryList, error)

func (fc *Fetcher) retry(fn tradeHistoryFetcher, symbol string, startTime, endTime time.Time, extras huobi.ExtrasTradeHistoryParams) (huobi.TradeHistoryList, error) {
	var (
		result huobi.TradeHistoryList
		err    error
		logger = fc.sugar.With("func", caller.GetCurrentFunctionName(),
			"symbol", symbol,
			"startTime", startTime,
			"endTime", endTime,
		)
	)
	for i := 0; i < fc.attempt; i++ {
		result, err = fn(symbol, startTime, endTime, extras)
		if err == nil {
			return result, nil
		}
		logger.Warnw("fail to fetch trade history", "error", err, "attempt", i+1)
		time.Sleep(fc.retryDelay)
	}
	return result, err
}

func (fc *Fetcher) getTradeHistoryWithSymbol(symbol string, from, to time.Time) ([]huobi.TradeHistory, error) {
	var (
		result []huobi.TradeHistory
		lastID string
		extras huobi.ExtrasTradeHistoryParams
	)
	latestTimeStored, err := fc.storage.GetLastStoredTimestamp(symbol)
	if err != nil {
		return result, err
	}
	if from.Before(latestTimeStored) {
		from = latestTimeStored
	}
	startTime := from
	for {
		endTime := startTime.Add(48 * time.Hour)
		if endTime.After(to) {
			endTime = to
		}
		for {
			tradeHistoriesResponse, err := fc.retry(fc.client.GetTradeHistory, symbol, startTime, endTime, extras)
			if err != nil {
				return result, err
			}
			// while result != empty, get trades latest time to toTime
			if len(tradeHistoriesResponse.Data) == 0 {
				break
			}
			//huobi returns tradelogs with latest trade in the beginning of the slice
			lastTrade := tradeHistoriesResponse.Data[0]
			if strconv.FormatInt(lastTrade.ID, 10) == lastID {
				break
			}
			lastID = strconv.FormatInt(lastTrade.ID, 10)

			result = append(result, tradeHistoriesResponse.Data...)

			startTime = timeutil.TimestampMsToTime(lastTrade.CreatedAt + 1)
			if endTime.Before(startTime) {
				break
			}
		}
		startTime = endTime
		if startTime == to {
			break
		}
	}

	return result, nil
}

//GetTradeHistory return all trade history between from-to and
func (fc *Fetcher) GetTradeHistory(from, to time.Time, symbols []huobi.Symbol) error {
	var (
		logger = fc.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", from,
			"to", to,
		)
		result      = make(map[string][]huobi.TradeHistory)
		fetchResult = sync.Map{}
		assertError error
		errGroup    errgroup.Group
		err         error
	)

	if len(symbols) == 0 {
		symbols, err = fc.client.GetSymbolsPair()
		if err != nil {
			return err
		}
	}
	for _, sym := range symbols {
		errGroup.Go(
			func(symbol string) func() error {
				return func() error {
					singleResult, err := fc.getTradeHistoryWithSymbol(symbol, from, to)
					if err != nil {
						return err
					}
					logger.Infow("Fetching done", "symbol", symbol, "error", err, "time", time.Now())
					if len(singleResult) > 0 {
						return fc.storage.UpdateTradeHistory(singleResult)
					}
					return nil
				}
			}(sym.SymBol),
		)
	}

	if err := errGroup.Wait(); err != nil {
		return err
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
	return assertError
}
