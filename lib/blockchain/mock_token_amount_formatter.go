package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// NewMockTokenAmountFormater creates a new instance of MockTokenAmountFormater, that implements blockchain.TokenAmountFormaterInterface, intended to use in tests.
func NewMockTokenAmountFormater() *MockTokenAmountFormater {
	return &MockTokenAmountFormater{}
}

// MockTokenAmountFormater is a simple mock client sending mock token from Core
type MockTokenAmountFormater struct {
}

// ToWei return the given human friendly number to wei unit.
func (mc *MockTokenAmountFormater) ToWei(_ common.Address, _ float64) (*big.Int, error) {
	return big.NewInt(100), nil
}

// FromWei formats the given amount in wei to human friendly
// number with preconfigured token decimals.
func (mc *MockTokenAmountFormater) FromWei(_ common.Address, _ *big.Int) (float64, error) {
	return 9.99, nil
}
