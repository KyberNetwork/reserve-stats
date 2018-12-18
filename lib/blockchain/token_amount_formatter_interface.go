package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// TokenAmountFormaterInterface interface convert token amount from/to wei
type TokenAmountFormaterInterface interface {
	FromWei(common.Address, *big.Int) (float64, error)
	ToWei(common.Address, float64) (*big.Int, error)
}
