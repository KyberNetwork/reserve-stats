package cq

import (
	"fmt"

	appnames "github.com/KyberNetwork/reserve-stats/app-names"
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	integrationVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/integrationvolume"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
)

// CreateIntegrationVolumeCq return a set of cqs required for KyberSwap and non KyberSwap Summary aggregation
func CreateIntegrationVolumeCq(dbName string) ([]*libcq.ContinuousQuery, error) {
	var result []*libcq.ContinuousQuery

	kyberSwapVol, err := libcq.NewContinuousQuery(
		"kyber_swap_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s INTO %[3]s FROM %[4]s WHERE %[5]s='%[6]s'`,
			logSchema.EthAmount,
			integrationVolumeSchema.KyberSwapVolume.String(),
			common.IntegrationVolumeMeasurement,
			common.TradeLogMeasurementName,
			logSchema.IntegrationApp.String(),
			appnames.KyberSwapAppName),
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, kyberSwapVol)

	nonKyberSwapVol, err := libcq.NewContinuousQuery(
		"non_kyber_swap_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s INTO %[3]s FROM %[4]s WHERE %[5]s!='%[6]s'`,
			logSchema.EthAmount,
			integrationVolumeSchema.NonKyberSwapVolume.String(),
			common.IntegrationVolumeMeasurement,
			common.TradeLogMeasurementName,
			logSchema.IntegrationApp,
			appnames.KyberSwapAppName),
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, nonKyberSwapVol)
	return result, nil
}
