package cq

import (
	"github.com/KyberNetwork/reserve-stats/lib/cq"
)

// CreateBurnFeeCqs return a set of cqs required for burnfee aggregation
func CreateBurnFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	burnfeeHourCqs, err := cq.NewContinuousQuery(
		"burn_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(amount) as sum_amount INTO burn_fee_hour FROM burn_fees GROUP BY reserve_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, burnfeeHourCqs)
	burnfeeDayCqs, err := cq.NewContinuousQuery(
		"burn_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(amount) as sum_amount INTO burn_fee_day FROM burn_fees GROUP BY reserve_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, burnfeeDayCqs)
	return result, nil
}
