package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	heatMapSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/heatmap"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//GetTokenHeatmap return list ordered country by asset volume
func (is *InfluxStorage) GetTokenHeatmap(asset ethereum.Address, from, to time.Time, timezone int8) (map[string]common.Heatmap, error) {
	var (
		err       error
		result    = make(map[string]common.Heatmap)
		tokenAddr = asset.Hex()
	)

	measurementName := getMeasurementName(common.HeatMapMeasurement, timezone)

	logger := is.sugar.With(
		"func", "tradelogs/storage/InfluxStorage.GetTokenHeatmap",
		"asset", asset.Hex(),
		"from", from,
		"to", to,
	)

	// TODO: review and consider keep or remove total_trade, unique_addresses and kyc_user here
	// tradeQuery := fmt.Sprintf(`
	// SELECT SUM(total_trade) as total_trade, SUM(unique_addresses) as unique_addresses FROM country_stats
	// WHERE time >= '%s' AND time <= '%s' GROUP BY country
	// `, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))

	// logger.Debug(tradeQuery)

	// tradeResponse, err := influxdb.QueryDB(is.influxClient, tradeQuery, is.dbName)
	// if err != nil {
	// 	return result, err
	// }

	// if len(tradeResponse[0].Series) == 0 {
	// 	return result, err
	// }

	// for _, s := range tradeResponse[0].Series {
	// 	country := s.Tags["country"]
	// 	if country == "" {
	// 		country = "unknown"
	// 	}
	// 	totalTrade, err := influxdb.GetInt64FromInterface(s.Values[0][1])
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	uniqueAddresses, err := influxdb.GetInt64FromInterface(s.Values[0][2])
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	result[country] = common.Heatmap{
	// 		TotalTrade:           totalTrade,
	// 		TotalUniqueAddresses: uniqueAddresses,
	// 	}
	// }

	volumeQuery := fmt.Sprintf(`
	SELECT SUM(%[1]s) as %[1]s, SUM(%[2]s) as %[2]s, SUM(%[3]s) as %[3]s from %[4]s
	WHERE (%[5]s='%[6]s' or %[7]s='%[6]s') and (time >= '%[8]s' AND time <= '%[9]s') GROUP BY %[10]s`,
		heatMapSchema.ETHVolume.String(),
		heatMapSchema.TokenVolume.String(),
		heatMapSchema.USDVolume.String(),
		measurementName,
		heatMapSchema.DstAddress.String(),
		tokenAddr,
		heatMapSchema.SrcAddress.String(),
		from.UTC().Format(time.RFC3339),
		to.UTC().Format(time.RFC3339),
		heatMapSchema.Country.String(),
	)

	logger.Debug(volumeQuery)

	volumeResponse, err := influxdb.QueryDB(is.influxClient, volumeQuery, is.dbName)
	if err != nil {
		return result, err
	}

	for _, s := range volumeResponse[0].Series {
		if len(s.Values[0]) != 4 {
			logger.Debug(s.Values)
			return result, errors.New("values field is invalid in len")
		}
		idxs, err := heatMapSchema.NewFieldsRegistrar(s.Columns)
		if err != nil {
			return result, err
		}
		country := s.Tags[heatMapSchema.Country.String()]
		if country == "" {
			country = "unknown"
		}
		ethVolume, err := influxdb.GetFloat64FromInterface(s.Values[0][idxs[heatMapSchema.ETHVolume]])
		if err != nil {
			return result, err
		}
		tokenVolume, err := influxdb.GetFloat64FromInterface(s.Values[0][idxs[heatMapSchema.TokenVolume]])
		if err != nil {
			return result, err
		}
		usdVolume, err := influxdb.GetFloat64FromInterface(s.Values[0][idxs[heatMapSchema.USDVolume]])
		if err != nil {
			return result, err
		}
		stat := result[country]
		stat.TotalETHValue = ethVolume
		stat.TotalTokenValue = tokenVolume
		stat.TotalFiatValue = usdVolume
		result[country] = stat
	}

	return result, err
}
