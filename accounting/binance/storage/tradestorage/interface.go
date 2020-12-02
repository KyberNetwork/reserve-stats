package tradestorage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
)

//Interface is inteface for binance storage
type Interface interface {
	UpdateTradeHistory([]binance.TradeHistory, string) error
	GetTradeHistory(fromTime, toTime time.Time) (map[string][]binance.TradeHistory, error)
	GetLastStoredID(symbol, account string) (uint64, error)

	GetLatestConvertToETHPrice() (uint64, error)
	GetNotETHTrades() (map[string][]binance.TradeHistory, error)
	GetTradeByTimestamp(symbol string, timestamp time.Time) (binance.TradeHistory, error)
	UpdateConvertToETHPrice(originalSymbols, symbols []string, price []float64, timestamp []uint64, originalTrade, trade []binance.TradeHistory) error
	GetConvertToETHPrice(fromTime, toTime uint64) ([]binance.ConvertToETHPrice, error)

	UpdateMarginTradeHistory([]binance.TradeHistory, string) error
	GetMarginTradeHistory(fromTime, toTime time.Time) (map[string][]binance.TradeHistory, error)
	GetLastStoredMarginTradeID(symbol, account string) (uint64, error)
}
