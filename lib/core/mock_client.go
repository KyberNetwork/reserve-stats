package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// NewMockClient creates a new instance of MockClient, that implements core.Interface, intended to use in tests.
func NewMockClient() *MockClient {
	return &MockClient{
		tokens: []Token{
			{ID: "KNC", Name: "KyberNetwork", Address: "0xdd974D5C2e2928deA5F71b9825b8b646686BD200", Decimals: 18},
			{ID: "ZRX", Name: "0x", Address: "0xe41d2489571d322189246dafa5ebde1f4699f498", Decimals: 18},
		},
	}
}

// MockClient is a simple mock client sending mock token from Core
type MockClient struct {
	tokens []Token
}

// Tokens returns all tokens of MockClient.
func (mc *MockClient) Tokens() ([]Token, error) {
	return mc.tokens, nil
}

// GetActiveTokens return list of Active Tokens from core
func (mc *MockClient) GetActiveTokens() ([]Token, error) {
	return mc.tokens, nil
}

// GetInternalTokens returns list of Internal Token from core
func (mc *MockClient) GetInternalTokens() ([]Token, error) {
	return mc.tokens, nil
}

// ToWei return the given human friendly number to wei unit.
func (mc *MockClient) ToWei(_ common.Address, _ float64) (*big.Int, error) {
	return big.NewInt(100), nil
}

// FromWei formats the given amount in wei to human friendly
// number with preconfigured token decimals.
func (mc *MockClient) FromWei(_ common.Address, _ *big.Int) (float64, error) {
	return 9.99, nil
}
