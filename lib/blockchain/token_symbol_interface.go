package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
)

// TokenSymbolInterface interface return token address
type TokenSymbolInterface interface {
	Symbol(common.Address) (string, error)
}
