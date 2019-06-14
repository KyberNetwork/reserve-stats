package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"
)

//TODO: implement this
func (tldb *TradeLogDB) GetCountryStats(countryCode string, from, to time.Time, timezone int8) (map[uint64]*common.CountryStats, error) {
	return nil, nil
}
