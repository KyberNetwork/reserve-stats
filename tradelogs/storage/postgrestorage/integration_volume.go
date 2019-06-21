package postgrestorage

import (
	appname "github.com/KyberNetwork/reserve-stats/app-names"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	"time"
)

func (tldb *TradeLogDB) GetIntegrationVolume(fromTime, toTime time.Time) (map[uint64]*common.IntegrationVolume, error) {
	logger := tldb.sugar.With("from", fromTime, "to", toTime, "func",
		"tradelogs/storage/postgresql/TradeLogDB.GetIntegrationVolume")

	integrationQuery :=
		`SELECT
			SUM(eth_amount * (CASE WHEN integration_app = '` + appname.KyberSwapAppName + `' then 1 else 0 end)) as integration_volume,
			SUM(eth_amount * (CASE WHEN integration_app != '` + appname.KyberSwapAppName + `' then 1 else 0 end)) as non_integration_volume,
			` + schema.BuildDateTruncField("day", 0) + ` AS time
		FROM ` + schema.TradeLogsTableName + ` WHERE timestamp >= $1 and timestamp < $2
		GROUP BY time
		`
	fromTime = schema.RoundTime(fromTime, "day", 0)
	toTime = schema.RoundTime(toTime, "day", 0).Add(time.Hour * 24)

	var datas []struct {
		Timestamp            time.Time `db:"time"`
		IntegrationVolume    float64   `db:"integration_volume"`
		NonIntegrationVolume float64   `db:"non_integration_volume"`
	}

	logger.Debugw("prepare statement", "stmt", integrationQuery)
	err := tldb.db.Select(&datas, integrationQuery, fromTime.UTC().Format(schema.DefaultDateFormat),
		toTime.UTC().Format(schema.DefaultDateFormat))
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		logger.Debugw("no trade found")
		return nil, nil
	}
	results := make(map[uint64]*common.IntegrationVolume)
	for _, data := range datas {
		ts := timeutil.TimeToTimestampMs(data.Timestamp)
		results[ts] = &common.IntegrationVolume{
			KyberSwapVolume:    data.IntegrationVolume,
			NonKyberSwapVolume: data.NonIntegrationVolume,
		}
	}
	return results, nil
}
