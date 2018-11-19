package cq

import (
	"fmt"

	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	//IntegrationVolumeMeasurement is the name for Integration Volume Measurement
	IntegrationVolumeMeasurement = "integration_volume"
)

// CreateIntergrationVoluemCq return a set of cqs required for KyberSwap and non KyberSwap Summary aggregation
func CreateIntergrationVoluemCq(dbName string) ([]*libcq.ContinuousQuery, error) {
	var result []*libcq.ContinuousQuery

	kyberSwapVol, err := libcq.NewContinuousQuery(
		"kyber_swap_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		fmt.Sprintf("SELECT SUM(eth_amount) AS kyber_swap_volume INTO %s FROM trades WHERE integration_app='%s'", IntegrationVolumeMeasurement, common.KyberSwapAppName),
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
		fmt.Sprintf("SELECT SUM(eth_amount) AS non_kyber_swap_volume INTO %s FROM trades WHERE integration_app!='%s'", IntegrationVolumeMeasurement, common.KyberSwapAppName),
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, nonKyberSwapVol)
	return result, nil
}
