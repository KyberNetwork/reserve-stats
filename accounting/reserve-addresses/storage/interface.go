package storage

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

// Interface is the common interface of reserve addresses backend storage.
type Interface interface {
	Create(address ethereum.Address, addressType common.AddressType, description string, ts time.Time) (uint64, error)
}
