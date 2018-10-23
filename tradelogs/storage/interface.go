package storage

import (
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	ethereum "github.com/ethereum/go-ethereum/common"
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	SaveTradeLogs(logs []common.TradeLog, rates []tokenrate.ETHUSDRate) error
	LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error)

	GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddr ethereum.Address) (map[string]float64, error)
}
