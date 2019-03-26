package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

//Interface is inteface for binance storage
type Interface interface {
	UpdateTradeHistory([]binance.TradeHistory) error
	GetTradeHistory(fromTime, toTime time.Time) ([]binance.TradeHistory, error)
}
