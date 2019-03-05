package huobi

import (
	"context"
	"time"
)

//Interface is for huobi api client
type Interface interface {
	GetTradeHistory(symbol string, startDate, endDate time.Time) (TradeHistoryList, error)
	GetWithdrawHistory(currency string, fromID uint64) (WithdrawHistoryList, error)
	GetSymbolsPair() ([]Symbol, error)
}

//Limiter define an inteface for huobi Limiter action
type Limiter interface {
	WaitN(context.Context, int) error
}
