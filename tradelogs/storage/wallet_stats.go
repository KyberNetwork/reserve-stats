package storage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//WalletStats return stats of a wallet address from time to time by a frequency
func (is *InfluxStorage) WalletStats(walletAddr ethereum.Address, from, to uint64, freq string) (common.WalletStats, error) {
	var (
		walletStats common.WalletStats
		err         error
	)
	return walletStats, err
}
