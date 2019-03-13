package fetcher

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"golang.org/x/sync/errgroup"
)

type tradeHistoryFetcher func(string, time.Time, time.Time, ...huobi.ExtrasTradeHistoryParams) (huobi.TradeHistoryList, error)

func (fc *Fetcher) retry(fn tradeHistoryFetcher, symbol string, startTime, endTime time.Time, extras huobi.ExtrasTradeHistoryParams) (huobi.TradeHistoryList, error) {
	var (
		result huobi.TradeHistoryList
		err    error
		logger = fc.sugar.With("func", "accounting/huobi/fetcher/trade-history.retry")
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
		startTime = from
		endTime   = to
		result    []huobi.TradeHistory
		lastID    string
		extras    huobi.ExtrasTradeHistoryParams
	)
	for {
		if lastID != "" {
			extras = huobi.ExtrasTradeHistoryParams{
				From:   lastID,
				Direct: "prev",
			}
		} else {
			extras = huobi.ExtrasTradeHistoryParams{}
		}
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

		startTime = timeutil.TimestampMsToTime(lastTrade.CreateAt)
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
					logger.Infow("Fetching done", "symbol", symbol, "error", err, "time", time.Now())
					return nil
				}
			}(sym.SymBol),
		)
	}

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
