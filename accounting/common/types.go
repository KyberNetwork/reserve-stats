package common

import (
	"encoding/json"
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

// AddressType is type of an address.
//go:generate stringer -type=AddressType -linecomment
type AddressType int

func (i *AddressType) UnmarshalJSON(input []byte) error {
	var val string
	if err := json.Unmarshal(input, &val); err != nil {
		return err
	}

	typ, ok := validAddressTypes[val]
	if !ok {
		return fmt.Errorf("invalid address type: '%s'", string(input))
	}
	*i = typ
	return nil
}

func (i *AddressType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

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

// MarshalJSON implements custom JSON marshaller for ReserveAddress to
// format timestamp in unix millis instead of RFC3339.
func (r *ReserveAddress) MarshalJSON() ([]byte, error) {
	type AliasReserveAddress ReserveAddress

	var ts *uint64
	if !r.Timestamp.IsZero() {
		millis := timeutil.TimeToTimestampMs(r.Timestamp)
		ts = &millis
	}

	return json.Marshal(struct {
		Timestamp *uint64 `json:"timestamp,omitempty"`
		*AliasReserveAddress
	}{
		AliasReserveAddress: (*AliasReserveAddress)(r),
		Timestamp:           ts,
	})
}

// UnmarshalJSON implements custom JSON unmarshaller for ReserveAddress to
// format timestamp in unix millis instead of RFC3339.
func (r *ReserveAddress) UnmarshalJSON(data []byte) error {
	type AliasReserveAddress ReserveAddress
	decoded := new(struct {
		Timestamp *uint64 `json:"timestamp,omitempty"`
		AliasReserveAddress
	})

	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	if decoded.Timestamp != nil {
		r.Timestamp = timeutil.TimestampMsToTime(*decoded.Timestamp)
	}
	r.ID = decoded.ID
	r.Address = decoded.Address
	r.Type = decoded.Type
	r.Description = decoded.Description
	return nil
}
