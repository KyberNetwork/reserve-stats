package cq

import (
	"bytes"
	"text/template"

	appnames "github.com/KyberNetwork/reserve-stats/app-names"
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	integrationVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/integrationvolume"
	logSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
)

// CreateIntegrationVolumeCq return a set of cqs required for KyberSwap and non KyberSwap Summary aggregation
func CreateIntegrationVolumeCq(dbName string) ([]*libcq.ContinuousQuery, error) {
	var result []*libcq.ContinuousQuery

	kyberSwapVolTemplate := `SELECT SUM({{.ETHAmount}}) AS {{.KyberSwapVolume}} INTO {{.IntegrationVolumeMeasurementName}} FROM {{.TradeLogMeasurementName}} WHERE {{.IntegrationApp}}='{{.KyberSwapAppName}}'`

	tmpl, err := template.New("kyberSwapVol").Parse(kyberSwapVolTemplate)
	if err != nil {
		return nil, err
	}

	var queryBuf bytes.Buffer
	if err = tmpl.Execute(&queryBuf, struct {
		ETHAmount                        string
		KyberSwapVolume                  string
		IntegrationVolumeMeasurementName string
		TradeLogMeasurementName          string
		IntegrationApp                   string
		KyberSwapAppName                 string
	}{
		ETHAmount:                        logSchema.EthAmount.String(),
		KyberSwapVolume:                  integrationVolumeSchema.KyberSwapVolume.String(),
		IntegrationVolumeMeasurementName: common.IntegrationVolumeMeasurement,
		TradeLogMeasurementName:          common.TradeLogMeasurementName,
		IntegrationApp:                   logSchema.IntegrationApp.String(),
		KyberSwapAppName:                 appnames.KyberSwapAppName,
	}); err != nil {
		return nil, err
	}

	kyberSwapVol, err := libcq.NewContinuousQuery(
		"kyber_swap_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryBuf.String(),
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, kyberSwapVol)

	nonKyberSwapVolTemplate := `SELECT SUM({{.ETHAmount}}) AS {{.NonKyberSwapVolume}} INTO {{.IntegrationVolumeMeasurementName}} FROM {{.TradeLogMeasurementName}} WHERE {{.IntegrationApp}}!='{{KyberSwapAppName}}'`

	tmpl, err = template.New("kyberSwapVol").Parse(nonKyberSwapVolTemplate)
	if err != nil {
		return nil, err
	}

	if err = tmpl.Execute(&queryBuf, struct {
		ETHAmount                        string
		NonKyberSwapVolume               string
		IntegrationVolumeMeasurementName string
		TradeLogMeasurementName          string
		IntegrationApp                   string
		KyberSwapAppName                 string
	}{
		ETHAmount:                        logSchema.EthAmount.String(),
		NonKyberSwapVolume:               integrationVolumeSchema.NonKyberSwapVolume.String(),
		IntegrationVolumeMeasurementName: common.IntegrationVolumeMeasurement,
		TradeLogMeasurementName:          common.TradeLogMeasurementName,
		IntegrationApp:                   logSchema.IntegrationApp.String(),
		KyberSwapAppName:                 appnames.KyberSwapAppName,
	}); err != nil {
		return nil, err
	}

	nonKyberSwapVol, err := libcq.NewContinuousQuery(
		"non_kyber_swap_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		queryBuf.String(),
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, nonKyberSwapVol)
	return result, nil
}
