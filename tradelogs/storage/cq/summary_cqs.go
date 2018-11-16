package cq

import (
	libcq "github.com/KyberNetwork/reserve-stats/lib/cq"
)

// CreateSummaryCqs return a set of cqs required for trade Summary aggregation
func CreateSummaryCqs(dbName string) ([]*libcq.ContinuousQuery, error) {
	var result []*libcq.ContinuousQuery

	uniqueAddrCqs, err := libcq.NewContinuousQuery(
		"unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT COUNT(record) AS unique_addresses INTO trade_summary FROM (SELECT SUM(eth_amount) AS record FROM trades GROUP BY user_addr)",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, uniqueAddrCqs)

	volCqs, err := libcq.NewContinuousQuery(
		"summary_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(eth_amount) AS total_eth_volume, SUM(usd_amount) AS total_usd_amount, COUNT(eth_amount) AS total_trade, MEAN(usd_amount) AS usd_per_trade, MEAN(eth_amount) AS eth_per_trade INTO trade_summary FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE'))",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, volCqs)

	totalBurnFeeCqs, err := libcq.NewContinuousQuery(
		"summary_total_burn_fee",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(amount) AS total_burn_fee INTO burn_fee_summary FROM burn_fees",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, totalBurnFeeCqs)

	newUnqAddressCq, err := libcq.NewContinuousQuery(
		"new_unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT COUNT(traded) as new_unique_addresses INTO trade_summary FROM first_trades",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, newUnqAddressCq)

	kyced, err := libcq.NewContinuousQuery(
		"kyced",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT COUNT(kyced) as kyced INTO trade_summary FROM (SELECT DISTINCT(kyced) AS kyced FROM kyced GROUP BY user_addr)",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, kyced)

	return result, nil
}
