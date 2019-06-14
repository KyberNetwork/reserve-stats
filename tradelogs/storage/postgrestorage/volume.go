package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

//TODO implement this
func (tldb *TradeLogDB) GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time,
	frequency string) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}

//TODO implement this
func (tldb *TradeLogDB) GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address,
	fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}
