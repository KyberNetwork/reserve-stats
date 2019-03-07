package storage

import (
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

//Interface defines a set of interface for huobi storage, which can be implemented by any DB
type Interface interface {
	UpdateTradeHistory(trade huobi.TradeHistory) error
}
