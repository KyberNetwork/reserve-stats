package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

//GetTokenHeatmap return list ordered country by asset volume
func (is *InfluxStorage) GetTokenHeatmap(asset core.Token, from, to time.Time, timezone int8) (map[string]common.Heatmap, error) {
	var (
		err             error
		result          = make(map[string]common.Heatmap)
		tokenAddr       = ethereum.HexToAddress(asset.Address).Hex()
		measurementName = "volume_country_stats"
	)

	measurementName = getMeasurementName(measurementName, timezone)

	logger := is.sugar.With(
		"func", "tradelogs/storage/InfluxStorage.GetTokenHeatmap",
		"asset", asset.ID,
		"from", from,
		"to", to,
	)

	// TODO: review and consider keep or remove total_trade, unique_addresses and kyc_user here
	// tradeQuery := fmt.Sprintf(`
	// SELECT SUM(total_trade) as total_trade, SUM(unique_addresses) as unique_addresses FROM country_stats
	// WHERE time >= '%s' AND time <= '%s' GROUP BY country
	// `, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))

	// logger.Debug(tradeQuery)

	// tradeResponse, err := is.queryDB(is.influxClient, tradeQuery)
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
	SELECT SUM(eth_volume) as eth_volume, SUM(token_volume) as token_volume, SUM(usd_volume) as usd_volume from %s
	WHERE (dst_addr='%s' or src_addr='%s') and (time >= '%s' AND time <= '%s') GROUP BY country
	`, measurementName, tokenAddr, tokenAddr, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))

	logger.Debug(volumeQuery)

	volumeResponse, err := is.queryDB(is.influxClient, volumeQuery)
	if err != nil {
		return result, err
	}

	for _, s := range volumeResponse[0].Series {
		if len(s.Values[0]) != 4 {
			logger.Debug(s.Values)
			return result, errors.New("values field is invalid in len")
		}
		country := s.Tags["country"]
		if country == "" {
			country = "unknown"
		}
		ethVolume, err := influxdb.GetFloat64FromInterface(s.Values[0][1])
		if err != nil {
			return result, err
		}
		tokenVolume, err := influxdb.GetFloat64FromInterface(s.Values[0][2])
		if err != nil {
			return result, err
		}
		usdVolume, err := influxdb.GetFloat64FromInterface(s.Values[0][3])
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
