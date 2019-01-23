package cq

import (
	"github.com/KyberNetwork/reserve-stats/lib/cq"
)

// CreateBurnFeeCqs return a set of cqs required for burnfee aggregation
func CreateBurnFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	srcBurnfeeHourCqs, err := cq.NewContinuousQuery(
		"src_burn_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(src_burn_fee) as sum_amount INTO burn_fee_hour FROM trades GROUP BY src_rsv_addr",
		"1h",
		[]string{},
	)

	if err != nil {
		return nil, err
	}
	result = append(result, srcBurnfeeHourCqs)

	dstBurnfeedstHourCqs, err := cq.NewContinuousQuery(
		"dst_burn_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(dst_burn_fee) as sum_amount INTO burn_fee_hour FROM trades GROUP BY dst_rsv_addr",
		"1h",
		[]string{},
	)

	if err != nil {
		return nil, err
	}
	result = append(result, dstBurnfeedstHourCqs)
	srcBurnfeeDayCqs, err := cq.NewContinuousQuery(
		"src_burn_fee_day",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(src_burn_fee) as sum_amount INTO burn_fee_day FROM trades GROUP BY src_rsv_addr",
		"1d",
		[]string{},
	)

	if err != nil {
		return nil, err
	}
	result = append(result, srcBurnfeeDayCqs)

	dstBurnfeedstDayCqs, err := cq.NewContinuousQuery(
		"dst_burn_fee_day",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(dst_burn_fee) as sum_amount INTO burn_fee_day FROM trades GROUP BY dst_rsv_addr",
		"1d",
		[]string{},
	)

	if err != nil {
		return nil, err
	}
	result = append(result, dstBurnfeedstDayCqs)
	return result, nil
}
