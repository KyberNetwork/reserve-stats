package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"
)

//TODO implement this
func (tldb *TradeLogDB) GetWalletStats(from, to time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error) {
	return nil, nil
}
