package storage

import (
	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

// Interface is the common interface of reserve addresses backend storage.
type Interface interface {
	Create(address ethereum.Address, addressType common.AddressType, description string) (uint64, error)
	Get(id uint64) (*common.ReserveAddress, error)
	GetAll() ([]*common.ReserveAddress, int64, error)
	Update(id uint64, address ethereum.Address, addressType *common.AddressType, description string) error
}
