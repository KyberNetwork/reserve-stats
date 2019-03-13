package binance

import (
	"context"
	"time"
)

//Interface is interface for binance api client
type Interface interface {
	GetTradeHistory(symbol string, fromID uint64) (TradeHistory, error)
	GetWithdrawHistory(fromTime, toTime time.Time) (WithdrawHistoryList, error)
	GetExchangeInfo() (ExchangeInfo, error)
}

// Limiter is the resource limiter for accessing Binance API.
type Limiter interface {
	WaitN(context.Context, int) error
}
