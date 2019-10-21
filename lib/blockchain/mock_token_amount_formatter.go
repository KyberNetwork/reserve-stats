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
func (mc *MockTokenAmountFormatter) ToWei(_ common.Address, amount float64) (*big.Int, error) {
	decimals := 18
	exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	result, _ := big.NewFloat(0).Mul(big.NewFloat(amount), big.NewFloat(0).SetInt(exp)).Int(nil)
	return result, nil
}

// FromWei formats the given amount in wei to human friendly
// number with preconfigured token decimals.
func (mc *MockTokenAmountFormatter) FromWei(_ common.Address, amount *big.Int) (float64, error) {
	if amount == nil {
		return 0, nil
	}
	decimals := 18
	floatAmount := new(big.Float).SetInt(amount)
	power := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10), big.NewInt(int64(decimals)), nil,
	))
	res := new(big.Float).Quo(floatAmount, power)
	result, _ := res.Float64()
	return result, nil
}
