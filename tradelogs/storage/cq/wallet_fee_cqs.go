package cq

import (
	"github.com/KyberNetwork/reserve-stats/lib/cq"
)

// CreateWalletFeeCqs return a set of cqs required for burnfee aggregation
func CreateWalletFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	srcWalletFeeHourCqs, err := cq.NewContinuousQuery(
		"src_wallet_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(src_wallet_fee_amount) as sum_amount INTO wallet_fee_hour FROM trades WHERE src_rsv_addr!='' GROUP BY src_rsv_addr, wallet_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, srcWalletFeeHourCqs)

	dstWalletFeeHourCqs, err := cq.NewContinuousQuery(
		"dst_wallet_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(dst_wallet_fee_amount) as sum_amount INTO wallet_fee_hour FROM trades WHERE dst_rsv_addr!='' GROUP BY dst_rsv_addr, wallet_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, dstWalletFeeHourCqs)
	srcWalletFeeDayCqs, err := cq.NewContinuousQuery(
		"src_wallet_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(src_wallet_fee_amount) as sum_amount INTO wallet_fee_day FROM trades WHERE src_rsv_addr!='' GROUP BY src_rsv_addr, wallet_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, srcWalletFeeDayCqs)
	dstWalletFeeDayCqs, err := cq.NewContinuousQuery(
		"dst_wallet_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(dst_wallet_fee_amount) as sum_amount INTO wallet_fee_day FROM trades WHERE dst_rsv_addr!='' GROUP BY dst_rsv_addr, wallet_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, dstWalletFeeDayCqs)

	return result, nil
}
