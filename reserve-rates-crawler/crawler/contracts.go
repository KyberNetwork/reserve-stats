package crawler

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
)

type wrapperForReserveRate interface {
	GetReserveRate(block uint64, rsvAddr ethereum.Address, srcs, dest []ethereum.Address) ([]*big.Int, []*big.Int, error)
}
