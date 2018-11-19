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
func (is *InfluxStorage) GetCountryStats(countryCode string, from, to time.Time) (map[uint64]*common.CountryStats, error) {
	var (
		err    error
		logger = is.sugar.With("country", countryCode,
			"fromTime", from, "toTime", to)
		timeFilter    = fmt.Sprintf("(time >='%s' AND time <= '%s')", from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
		countryFilter = fmt.Sprintf("(country='%s')", countryCode)
	)

	cmd := fmt.Sprintf(`
	SELECT eth_per_trade,
	total_eth_volume,total_trade,total_usd_amount,
	usd_per_trade,unique_addresses,new_unique_addresses, kyced
	FROM country_stats WHERE %s AND %s`, timeFilter, countryFilter)

	logger.Debugw("get country stats", "query", cmd)

	response, err := is.queryDB(is.influxClient, cmd)
	if err != nil {
		return nil, err
	}

	logger.Debugw("got result for country stats query", "response", response)
	result := map[uint64]*common.CountryStats{}
	if len(response) == 0 || len(response[0].Series) == 0 {
		return result, nil
	}
	return convertQueryResultToCountry(response[0].Series[0])
}

func convertQueryResultToCountry(row influxModel.Row) (map[uint64]*common.CountryStats, error) {
	result := make(map[uint64]*common.CountryStats)
	if len(row.Values) == 0 {
		return nil, nil
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToCountrySummary(v)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToCountrySummary(v []interface{}) (uint64, *common.CountryStats, error) {
	if len(v) != 9 {
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
	newUniqueAddr, err := influxdb.GetUint64FromInterface(v[7])
	if err != nil {
		return 0, nil, err
	}
	kyced, err := influxdb.GetUint64FromInterface(v[8])
	if err != nil {
		return 0, nil, err
	}
	return tsUint64, &common.CountryStats{
		TotalETHVolume:     ethVolume,
		TotalUSDVolume:     usdVolume,
		TotalTrade:         totalTrade,
		UniqueAddresses:    uniqueAddr,
		NewUniqueAddresses: newUniqueAddr,
		USDPerTrade:        usdPerTrade,
		ETHPerTrade:        ethPerTrade,
		KYCEDAddresses:     kyced,
	}, nil
}
