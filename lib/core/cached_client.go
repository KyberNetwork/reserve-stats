package core

import (
	"errors"
	"math"
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

// token returns the configured Token from core API.
// If token not found, trying to purging the cache and
// retry before returning an error.
func (cc *CachedClient) token(address common.Address) (Token, error) {
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
	token, err := cc.token(address)
	if err != nil {
		return 0, err
	}
	return token.FormatAmount(amount), nil
}

// ToWei return the given human friendly number to wei unit.
func (cc *CachedClient) ToWei(address common.Address, amount float64) (*big.Int, error) {
	token, err := cc.token(address)
	if err != nil {
		return big.NewInt(0), err
	}

	decimals := token.Decimals
	// 6 is our smallest precision,
	if decimals < 6 {
		return big.NewInt(int64(amount * math.Pow10(int(decimals)))), nil
	}

	result := big.NewInt(int64(amount * math.Pow10(6)))
	return result.Mul(result, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(decimals-6), nil)), nil
}
