package storage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	SaveTradeLogs([]common.TradeLog) error
	LoadTradeLogs(from, to time.Time) []common.TradeLog
}
