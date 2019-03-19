package fetcher

import (
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"golang.org/x/sync/errgroup"
)

//Fetcher is a fetcher for get binance data
type Fetcher struct {
	sugar      *zap.SugaredLogger
	client     *binance.Client
	retryDelay time.Duration
	attempt    int
	batchSize  int
}

//NewFetcher return a new fetcher instance
func NewFetcher(sugar *zap.SugaredLogger, client *binance.Client, retryDelay, attempt, batchSize int) *Fetcher {
	retryDelayTime := time.Duration(retryDelay) * time.Minute
	return &Fetcher{
		sugar:      sugar,
		client:     client,
		retryDelay: retryDelayTime,
		attempt:    attempt,
		batchSize:  batchSize,
	}
}

func (f *Fetcher) getTradeHistoryWithRetry(symbol string, fromID uint64) ([]binance.TradeHistory, error) {
	var (
		tradeHistoriesResponse []binance.TradeHistory
		err                    error
		logger                 = f.sugar.With("func", "accounting/binance-fetcher/fetcher.getTradeHistoryWithRetry")
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		tradeHistoriesResponse, err = f.client.GetTradeHistory(symbol, fromID)
		if err == nil {
			return tradeHistoriesResponse, nil
		}
		logger.Warnw("get trade history failed", "error", err, "attempt", attempt)
		time.Sleep(f.retryDelay)
	}
	return tradeHistoriesResponse, err
}

func (f *Fetcher) getTradeHistoryForOneSymBol(fromID uint64, symbol string) ([]binance.TradeHistory, error) {
	var (
		logger = f.sugar.With("func", "accounting/binance-fetcher.getTradeHistoryForOneSymbol")
		result = []binance.TradeHistory{}
	)
	for {
		tradeHistoriesResponse, err := f.getTradeHistoryWithRetry(symbol, fromID)
		if err != nil {
			logger.Debugw("get trade history error", "symbol", symbol, "error", err)
			return result, err
		}
		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse) == 0 {
			break
		}
		logger.Debugw("trade history for", "symbol", symbol, "history", tradeHistoriesResponse)
		result = append(result, tradeHistoriesResponse...)
		lastTrade := tradeHistoriesResponse[len(tradeHistoriesResponse)-1]
		fromID = lastTrade.ID + 1
	}
	return result, nil
}

//GetTradeHistory get all trade history from trades for all token
func (f *Fetcher) GetTradeHistory(fromID uint64) ([]binance.TradeHistory, error) {
	var (
		tradeHistories sync.Map
		logger         = f.sugar.With("func", "accounting/binance-fetcher.getTradeHistory")
		errGroup       errgroup.Group
		result         []binance.TradeHistory
	)
	// get list token
	exchangeInfo, err := f.client.GetExchangeInfo()
	if err != nil {
		return result, err
	}
	tokenPairs := exchangeInfo.Symbols
	index := 0
	for index < len(tokenPairs) {
		for count := 0; count < f.batchSize && index+count < len(tokenPairs); count++ {
			pair := tokenPairs[index+count]
			errGroup.Go(
				func(pair binance.Symbol) func() error {
					return func() error {
						logger.Debugw("token", "pair", pair.Symbol)
						oneSymbolTradeHistory, err := f.getTradeHistoryForOneSymBol(fromID, pair.Symbol)
						if err != nil {
							return err
						}
						if len(oneSymbolTradeHistory) != 0 {
							tradeHistories.Store(pair.Symbol, oneSymbolTradeHistory)
						}
						return nil
					}
				}(pair),
			)
		}
		if err := errGroup.Wait(); err != nil {
			return result, err
		}
		index += f.batchSize
	}

	// log here for test get trade history without persistence storage
	tradeHistories.Range(func(key, value interface{}) bool {
		tradeHistory, ok := value.([]binance.TradeHistory)
		if !ok {
			return false
		}
		if len(tradeHistory) != 0 {
			result = append(result, tradeHistory...)
		}
		return true
	})
	return result, nil
}

func (f *Fetcher) getWithdrawHistoryWithRetry(startTime, endTime time.Time) (binance.WithdrawHistoryList, error) {
	var (
		withdrawHistory binance.WithdrawHistoryList
		err             error
		logger          = f.sugar.With(
			"func", "accounting/binance-fetcher.getWithdrawHistoryWithRetry",
		)
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		logger.Debugw("attempt to get withdraw history", "attempt", attempt, "startTime", startTime, "endTime", endTime)
		withdrawHistory, err = f.client.GetWithdrawalHistory(startTime, endTime)
		if err == nil {
			return withdrawHistory, nil
		}
		logger.Warnw("get withdraw history failed", "error", err, "attempt", attempt)
		time.Sleep(f.retryDelay)
	}
	return withdrawHistory, err
}

func appendResult(result map[string]binance.WithdrawHistory, withdrawList []binance.WithdrawHistory) time.Time {
	var (
		startTime time.Time
	)
	for _, withdrawHistory := range withdrawList {
		if _, exist := result[withdrawHistory.ID]; !exist {
			result[withdrawHistory.ID] = withdrawHistory
			withdrawHistoryTime := timeutil.TimestampMsToTime(withdrawHistory.ApplyTime)
			if withdrawHistoryTime.After(startTime) {
				startTime = withdrawHistoryTime
			}
		}
	}
	if startTime.IsZero() {
		startTime = time.Now()
	}
	return startTime
}

//GetWithdrawHistory get all withdraw history in time range fromTime to toTime
func (f *Fetcher) GetWithdrawHistory(fromTime, toTime time.Time) (map[string]binance.WithdrawHistory, error) {
	var (
		result = make(map[string]binance.WithdrawHistory)
		logger = f.sugar.With("func", "accounting/binance-fetcher.GetWithdrawHistory")
	)
	logger.Info("Start get withdraw history")
	startTime := fromTime
	endTime := toTime
	for {
		withdrawHistory, err := f.getWithdrawHistoryWithRetry(startTime, endTime)
		if err != nil {
			return result, err
		}
		startTime = appendResult(result, withdrawHistory.WithdrawList)
		if startTime.After(endTime) {
			break
		}
	}
	// log for test get withdraw history successfully
	logger.Debugw("withdraw history", "list", result)

	return result, nil
}
