package fetcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"

	"golang.org/x/sync/errgroup"
)

type withdrawalHistoryFetcher func(string, uint64) (huobi.WithdrawHistoryList, error)

func (fc *Fetcher) retryGetWithdrawal(fn withdrawalHistoryFetcher, symbol string, fromID uint64) (huobi.WithdrawHistoryList, error) {
	var (
		result huobi.WithdrawHistoryList
		err    error
		logger = fc.sugar.With("func", "accounting/huobi/fetcher/withdrawal-history.retryGetWithdrawal")
	)
	for i := 0; i < fc.attempt; i++ {
		result, err = fn(symbol, fromID)
		if err == nil {
			return result, nil
		}
		logger.Warnw("fail to fetch withdrawal history", "error", err, "attempt", i+1)

		time.Sleep(fc.retryDelay)
	}
	return result, err
}

func (fc *Fetcher) getWithdrawHistoryWithSymbol(symbol string, fromID uint64) ([]huobi.WithdrawHistory, error) {
	var (
		nextFromID = fromID
		result     []huobi.WithdrawHistory
	)
	for {
		tradeHistoriesResponse, err := fc.retryGetWithdrawal(fc.client.GetWithdrawHistory, symbol, nextFromID)
		if err != nil {
			return result, err
		}

		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse.Data) == 0 {
			break
		}
		result = append(result, tradeHistoriesResponse.Data...)
		lastWithdrawal := tradeHistoriesResponse.Data[0]
		nextFromID = lastWithdrawal.ID + 1
	}

	return result, nil
}

//GetWithdrawHistory return all trade history between fromID and latest withdrawal
func (fc *Fetcher) GetWithdrawHistory(fromID uint64) (map[string][]huobi.WithdrawHistory, error) {
	var (
		logger = fc.sugar.With(
			"func", "accounting/accounting-huobi-fetcher/GetWithdrawHistory",
			"from", fromID,
		)
		result      = make(map[string][]huobi.WithdrawHistory)
		fetchResult = sync.Map{}
		assertError error
		errGroup    errgroup.Group
	)

	symbols, err := fc.client.GetCurrencies()
	if err != nil {
		return result, err
	}
	for _, sym := range symbols {
		errGroup.Go(
			func(symbol string) func() error {
				return func() error {
					singleResult, err := fc.getWithdrawHistoryWithSymbol(symbol, fromID)
					if err != nil {
						return err
					}
					if len(singleResult) > 0 {
						fetchResult.Store(symbol, singleResult)
					}
					logger.Infow("Fetching done", "symbol", symbol, "error", err, "time", time.Now())
					return nil
				}
			}(sym),
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
		historyList, ok := value.([]huobi.WithdrawHistory)
		if !ok {
			assertError = fmt.Errorf("cannot assert value (value: %v) of map result to WithdrawHistoryList", historyList)
			return false
		}
		result[symbol] = historyList
		return true
	})
	return result, assertError
}
