package common

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// AddressType is type of an address.
//go:generate stringer -type=AddressType -linecomment
type AddressType int

const (
	// Reserve is address of KyberNetwork's reserve contract.
	Reserve AddressType = iota // reserve
	// PricingOperator is operator addresses of conversion rates contract.
	PricingOperator // pricing_operator
	// SanityOperator is operator address of sanity rates contract.
	SanityOperator // sanity_operator
	// IntermediateOperator is the proxy address used when withdraw, deposit to centralized exchanges.
	// It is necessary as some exchanges don't allow withdraw, deposit directly to a contract account.
	IntermediateOperator // intermediate_operator
	// CEXDepositAddress is Ethereum address used to deposit funds to centralized exchanges.
	CEXDepositAddress // cex_deposit_address
)

var validAddressTypes = map[string]AddressType{
	Reserve.String():              Reserve,
	PricingOperator.String():      PricingOperator,
	SanityOperator.String():       SanityOperator,
	IntermediateOperator.String(): IntermediateOperator,
	CEXDepositAddress.String():    CEXDepositAddress,
}

// IsValidAddressType returns true if given
func IsValidAddressType(typ string) (AddressType, bool) {
	addressType, ok := validAddressTypes[typ]
	return addressType, ok
}

// ReserveAddress represents a reserve address.
type ReserveAddress struct {
	ID          uint64           `json:"id"`
	Address     ethereum.Address `json:"address"`
	Type        AddressType      `json:"type"`
	Description string           `json:"description"`
	Timestamp   time.Time        `json:"timestamp"`
}
