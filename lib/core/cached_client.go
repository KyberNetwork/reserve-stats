package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// CachedClient is the wrapper of Core Client with caching ability.
type CachedClient struct {
	*Client

	mu     sync.RWMutex
	cached map[common.Address]Token
}

// NewCachedClient creates a new core client instance.
func NewCachedClient(client *Client) *CachedClient {
	return &CachedClient{
		Client: client,
		cached: make(map[common.Address]Token),
	}
}

// Tokens purges the current cached tokens and fetching from Core server.
func (cc *CachedClient) Tokens() ([]Token, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.cached = make(map[common.Address]Token)

	tokens, err := cc.Client.Tokens()
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		cc.cached[common.HexToAddress(token.Address)] = token
	}
	return tokens, nil
}
