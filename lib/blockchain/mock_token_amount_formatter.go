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

//KNCAddr return mock knc addr
func (mc *MockTokenAmountFormatter) KNCAddr() common.Address {
	return common.HexToAddress("0xdd974D5C2e2928deA5F71b9825b8b646686BD200")
}

//ETHAddr return mock eth addr
func (mc *MockTokenAmountFormatter) ETHAddr() common.Address {
	return common.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
}

//IsBurnable always return false
func (mc *MockTokenAmountFormatter) IsBurnable(token common.Address) bool {
	return false
}
