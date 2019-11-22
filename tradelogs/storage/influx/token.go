package influx

import (
	"strings"

	ethereum "github.com/ethereum/go-ethereum/common"
)

var tokenSymbol = map[string]string{
	ethereum.HexToAddress("0x89d24a6b4ccb1b6faa2625fe562bdd9a23260359").Hex(): "SAI",
}

// GetTokenSymbol return token symbol of given address
func (is *Storage) GetTokenSymbol(address string) (string, error) {
	return tokenSymbol[ethereum.HexToAddress(address).Hex()], nil
}

// UpdateTokens update token symbol
func (is *Storage) UpdateTokens(tokenAddresses, symbols []string) error {
	/// update the tokenSymbol map
	for index, tokenAddress := range tokenAddresses {
		tokenSymbol[ethereum.HexToAddress(tokenAddress).Hex()] = strings.ToUpper(symbols[index])
	}
	return nil
}
