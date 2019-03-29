package client

import (
	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
)

// FixedAddresses is an implementation of client.Interface that always returns
// a pre supplied list of addresses.
type FixedAddresses struct {
	addresses []common.ReserveAddress
}

// NewFixedAddresses creates a FixedAddresses instance with given list of reserve addresses.
func NewFixedAddresses(addresses []string, resolv blockchain.ContractTimestampResolver) (*FixedAddresses, error) {
	var ras []common.ReserveAddress
	for i, address := range addresses {
		ra := common.ReserveAddress{
			ID:          uint64(i),
			Address:     ethereum.HexToAddress(address),
			Type:        common.Reserve,
			Description: "fixed address",
		}
		timestamp, err := resolv.Resolve(ra.Address)
		if err != nil {
			return nil, err
		}
		ra.Timestamp = timestamp
		ras = append(ras, ra)
	}

	return &FixedAddresses{addresses: ras}, nil
}

// ReserveAddresses returns all available reserve addresses.
func (fa *FixedAddresses) ReserveAddresses(_ ...common.AddressType) ([]common.ReserveAddress, error) {
	return fa.addresses, nil
}
