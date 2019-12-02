package influx

import (
	"errors"
	"strings"

	ethereum "github.com/ethereum/go-ethereum/common"
)

var (
	errSymbolNotFound = errors.New("token symbol not found")
)

// GetTokenSymbol return token symbol of given address
func (is *Storage) GetTokenSymbol(address string) (string, error) {
	if v, ok := is.tokenSymbol.Load(ethereum.HexToAddress(address).Hex()); ok {
		return v.(string), nil
	}
	return "", errSymbolNotFound
}

// UpdateTokens update token symbol
func (is *Storage) UpdateTokens(tokenAddresses, symbols []string) error {
	// update the tokenSymbol map
	for index, tokenAddress := range tokenAddresses {
		is.tokenSymbol.Store(ethereum.HexToAddress(tokenAddress).Hex(), strings.ToUpper(symbols[index]))
	}
	return nil
}
