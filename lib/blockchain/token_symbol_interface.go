package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
)

// TokenSymbolResolver is the common interface of resolver that
// resolve the symbol of a ERC20 Ethereum Token.
type TokenSymbolResolver interface {
	Symbol(common.Address) (string, error)
}
