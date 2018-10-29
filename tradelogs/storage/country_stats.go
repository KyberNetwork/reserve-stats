package storage

import (
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

//GetCountryStats return stats of a country from time to time by daily fred in provided timezone
func (is *InfluxStorage) GetCountryStats(countryCode string, timezone int64, fromTime, toTime time.Time) (map[uint64]common.CountryStats, error) {
	var (
		result map[uint64]common.CountryStats
		err    error
	)

	logger := is.sugar.With("country", countryCode,
		"fromTime", fromTime, "toTime", toTime, "timezone", timezone)

	q := fmt.Sprintf(`
		SELECT sum_amount from %s
		WHERE country = '%s' and time >= '%s' and time <= '%s'
	`, countryCode, countryCode, fromTime, toTime)

	res, err := is.queryDB(is.influxClient, q)
	if err != nil {
		logger.Error(err)
		return result, err
	}

	for _, row := range res[0].Series[0].Values {
		ts, stats, err := is.rowToCountryStats(row)
		if err != nil {
			return result, err
		}
		key := timeutil.TimeToTimestampMs(ts)
		result[key] = stats
	}

	return result, err
}
