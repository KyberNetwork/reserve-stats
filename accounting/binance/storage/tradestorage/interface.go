package tradestorage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

//Interface is inteface for binance storage
type Interface interface {
	UpdateTradeHistory([]binance.TradeHistory, string) error
	GetTradeHistory(fromTime, toTime time.Time) (map[string]binance.TradeHistory, error)
	GetLastStoredID(symbol, account string) (uint64, error)
}
