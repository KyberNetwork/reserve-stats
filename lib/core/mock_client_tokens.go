package core

// MockClient is a simple mock client sending mock token from Core
type MockClient struct {
}

// GetActiveTokens return list of Active Tokens from core
func (mc *MockClient) GetActiveTokens() ([]Token, error) {
	token1 := Token{"KNC", "KyberNetwork", "0xdd974D5C2e2928deA5F71b9825b8b646686BD200", 18}
	token2 := Token{"ZRX", "0x", "0xe41d2489571d322189246dafa5ebde1f4699f498", 18}
	result := []Token{token1, token2}
	return result, nil
}

// GetInternalTokens returns list of Internal Token from core
func (mc *MockClient) GetInternalTokens() ([]Token, error) {
	token1 := Token{"KNC", "KyberNetwork", "0xdd974D5C2e2928deA5F71b9825b8b646686BD200", 18}
	token2 := Token{"ZRX", "0x", "0xe41d2489571d322189246dafa5ebde1f4699f498", 18}
	result := []Token{token1, token2}
	return result, nil
}
