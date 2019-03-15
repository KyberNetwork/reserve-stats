package binancestorage

import (
	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

//Interface is inteface for binance storage
type Interface interface {
	UpdateTradeHistory(map[string]binance.TradeHistory) error
}
