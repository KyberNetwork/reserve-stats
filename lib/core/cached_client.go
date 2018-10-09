package core

import (
	"errors"
	"math/big"
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

// Token returns the configured Token from core API.
// If token not found, trying to purging the cache and
// retry before returning an error.
func (cc *CachedClient) Token(address common.Address) (Token, error) {
	logger := cc.sugar.With(
		"func", "lib/core/CachedClient.Token",
		"address", address.Hex(),
	)
	cc.mu.RLock()
	token, ok := cc.cached[address]
	if ok {
		cc.mu.RUnlock()
		logger.Debug("cache hit")
		return token, nil
	}
	cc.mu.RUnlock()

	logger.Debug("cache miss, purging")
	_, err := cc.Tokens()
	if err != nil {
		return Token{}, err
	}

	cc.mu.RLock()
	defer cc.mu.RUnlock()
	token, ok = cc.cached[address]
	if ok {
		logger.Debug("cache hit after refreshing")
		return token, nil
	}
	logger.Debug("cache miss after refreshing")
	return Token{}, errors.New("token not found")
}

// FormatAmount formats the given amount in wei to human friendly
// number with preconfigured token decimals.
func (cc *CachedClient) FormatAmount(address common.Address, amount *big.Int) (float64, error) {
	token, err := cc.Token(address)
	if err != nil {
		return 0, err
	}
	return token.FormatAmount(amount), nil
}
