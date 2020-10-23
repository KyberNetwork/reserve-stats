package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

//Interface defines a set of interface for huobi storage, which can be implemented by any DB
type Interface interface {
	UpdateTradeHistory(trade map[int64]huobi.TradeHistory, account string) error
	GetTradeHistory(from, to time.Time) (map[string][]huobi.TradeHistory, error)
	GetLastStoredTimestamp(account string) (time.Time, error)
}
