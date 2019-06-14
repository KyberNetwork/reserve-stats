package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"
)

//TODO implement this
func (tldb *TradeLogDB) GetTradeSummary(from, to time.Time, timezone int8) (map[uint64]*common.TradeSummary, error) {
	return nil, nil
}
