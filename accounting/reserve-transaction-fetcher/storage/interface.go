package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/common"
)

// ReserveTransactionStorage is the common interface of accounting-reserve-transaction persistent storage.
type ReserveTransactionStorage interface {
	StoreNormalTx([]common.NormalTx) error
	GetNormalTx(from time.Time, to time.Time) ([]common.NormalTx, error)

	StoreInternalTx([]common.InternalTx) error
	GetInternalTx(from time.Time, to time.Time) ([]common.InternalTx, error)

	StoreERC20Transfer([]common.ERC20Transfer) error
	GetERC20Transfer(from time.Time, to time.Time) ([]common.ERC20Transfer, error)
}
