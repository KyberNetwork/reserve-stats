package cq

import "github.com/KyberNetwork/reserve-stats/lib/cq"

// CreateWalletFeeCqs return a set of cqs required for burnfee aggregation
func CreateWalletFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	walletFeeHourCqs, err := cq.NewContinuousQuery(
		"wallet_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(amount) as sum_amount INTO wallet_fee_hour FROM wallet_fees GROUP BY reserve_addr, wallet_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, walletFeeHourCqs)
	walletFeeDayCqs, err := cq.NewContinuousQuery(
		"wallet_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(amount) as sum_amount INTO wallet_fee_day FROM wallet_fees GROUP BY reserve_addr, wallet_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, walletFeeDayCqs)
	return result, nil
}
