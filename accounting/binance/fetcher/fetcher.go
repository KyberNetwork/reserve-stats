package fetcher

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/tradestorage"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/marketdata"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	withdrawalTimeLimit = time.Hour * 24 * 90 // 90 days
	depositTimeLimit    = time.Hour * 24 * 90 // 90 days
)

//Fetcher is a fetcher for get binance data
type Fetcher struct {
	sugar            *zap.SugaredLogger
	client           *binance.Client
	retryDelay       time.Duration
	storage          tradestorage.Interface
	attempt          int
	batchSize        int
	marketDataClient *marketdata.Client
}

//NewFetcher return a new fetcher instance
func NewFetcher(sugar *zap.SugaredLogger, client *binance.Client, retryDelay time.Duration, attempt, batchSize int, storage tradestorage.Interface,
	accountName string, marketDataClient *marketdata.Client) *Fetcher {
	return &Fetcher{
		sugar:            sugar,
		client:           client,
		retryDelay:       retryDelay,
		attempt:          attempt,
		batchSize:        batchSize,
		storage:          storage,
		marketDataClient: marketDataClient,
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

func (f *Fetcher) getGetAggregatedTradesWithRetry(symbol string, startTime, endTime uint64) ([]binance.AggregatedTrade, error) {
	var (
		aggregatedTrades []binance.AggregatedTrade
		err              error
		logger           = f.sugar.With("func", caller.GetCurrentFunctionName(), "symbol", symbol)
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		aggregatedTrades, err = f.client.GetAggregatedTrades(symbol, startTime, endTime, 1) // we only want to get 1 nearest price
		switch err {
		case binance.ErrBadAPIKeyFormat, binance.ErrRejectedMBxKey:
			return nil, err
		case nil:
			return aggregatedTrades, nil
		default:
			logger.Warnw("get aggregated trade failed", "error", err, "attempt", attempt)
			time.Sleep(f.retryDelay)
		}
	}
	return aggregatedTrades, err
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
		logger.Infow("trade history for", "symbol", symbol, "history", tradeHistoriesResponse)
		result = append(result, tradeHistoriesResponse...)
		lastTrade := tradeHistoriesResponse[len(tradeHistoriesResponse)-1]
		fromID = lastTrade.ID + 1
	}
	return result, nil
}

// GetTradeHistory get all trade history from trades for all token and save them into database
func (f *Fetcher) GetTradeHistory(fromIDs map[string]uint64, tokenPairs []binance.Symbol, account string) error {
	var (
		logger   = f.sugar.With("func", caller.GetCurrentFunctionName())
		errGroup errgroup.Group
		// notETHTrade = make(map[*binance.Symbol][]binance.TradeHistory)
	)
	index := 0
	for index < len(tokenPairs) {
		for count := 0; count < f.batchSize && index+count < len(tokenPairs); count++ {
			pair := tokenPairs[index+count]
			errGroup.Go(
				func(pair binance.Symbol) func() error {
					return func() error {
						logger.Infow("token", "pair", pair.Symbol)
						oneSymbolTradeHistory, err := f.getTradeHistoryForOneSymBol(fromIDs[pair.Symbol], pair.Symbol)
						if err != nil {
							return err
						}
						return f.storage.UpdateTradeHistory(oneSymbolTradeHistory, account)
					}
				}(pair),
			)
		}
		if err := errGroup.Wait(); err != nil {
			return err
		}
		index += f.batchSize
	}
	return nil
}

// UpdateTradeNotETH ...
func (f *Fetcher) UpdateTradeNotETH(originalSymbol, symbol string, notETHTrades []binance.TradeHistory) error {
	var (
		logger = f.sugar.With(
			"func", caller.GetCurrentFunctionName(),
		)
		prices            []float64
		endTimes          []uint64
		trades, ethTrades []binance.TradeHistory
	)
	isPairValid, err := f.marketDataClient.PairSupported("binance", strings.ToLower(symbol))
	if err != nil {
		return err
	}
	if !isPairValid {
		logger.Infow("pair is not supported", "pair", symbol)
		return nil
	}
	for _, trade := range notETHTrades {
		endTime := trade.Time

		// get aggregated trade for that timestamp
		var (
			delta uint64 = 60 * 60 * 1000 // default 5 seconds
			res   []binance.AggregatedTrade
			err   error
		)
		startTime := endTime - delta + 1000
		res, err = f.getGetAggregatedTradesWithRetry(symbol, startTime, endTime)
		if err != nil {
			logger.Errorw("failed to get aggregated trades from binance", "error", err)
			return err
		}
		// increase delta if there is no result
		if len(res) == 0 {
			return fmt.Errorf("no trade in 1 hour")
		}
		price, err := strconv.ParseFloat(res[0].Price, 64)
		if err != nil {
			logger.Errorw("failed to parse price", "error", err)
			return err
		}
		// find the trade which is timestamp < endTime - 2mins
		timestampMillis := endTime - 2*60*1000 // 2 min in millisecond
		timestamp := timeutil.TimestampMsToTime(timestampMillis)
		ethTrade, err := f.storage.GetTradeByTimestamp(symbol, timestamp)
		if err != nil {
			logger.Errorw("failed to get trade by timestamp", "error", err)
			return err
		}
		prices = append(prices, price)
		endTimes = append(endTimes, endTime)
		trades = append(trades, trade)
		ethTrades = append(ethTrades, ethTrade)
	}
	// store info to persistent storage
	return f.storage.UpdateConvertToETHPrice(originalSymbol, symbol, prices, endTimes, trades, ethTrades)
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
			return withdrawHistory, nil
		}
		logger.Warnw("get withdraw history failed", "error", err, "attempt", attempt)
		time.Sleep(f.retryDelay)
	}
	return withdrawHistory, err
}

