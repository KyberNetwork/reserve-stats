package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

func (tldb *TradeLogDB) GetUserVolume(userAddress ethereum.Address, from, to time.Time, freq string) (map[uint64]common.UserVolume, error) {
	return nil, nil
}
