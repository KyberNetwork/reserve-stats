package storage

import (
	"math/big"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//Interface represent interface for accounting lsited token service
type Interface interface {
	CreateOrUpdate(tokens []common.ListedToken, blockNumber *big.Int, reserve ethereum.Address) error
	GetTokens(reserve string) ([]common.ListedToken, uint64, uint64, error)
}