// GetWithdrawHistory get all withdraw history in time range fromTime to toTime
func (f *Fetcher) GetWithdrawHistory(fromTime, toTime time.Time) ([]binance.WithdrawHistory, error) {
	var (
		result []binance.WithdrawHistory
		logger = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	logger.Info("Start get withdraw history")
	endTime := toTime
	for toTime.After(fromTime) {
		endTime = fromTime.Add(withdrawalTimeLimit)
		if endTime.After(toTime) {
			endTime = toTime
		}
		withdrawHistory, err := f.getWithdrawHistoryWithRetry(fromTime, endTime)
		if err != nil {
			logger.Errorw("get withdraw history failed after retry", "attempts", f.attempt, "error", err)
			return result, err
		}
		result = append(result, withdrawHistory...)
		fromTime = endTime
	}
	// log for test get withdraw history successfully
	logger.Infow("withdraw history", "list", result)

	return result, nil
}

func (f *Fetcher) getMarginTradeHistoryWithRetry(symbol string, fromID uint64) ([]binance.TradeHistory, error) {
	var (
		tradeHistoriesResponse []binance.TradeHistory
		err                    error
		logger                 = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		tradeHistoriesResponse, err = f.client.GetMarginTradeHistory(symbol, fromID)
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

func (f *Fetcher) getMarginTradeHistoryForOneSymBol(fromID uint64, symbol string) ([]binance.TradeHistory, error) {
	var (
		logger = f.sugar.With("func", caller.GetCurrentFunctionName())
		result []binance.TradeHistory
	)
	for {
		tradeHistoriesResponse, err := f.getMarginTradeHistoryWithRetry(symbol, fromID)
		if err != nil {
			logger.Errorw("get trade history error", "symbol", symbol, "error", err)
			return result, err
		}
		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse) == 0 {
			break
		}
		logger.Infow("trade history for", "symbol", symbol, "history", tradeHistoriesResponse)
		result = append(result, tradeHistoriesResponse...)
		lastTrade := tradeHistoriesResponse[len(tradeHistoriesResponse)-1]
		fromID = lastTrade.ID + 1
	}
	return result, nil
}

//GetMarginTradeHistory get all trade history from trades for all token and save them into database
func (f *Fetcher) GetMarginTradeHistory(fromIDs map[string]uint64, tokenPairs []binance.Symbol, account string) error {
	var (
		logger   = f.sugar.With("func", caller.GetCurrentFunctionName())
		errGroup errgroup.Group
	)
	index := 0
	for index < len(tokenPairs) {
		for count := 0; count < f.batchSize && index+count < len(tokenPairs); count++ {
			pair := tokenPairs[index+count]
			errGroup.Go(
				func(pair binance.Symbol) func() error {
					return func() error {
						logger.Infow("token", "pair", pair.Symbol)
						oneSymbolTradeHistory, err := f.getMarginTradeHistoryForOneSymBol(fromIDs[pair.Symbol], pair.Symbol)
						if err != nil {
							return err
						}
						return f.storage.UpdateMarginTradeHistory(oneSymbolTradeHistory, account)
					}
				}(pair),
			)
		}
		if err := errGroup.Wait(); err != nil {
			return err
		}
		index += f.batchSize
	}
	return nil
}

func (f *Fetcher) getDepositHistoryWithRetry(startTime, endTime time.Time) ([]binance.DepositHistory, error) {
	var (
		depositHistory []binance.DepositHistory
		err            error
		logger         = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	for attempt := 0; attempt < f.attempt; attempt++ {
		logger.Infow("attempt to get deposit history", "attempt", attempt, "startTime", startTime, "endTime", endTime)
		depositHistory, err = f.client.GetDepositHistory(startTime, endTime)
		if err == nil {
			return depositHistory, nil
		}
		logger.Warnw("get deposit history failed", "error", err, "attempt", attempt)
		time.Sleep(f.retryDelay)
	}
	return depositHistory, err
}

// GetDepositHistory get all deposit history in time range fromTime to toTime
func (f *Fetcher) GetDepositHistory(fromTime, toTime time.Time) ([]binance.DepositHistory, error) {
	var (
		result []binance.DepositHistory
		logger = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	logger.Info("Start get deposit history")
	endTime := toTime
	for toTime.After(fromTime) {
		endTime = fromTime.Add(depositTimeLimit)
		if endTime.After(toTime) {
			endTime = toTime
		}
		depositHistory, err := f.getDepositHistoryWithRetry(fromTime, endTime)
		if err != nil {
			logger.Errorw("get deposit history failed after retry", "attempts", f.attempt, "error", err)
			return result, err
		}
		result = append(result, depositHistory...)
		fromTime = endTime
	}
	// log for test get deposit history successfully
	logger.Infow("deposit history", "list", result)

	return result, nil
}
