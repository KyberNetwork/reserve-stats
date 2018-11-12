package storage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

//GetWalletStats return stats of a wallet address from time to time by a frequency
func (is *InfluxStorage) GetWalletStats(from, to uint64, walletAddr string) (map[uint64]common.WalletStats, error) {
	var (
		walletStats map[uint64]common.WalletStats
		err         error
	)
	return walletStats, err
}
