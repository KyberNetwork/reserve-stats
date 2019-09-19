package postgrestorage

import (
	"fmt"
	"time"

	appname "github.com/KyberNetwork/reserve-stats/app-names"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
)

// GetIntegrationVolume returns integration_volume and non_integration_volume groups by day
func (tldb *TradeLogDB) GetIntegrationVolume(fromTime, toTime time.Time) (map[uint64]*common.IntegrationVolume, error) {
	logger := tldb.sugar.With("from", fromTime, "to", toTime, "func",
		"tradelogs/storage/postgresql/TradeLogDB.GetIntegrationVolume")
	integrationQuery := fmt.Sprintf(
		`SELECT
			SUM(eth_amount * (CASE WHEN integration_app = '%[1]s' then 1 else 0 end)) as integration_volume,
			SUM(eth_amount * (CASE WHEN integration_app != '%[1]s' then 1 else 0 end)) as non_integration_volume,
			%[2]s AS time
		FROM "%[3]s" 
		WHERE timestamp >= $1 and timestamp < $2
		GROUP BY time`,
		appname.KyberSwapAppName, schema.BuildDateTruncField("day", 0), schema.TradeLogsTableName)
	logger.Debugw("prepare statement", "stmt", integrationQuery)
	fromTime = schema.RoundTime(fromTime, "day", 0)
	toTime = schema.RoundTime(toTime, "day", 0).Add(time.Hour * 24)
	var records []struct {
		Timestamp            time.Time `db:"time"`
		IntegrationVolume    float64   `db:"integration_volume"`
		NonIntegrationVolume float64   `db:"non_integration_volume"`
	}
	err := tldb.db.Select(&records, integrationQuery, fromTime, toTime)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		logger.Debugw("no trade found")
		return nil, nil
	}
	results := make(map[uint64]*common.IntegrationVolume)
	for _, r := range records {
		ts := timeutil.TimeToTimestampMs(r.Timestamp)
		results[ts] = &common.IntegrationVolume{
			KyberSwapVolume:    r.IntegrationVolume,
			NonKyberSwapVolume: r.NonIntegrationVolume,
		}
	}
	return results, nil
}
