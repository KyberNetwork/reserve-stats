package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

//GetWalletStats return stats of a wallet address from time to time by a frequency
func (is *InfluxStorage) GetWalletStats(from, to time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error) {
	var (
		logger = is.sugar.With(
			"func", "tradelogs/storage/InfluxStorage.GetWalletStats",
			"from", from,
			"to", to,
		)
		result          = make(map[uint64]common.WalletStats)
		err             error
		measurementName = "wallet_stats"
	)

	measurementName = getMeasurementName(measurementName, timezone)

	query := fmt.Sprintf(`
	SELECT eth_volume, usd_volume, total_trade, unique_addresses, usd_per_trade, eth_per_trade, kyced, new_unique_addresses, total_burn_fee
	FROM %s WHERE time >= '%s' and time <= '%s' and wallet_addr='%s'
	`, measurementName, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339), walletAddr)

	logger.Debug(query)

	response, err := influxdb.QueryDB(is.influxClient, query, is.dbName)
	if err != nil {
		return result, err
	}
	if len(response[0].Series) == 0 {
		return result, nil
	}

	for _, v := range response[0].Series[0].Values {
		ts, walletStats, err := convertQueryToWalletStats(v)
		if err != nil {
			return result, err
		}
		key := timeutil.TimeToTimestampMs(ts)
		result[key] = walletStats
	}
	return result, err
}

func convertQueryToWalletStats(v []interface{}) (time.Time, common.WalletStats, error) {
	var (
		ts          time.Time
		walletStats common.WalletStats
		err         error
	)
	if len(v) != 10 {
		return ts, walletStats, errors.New("value fields is invalid in len")
	}
	ts, err = influxdb.GetTimeFromInterface(v[0])
	if err != nil {
		return ts, walletStats, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[1])
	if err != nil {
		return ts, walletStats, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[2])
	if err != nil {
		return ts, walletStats, err
	}
	totalTrade, err := influxdb.GetInt64FromInterface(v[3])
	if err != nil {
		return ts, walletStats, err
	}
	uniqueAddresses, err := influxdb.GetInt64FromInterface(v[4])
	if err != nil {
		return ts, walletStats, err
	}
	usdPerTrade, err := influxdb.GetFloat64FromInterface(v[5])
	if err != nil {
		return ts, walletStats, err
	}
	ethPerTrade, err := influxdb.GetFloat64FromInterface(v[6])
	if err != nil {
		return ts, walletStats, err
	}

	kyced, err := influxdb.GetInt64FromInterface(v[7])
	if err != nil {
		return ts, walletStats, err
	}

	newUniqueAddress, err := influxdb.GetInt64FromInterface(v[8])
	if err != nil {
		return ts, walletStats, err
	}

	totalBurnFee, err := influxdb.GetFloat64FromInterface(v[9])
	if err != nil {
		return ts, walletStats, err
	}
	walletStats = common.WalletStats{
		ETHVolume:          ethVolume,
		USDVolume:          usdVolume,
		TotalTrade:         totalTrade,
		UniqueAddresses:    uniqueAddresses,
		USDPerTrade:        usdPerTrade,
		ETHPerTrade:        ethPerTrade,
		KYCEDAddresses:     kyced,
		NewUniqueAddresses: newUniqueAddress,
		BurnFee:            totalBurnFee,
	}
	return ts, walletStats, err
}
