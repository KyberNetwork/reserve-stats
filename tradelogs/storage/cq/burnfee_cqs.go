package cq

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	burnschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/burnfee"
	burnVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/burnfee_volume"
)

const (
	// DayMeasurement is the measure to store aggregatedBurnFee in Day Frequency
	DayMeasurement = "burn_fee_day"
	// HourMeasurement is the measure to store aggregatedBurnFee in Hour Frequency
	HourMeasurement = "burn_fee_hour"
)

func prepareBurnfeeAggregationQuery(measurement string) string {
	q := fmt.Sprintf(
		`SELECT SUM(%s) as %s INTO %s FROM %s GROUP BY %s`,
		burnschema.Amount.String(),
		burnVolumeSchema.SumAmount.String(),
		measurement,
		common.BurnfeeMeasurementName,
		burnschema.ReserveAddr.String(),
	)
	return q
}

// CreateBurnFeeCqs return a set of cqs required for burnfee aggregation
func CreateBurnFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	burnfeeHourCqs, err := cq.NewContinuousQuery(
		HourMeasurement,
		dbName,
		hourResampleInterval,
		hourResampleFor,
		prepareBurnfeeAggregationQuery(HourMeasurement),
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, burnfeeHourCqs)
	burnfeeDayCqs, err := cq.NewContinuousQuery(
		DayMeasurement,
		dbName,
		dayResampleInterval,
		dayResampleFor,
		prepareBurnfeeAggregationQuery(DayMeasurement),
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, burnfeeDayCqs)
	return result, nil
}
