package storage

import (
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

// ReserveTransactionStorage is the common interface of accounting-reserve-transaction persistent storage.
type ReserveTransactionStorage interface {
	StoreNormalTx([]common.NormalTx) error
	GetNormalTx(from time.Time, to time.Time) ([]common.NormalTx, error)

	StoreInternalTx([]common.InternalTx) error
	GetInternalTx(from, to time.Time) ([]common.InternalTx, error)

	StoreERC20Transfer([]common.ERC20Transfer, string) error
	GetERC20Transfer(from, to time.Time) ([]common.ERC20Transfer, error)

	GetWalletERC20Transfers(WalletAddr, TokenAddr ethereum.Address, from, to time.Time) ([]common.ERC20Transfer, error)

	StoreLastInserted(ethereum.Address, *big.Int) error
	GetLastInserted(ethereum.Address) (*big.Int, error)
}
