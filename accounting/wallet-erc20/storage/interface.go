package storage

import (
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
)

//Interface defines required functions for wallet-erc20 storage
type Interface interface {
	UpdateERC20Transfers([]common.ERC20Transfer) error
	GetERC20Transfers(WalletAddr, TokenAddr ethereum.Address, from, to time.Time) ([]common.ERC20Transfer, error)
	GetLastStoredBlock(wallet ethereum.Address) (int, error)
}
