package influx

import (
	"errors"
	"fmt"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	walletStatSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/walletstats"
)

//GetWalletStats return stats of a wallet address from time to time by a frequency
func (is *Storage) GetWalletStats(from, to time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error) {
	var (
		logger = is.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", from,
			"to", to,
		)
		result          = make(map[uint64]common.WalletStats)
		err             error
		measurementName = common.WalletStatsMeasurement
	)

	measurementName = getMeasurementName(measurementName, timezone)

	query := fmt.Sprintf(
		`SELECT %[1]s, %[2]s, %[3]s, %[4]s, %[5]s, %[6]s, %[7]s, %[8]s, %[9]s `+
			`FROM %[10]s WHERE %[11]s >= '%[12]s' and %[11]s <= '%[13]s' and %[14]s='%[15]s'`,
		walletStatSchema.ETHVolume.String(),
		walletStatSchema.USDVolume.String(),
		walletStatSchema.TotalTrade.String(),
		walletStatSchema.UniqueAddresses.String(),
		walletStatSchema.USDPerTrade.String(),
		walletStatSchema.ETHPerTrade.String(),
		walletStatSchema.KYCedAddresses.String(),
		walletStatSchema.NewUniqueAddresses.String(),
		walletStatSchema.TotalBurnFee.String(),
		measurementName,
		walletStatSchema.Time.String(),
		from.UTC().Format(time.RFC3339),
		to.UTC().Format(time.RFC3339),
		walletStatSchema.WalletAddress.String(),
		walletAddr,
	)

	logger.Debug(query)

	response, err := influxdb.QueryDB(is.influxClient, query, is.dbName)
	if err != nil {
		return result, err
	}
	if len(response[0].Series) == 0 {
		return result, nil
	}

	for _, v := range response[0].Series[0].Values {

		idxs, err := walletStatSchema.NewFieldsRegistrar(response[0].Series[0].Columns)
		if err != nil {
			return result, err
		}
		ts, walletStats, err := convertQueryToWalletStats(v, idxs)
		if err != nil {
			return result, err
		}

		key := timeutil.TimeToTimestampMs(ts)
		result[key] = walletStats
	}
	return result, err
}

func convertQueryToWalletStats(v []interface{}, idxs walletStatSchema.FieldsRegistrar) (time.Time, common.WalletStats, error) {
	var (
		ts          time.Time
		walletStats common.WalletStats
		err         error
	)
	if len(v) != 10 {
		return ts, walletStats, errors.New("value fields is invalid in len")
	}
	ts, err = influxdb.GetTimeFromInterface(v[idxs[walletStatSchema.Time]])
	if err != nil {
		return ts, walletStats, err
	}
	ethVolume, err := influxdb.GetFloat64FromInterface(v[idxs[walletStatSchema.ETHVolume]])
	if err != nil {
		return ts, walletStats, err
	}
	usdVolume, err := influxdb.GetFloat64FromInterface(v[idxs[walletStatSchema.USDVolume]])
	if err != nil {
		return ts, walletStats, err
	}
	totalTrade, err := influxdb.GetInt64FromInterface(v[idxs[walletStatSchema.TotalTrade]])
	if err != nil {
		return ts, walletStats, err
	}
	uniqueAddresses, err := influxdb.GetInt64FromInterface(v[idxs[walletStatSchema.UniqueAddresses]])
	if err != nil {
		return ts, walletStats, err
	}
	usdPerTrade, err := influxdb.GetFloat64FromInterface(v[idxs[walletStatSchema.USDPerTrade]])
	if err != nil {
		return ts, walletStats, err
	}
	ethPerTrade, err := influxdb.GetFloat64FromInterface(v[idxs[walletStatSchema.ETHPerTrade]])
	if err != nil {
		return ts, walletStats, err
	}

	kyced, err := influxdb.GetInt64FromInterface(v[idxs[walletStatSchema.KYCedAddresses]])
	if err != nil {
		return ts, walletStats, err
	}

	newUniqueAddress, err := influxdb.GetInt64FromInterface(v[idxs[walletStatSchema.NewUniqueAddresses]])
	if err != nil {
		return ts, walletStats, err
	}

	totalBurnFee, err := influxdb.GetFloat64FromInterface(v[idxs[walletStatSchema.TotalBurnFee]])
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
