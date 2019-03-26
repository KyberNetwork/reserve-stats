package storage

import (
	"math/big"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

//Interface represent interface for accounting lsited token service
type Interface interface {
	CreateOrUpdate(tokens []common.ListedToken, blockNumber *big.Int) error
	GetTokens() ([]common.ListedToken, uint64, uint64, error)
}
