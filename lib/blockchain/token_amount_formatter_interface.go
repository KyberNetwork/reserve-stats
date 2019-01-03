package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// TokenAmountFormatterInterface interface convert token amount from/to wei
type TokenAmountFormatterInterface interface {
	FromWei(common.Address, *big.Int) (float64, error)
	ToWei(common.Address, float64) (*big.Int, error)
	KNCAddr() common.Address
	ETHAddr() common.Address
	IsBurnable(token common.Address) bool
}
