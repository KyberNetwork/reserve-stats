package binance

import (
	"context"
)

//Interface is interface for binance api client
type Interface interface {
	GetTradeHistory(symbol string, fromID int64) (TradeHistory, error)
	GetWithdrawHistory(fromTime, toTime uint64) (WithdrawHistoryList, error)
	GetExchangeInfo() (ExchangeInfo, error)
}

// Limiter is the resource limiter for accessing Binance API.
type Limiter interface {
	WaitN(context.Context, int) error
}
