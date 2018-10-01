package setting

import "github.com/KyberNetwork/reserve-stats/common"

// Client is a simple http client to query for setting from Core
type Client struct {
	coreURL string
}

// NewSettingClient returns a new SettingClient pointer
func NewSettingClient(coreURL string) (*Client, error) {
	return &Client{}, nil
}

// GetActiveTokens return list of Active Tokens from core
func (sc *Client) GetActiveTokens() ([]common.Token, error) {
	token1 := common.NewToken("KNC", "KyberNetwork", "0xdd974D5C2e2928deA5F71b9825b8b646686BD200", 18, true, true, 1535021910158)
	token2 := common.NewToken("ZRX", "0x", "0xe41d2489571d322189246dafa5ebde1f4699f498", 18, true, true, 1537555947208)
	result := []common.Token{token1, token2}
	return result, nil
}

// GetInternalTokens returns list of Internal Token from core
func (sc *Client) GetInternalTokens() ([]common.Token, error) {
	token1 := common.NewToken("KNC", "KyberNetwork", "0xdd974D5C2e2928deA5F71b9825b8b646686BD200", 18, true, true, 1535021910158)
	token2 := common.NewToken("ZRX", "0x", "0xe41d2489571d322189246dafa5ebde1f4699f498", 18, true, true, 1537555947208)
	result := []common.Token{token1, token2}
	return result, nil
}
