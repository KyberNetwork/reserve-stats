package fetcher

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/common"
)

// TransactionFetcher is the ethereum interface of transaction fetchers.
// A transaction fetcher should supports 3 kind of transactions: normal, internal and ERC20.
type TransactionFetcher interface {
	NormalTx(addr ethereum.Address, from, to *big.Int) ([]common.NormalTx, error)
	InternalTx(addr ethereum.Address, from, to *big.Int) ([]common.InternalTx, error)
	ERC20Transfer(addr ethereum.Address, from, to *big.Int) ([]common.ERC20Transfer, error)
}
