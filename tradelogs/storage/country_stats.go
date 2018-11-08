package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	influxModel "github.com/influxdata/influxdb/models"
)

//GetCountryStats return stats of a country from time to time by daily fred in provided timezone
func (is *InfluxStorage) GetCountryStats(countryCode string, from, to uint64) (map[uint64]*common.CountryStats, error) {
	var (
		err    error
		logger = is.sugar.With("country", countryCode,
			"fromTime", from, "toTime", to)
		timeFilter    = fmt.Sprintf("(time >=%d%s AND time <= %d%s)", from, timePrecision, to, timePrecision)
		countryFilter = fmt.Sprintf("(country='%s')", countryCode)
	)

	cmd := fmt.Sprintf("SELECT time,eth_per_trade,total_eth_volume,total_trade,total_usd_amount,usd_per_trade,unique_addresses FROM country_stats WHERE %s AND %s", timeFilter, countryFilter)
	logger.Debugw("get country stats", "query", cmd)

	response, err := is.queryDB(is.influxClient, cmd)
	if err != nil {
		return nil, err
	}

	logger.Debugw("got result for country stats query", "response", response)
	if len(response) == 0 || len(response[0].Series) == 0 {
		return nil, nil
	}
	return convertQueryResultToCountry(response[0].Series[0])
}

func convertQueryResultToCountry(row influxModel.Row) (map[uint64]*common.CountryStats, error) {
	result := make(map[uint64]*common.CountryStats)
	if len(row.Values) == 0 {
		return nil, nil
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToSummary(v)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToSummary(v []interface{}) (uint64, *common.CountryStats, error) {
	if len(v) != 7 {
		return 0, nil, errors.New("value fields is invalid in len")
	}
	timestampString, ok := v[0].(string)
	if !ok {
		return 0, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, nil, err
	}
	tsUint64 := timeutil.TimeToTimestampMs(ts)
	ethPerTrade, err := influxdb.GetFloat64FromInterface(v[1])
	if err != nil {
		return 0, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[2])
	if err != nil {
		return 0, nil, err
	}
	totalTrade, err := influxdb.GetUint64FromInterface(v[3])
	if err != nil {
		return 0, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[4])
	if err != nil {
		return 0, nil, err
	}
	usdPerTrade, err := influxdb.GetFloat64FromInterface(v[5])
	if err != nil {
		return 0, nil, err
	}
	uniqueAddr, err := influxdb.GetUint64FromInterface(v[6])
	if err != nil {
		return 0, nil, err
	}

	return tsUint64, &common.CountryStats{
		TotalETHVolume:  ethVolume,
		TotalUSDVolume:  usdVolume,
		TotalTrade:      totalTrade,
		UniqueAddresses: uniqueAddr,
		USDPerTrade:     usdPerTrade,
		ETHPerTrade:     ethPerTrade,
	}, nil
}
