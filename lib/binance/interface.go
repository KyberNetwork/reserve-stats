package binance

import (
	"context"
)

//Interface is interface for binance api client
type Interface interface {
	GetTradeHistory(symbol string, fromID int64) (TradeHistory, error)
	GetWithdrawHistory(fromTime, toTime uint64) (WithdrawHistoryList, error)
}

//BinanceLimiter define an inteface for binance Limiter action
type BinanceLimiter interface {
	WaitN(context.Context, int) error
}