package influx

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
	"time"

	influxModel "github.com/influxdata/influxdb/models"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	countryStatSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/country_stats"
)

//GetCountryStats return stats of a country from time to time by daily fred in provided timezone
func (is *Storage) GetCountryStats(countryCode string, from, to time.Time, timezone int8) (map[uint64]*common.CountryStats, error) {
	var (
		err    error
		logger = is.sugar.With("country", countryCode,
			"fromTime", from, "toTime", to)
		timeFilter      = fmt.Sprintf("(time >='%s' AND time <= '%s')", from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
		countryFilter   = fmt.Sprintf("(country='%s')", countryCode)
		measurementName = common.CountryStatsMeasurementName
	)
	measurementName = getMeasurementName(measurementName, timezone)

	queryTmpl := `SELECT {{.ETHPerTrade}}, {{.TotalETHVolume}}, {{.TotalTrade}}, {{.TotalUSDAmount}}, {{.USDPerTrade}},
		 {{.UniqueAddresses}}, {{.NewUniqueAddresses}}, {{.KYCedAddresses}}, {{.TotalBurnFee}} FROM 
		 {{.MeasurementName}} WHERE {{.TimeFilter}} AND {{.CountryFilter}}`

	tmpl, err := template.New("countryStatsTemplate").Parse(queryTmpl)
	if err != nil {
		return nil, err
	}
	var queryStmtBuf bytes.Buffer
	if err = tmpl.Execute(&queryStmtBuf, struct {
		ETHPerTrade        string
		TotalETHVolume     string
		TotalTrade         string
		TotalUSDAmount     string
		USDPerTrade        string
		UniqueAddresses    string
		NewUniqueAddresses string
		KYCedAddresses     string
		TotalBurnFee       string
		MeasurementName    string
		TimeFilter         string
		CountryFilter      string
	}{
		ETHPerTrade:        countryStatSchema.ETHPerTrade.String(),
		TotalETHVolume:     countryStatSchema.TotalETHVolume.String(),
		TotalTrade:         countryStatSchema.TotalTrade.String(),
		TotalUSDAmount:     countryStatSchema.TotalUSDAmount.String(),
		USDPerTrade:        countryStatSchema.USDPerTrade.String(),
		UniqueAddresses:    countryStatSchema.UniqueAddresses.String(),
		NewUniqueAddresses: countryStatSchema.NewUniqueAddresses.String(),
		KYCedAddresses:     countryStatSchema.KYCedAddresses.String(),
		TotalBurnFee:       countryStatSchema.TotalBurnFee.String(),
		MeasurementName:    measurementName,
		TimeFilter:         timeFilter,
		CountryFilter:      countryFilter,
	}); err != nil {
		return nil, err
	}

	logger.Debugw("get country stats", "query", queryStmtBuf.String())

	response, err := influxdb.QueryDB(is.influxClient, queryStmtBuf.String(), is.dbName)
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
		return nil, err
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
