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
	//TODO: storage will be add in another PR
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

func (f *Fetcher) getTradeHistoryForOneSymBol(fromTime, toTime time.Time, symbol string) ([]binance.TradeHistory, error) {
	var (
		logger = f.sugar.With("func", "accounting/binance-fetcher.getTradeHistoryForOneSymbol")
		result = []binance.TradeHistory{}
	)
	startTime := fromTime
	endTime := toTime
	for {
		tradeHistoriesResponse, err := f.getTradeHistoryWithRetry(symbol, startTime, endTime)
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
		startTime = timeutil.TimestampMsToTime(lastTrade.Time + 1)
	}
	return result, nil
}

//GetTradeHistory get all trade history from trades for all token
func (f *Fetcher) GetTradeHistory(fromTime, toTime time.Time) error {
	var (
		tradeHistories sync.Map
		logger         = f.sugar.With("func", "accounting/binance-fetcher.getTradeHistory")
		errGroup       errgroup.Group
	)
	// get list token
	exchangeInfo, err := f.client.GetExchangeInfo()
	if err != nil {
		return err
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
						oneSymbolTradeHistory, err := f.getTradeHistoryForOneSymBol(fromTime, toTime, pair.Symbol)
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
			return err
		}
		index += f.batchSize
	}

	// log here for test get trade history without persistence storage
	tradeHistories.Range(func(key, value interface{}) bool {
		logger.Info("symbol", key, "history", value)
		return true
	})
	// TODO: save to storage
	return nil
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
		time.Sleep(f.retryDelay)
	}
	return withdrawHistory, err
}

//GetWithdrawHistory get all withdraw history in time range fromTime to toTime
func (f *Fetcher) GetWithdrawHistory(fromTime, toTime time.Time) error {
	var (
		result []binance.WithdrawHistory
		logger = f.sugar.With("func", "accounting/binance-fetcher.GetWithdrawHistory")
	)
	logger.Info("Start get withdraw history")
	startTime := fromTime
	endTime := toTime
	for {
		withdrawHistory, err := f.getWithdrawHistoryWithRetry(startTime, endTime)
		if err != nil {
			return err
		}
		if len(withdrawHistory.WithdrawList) == 0 {
			break
		}
		result = append(result, withdrawHistory.WithdrawList...)

		// set start equal to latest withdraw apply time + 1
		latestWithdraw := len(withdrawHistory.WithdrawList) - 1
		latestTimeStamp := withdrawHistory.WithdrawList[latestWithdraw].ApplyTime
		startTime = timeutil.TimestampMsToTime(latestTimeStamp + 1)
	}
	// log for test get withdraw history successfully
	logger.Debugw("withdraw history", "list", result)

	// TODO: save to storage
	return nil
}
