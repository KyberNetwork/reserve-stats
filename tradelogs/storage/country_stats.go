package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	countryStatSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/country_stats"
	influxModel "github.com/influxdata/influxdb/models"
)

//GetCountryStats return stats of a country from time to time by daily fred in provided timezone
func (is *InfluxStorage) GetCountryStats(countryCode string, from, to time.Time, timezone int8) (map[uint64]*common.CountryStats, error) {
	var (
		err    error
		logger = is.sugar.With("country", countryCode,
			"fromTime", from, "toTime", to)
		timeFilter      = fmt.Sprintf("(time >='%s' AND time <= '%s')", from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
		countryFilter   = fmt.Sprintf("(country='%s')", countryCode)
		measurementName = common.CountryStatsMeasurementName
	)
	measurementName = getMeasurementName(measurementName, timezone)

	cmd := fmt.Sprintf(
		`SELECT %[1]s, %[2]s, %[3]s, %[4]s, %[5]s, %[6]s, %[7]s, %[8]s, %[9]s FROM %[10]s WHERE %[11]s AND %[12]s`,
		countryStatSchema.ETHPerTrade.String(),
		countryStatSchema.TotalETHVolume.String(),
		countryStatSchema.TotalTrade.String(),
		countryStatSchema.TotalUSDAmount.String(),
		countryStatSchema.USDPerTrade.String(),
		countryStatSchema.UniqueAddresses.String(),
		countryStatSchema.NewUniqueAddresses.String(),
		countryStatSchema.KYCedAddresses.String(),
		countryStatSchema.TotalBurnFee.String(),
		measurementName,
		timeFilter,
		countryFilter)

	logger.Debugw("get country stats", "query", cmd)

	response, err := influxdb.QueryDB(is.influxClient, cmd, is.dbName)
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
	idxs, err := countryStatSchema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return nil, nil
	}
	for _, v := range row.Values {
		ts, vol, err := convertRowValueToCountrySummary(v, idxs)
		if err != nil {
			return nil, err
		}
		result[ts] = vol
	}
	return result, nil
}

func convertRowValueToCountrySummary(v []interface{}, idxs map[countryStatSchema.FieldName]int) (uint64, *common.CountryStats, error) {
	if len(v) != 10 {
		return 0, nil, errors.New("value fields is invalid in len")
	}
	timestampString, ok := v[idxs[countryStatSchema.Time]].(string)
	if !ok {
		return 0, nil, errCantConvert
	}
	ts, err := time.Parse(time.RFC3339, timestampString)
	if err != nil {
		return 0, nil, err
	}
	tsUint64 := timeutil.TimeToTimestampMs(ts)
	ethPerTrade, err := influxdb.GetFloat64FromInterface(v[idxs[countryStatSchema.ETHPerTrade]])
	if err != nil {
		return 0, nil, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[idxs[countryStatSchema.TotalETHVolume]])
	if err != nil {
		return 0, nil, err
	}
	totalTrade, err := influxdb.GetUint64FromInterface(v[idxs[countryStatSchema.TotalTrade]])
	if err != nil {
		return 0, nil, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[idxs[countryStatSchema.TotalUSDAmount]])
	if err != nil {
		return 0, nil, err
	}
	usdPerTrade, err := influxdb.GetFloat64FromInterface(v[idxs[countryStatSchema.USDPerTrade]])
	if err != nil {
		return 0, nil, err
	}
	uniqueAddr, err := influxdb.GetUint64FromInterface(v[idxs[countryStatSchema.UniqueAddresses]])
	if err != nil {
		return 0, nil, err
	}
	newUniqueAddr, err := influxdb.GetUint64FromInterface(v[idxs[countryStatSchema.NewUniqueAddresses]])
	if err != nil {
		return 0, nil, err
	}
	kyced, err := influxdb.GetUint64FromInterface(v[idxs[countryStatSchema.KYCedAddresses]])
	if err != nil {
		return 0, nil, err
	}
	totalBurnFee, err := influxdb.GetFloat64FromInterface(v[idxs[countryStatSchema.TotalBurnFee]])
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
		TotalBurnFee:       totalBurnFee,
	}, nil
}
