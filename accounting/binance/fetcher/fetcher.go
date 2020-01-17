package fetcher

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
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
func NewFetcher(sugar *zap.SugaredLogger, client *binance.Client, retryDelay time.Duration, attempt, batchSize int) *Fetcher {
	return &Fetcher{
		sugar:      sugar,
		client:     client,
		retryDelay: retryDelay,
		attempt:    attempt,
		batchSize:  batchSize,
	}
}

func (f *Fetcher) getTradeHistoryWithRetry(symbol string, fromID uint64) ([]binance.TradeHistory, error) {
	var (
		tradeHistoriesResponse []binance.TradeHistory
		err                    error
		logger                 = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		tradeHistoriesResponse, err = f.client.GetTradeHistory(symbol, fromID)
		switch err {
		case binance.ErrBadAPIKeyFormat, binance.ErrRejectedMBxKey:
			return nil, err
		case nil:
			return tradeHistoriesResponse, nil
		default:
			logger.Warnw("get trade history failed", "error", err, "attempt", attempt)
			time.Sleep(f.retryDelay)
		}
	}
	return tradeHistoriesResponse, err
}

func (f *Fetcher) getTradeHistoryForOneSymBol(fromID uint64, symbol string) ([]binance.TradeHistory, error) {
	var (
		logger = f.sugar.With("func", caller.GetCurrentFunctionName())
		result []binance.TradeHistory
	)
	for {
		tradeHistoriesResponse, err := f.getTradeHistoryWithRetry(symbol, fromID)
		if err != nil {
			logger.Errorw("get trade history error", "symbol", symbol, "error", err)
			return result, err
		}
		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse) == 0 {
			break
		}
		result = append(result, tradeHistoriesResponse...)
		lastTrade := tradeHistoriesResponse[len(tradeHistoriesResponse)-1]
		fromID = lastTrade.ID + 1
	}
	return result, nil
}

//GetTradeHistory get all trade history from trades for all token
func (f *Fetcher) GetTradeHistory(fromIDs map[string]uint64) ([]binance.TradeHistory, error) {
	var (
		tradeHistories sync.Map
		logger         = f.sugar.With("func", caller.GetCurrentFunctionName())
		errGroup       errgroup.Group
		result         []binance.TradeHistory
	)
	// get list token
	exchangeInfo, err := f.client.GetExchangeInfo()
	if err != nil {
		logger.Errorw("failed to get exchange info from binance", "error", err)
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
						logger.Infow("token",
							"pair", pair.Symbol,
							"fromID", fromIDs[pair.Symbol],
						)
						oneSymbolTradeHistory, err := f.getTradeHistoryForOneSymBol(fromIDs[pair.Symbol], pair.Symbol)
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
			logger.Errorw("failed to get trade history from binance", "error", err)
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
		logger          = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		logger.Infow("attempt to get withdraw history", "attempt", attempt, "startTime", startTime, "endTime", endTime)
		withdrawHistory, err = f.client.GetWithdrawalHistory(startTime, endTime)
		if err == nil {
			logger.Errorw("failed to get withdraw history from binance", "error", err)
			return withdrawHistory, nil
		}
		logger.Errorw("get withdraw history failed", "error", err, "attempt", attempt)
		time.Sleep(f.retryDelay)
	}
	return withdrawHistory, err
}

//GetWithdrawHistory get all withdraw history in time range fromTime to toTime
func (f *Fetcher) GetWithdrawHistory(fromTime, toTime time.Time) ([]binance.WithdrawHistory, error) {
	var (
		result []binance.WithdrawHistory
		logger = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	logger.Info("Start get withdraw history")
	withdrawHistory, err := f.getWithdrawHistoryWithRetry(fromTime, toTime)
	if err != nil {
		logger.Errorw("failed to get withdraw history from binance", "error", err)
		return result, err
	}
	result = append(result, withdrawHistory.WithdrawList...)
	// log for test get withdraw history successfully
	logger.Debugw("withdraw history", "list", result)

	return result, nil
}
