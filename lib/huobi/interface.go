package huobi

import "time"

//Interface is for huobi api client
type Interface interface {
	GetTradeHistory(symbol string, startDate, endDate time.Time) (TradeHistoryList, error)
	GetWithdrawHistory(currency string, fromID uint64) (WithdrawHistoryList, error)
}
