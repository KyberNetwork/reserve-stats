package core

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Interface represents a client o interact with KyberNetwork core APIs.
type Interface interface {
	Tokens() ([]Token, error)
	FromWei(common.Address, *big.Int) (float64, error)
	ToWei(common.Address, float64) (*big.Int, error)
}
