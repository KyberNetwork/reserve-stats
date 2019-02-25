package postgresql

import (
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

// ReserveAddress represents a reserve address.
type ReserveAddress struct {
	ID          uint64    `json:"id" db:"id"`
	Address     string    `json:"address" db:"address"`
	Type        string    `json:"type" db:"type"`
	Description string    `json:"description" db:"description"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

// Common converts the database presentation of ReserveAddress to common type.
func (ra *ReserveAddress) Common() (*common.ReserveAddress, error) {
	addressType, ok := common.IsValidAddressType(ra.Type)
	if !ok {
		return nil, fmt.Errorf("unknown type: %s", ra.Type)
	}

	return &common.ReserveAddress{
		ID:          ra.ID,
		Address:     ethereum.HexToAddress(ra.Address),
		Type:        addressType,
		Description: ra.Description,
		Timestamp:   ra.Timestamp,
	}, nil
}
