package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"
)

//TODO: implement this
func (tldb *TradeLogDB) GetIntegrationVolume(fromTime, toTime time.Time) (map[uint64]*common.IntegrationVolume, error) {
	return nil, nil
}
