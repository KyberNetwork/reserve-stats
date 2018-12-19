package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// NewMockTokenAmountFormatter creates a new instance of MockTokenAmountFormatter, that implements blockchain.TokenAmountFormatterInterface, intended to use in tests.
func NewMockTokenAmountFormatter() *MockTokenAmountFormatter {
	return &MockTokenAmountFormatter{}
}

// MockTokenAmountFormatter is a simple mock client sending mock token from Core
type MockTokenAmountFormatter struct {
}

// ToWei return the given human friendly number to wei unit.
func (mc *MockTokenAmountFormatter) ToWei(_ common.Address, _ float64) (*big.Int, error) {
	return big.NewInt(100), nil
}

// FromWei formats the given amount in wei to human friendly
// number with preconfigured token decimals.
func (mc *MockTokenAmountFormatter) FromWei(_ common.Address, _ *big.Int) (float64, error) {
	return 9.99, nil
}
