package binance

//Interface is interface for binance api client
type Interface interface {
	GetTradeHistory(symbol string, fromID int64) (TradeHistory, error)
	GetWithdrawHistory(fromTime, toTime uint64) (WithdrawHistoryList, error)
}
