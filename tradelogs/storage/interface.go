package storage

import (
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// Interface represent a storage for TradeLogs data
type Interface interface {
	SaveTradeLogs([]common.TradeLog) error
	LoadTradeLogs(from, to time.Time) []common.TradeLog
}

// tokenAmountFormatter is the formatter used to format the amount from big number
// to float using preconfigured decimals. The intended implementation is from Core API client.
type tokenAmountFormatter interface {
	FormatAmount(address ethereum.Address, amount *big.Int) (float64, error)
}
