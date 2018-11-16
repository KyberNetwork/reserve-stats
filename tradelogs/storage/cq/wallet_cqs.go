package cq

import (
	"github.com/KyberNetwork/reserve-stats/lib/cq"
)

//CreateWalletStatsCqs return a new set of cqs required for wallet stats aggregation
func CreateWalletStatsCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	uniqueAddrCqs, err := cq.NewContinuousQuery(
		"unique_addr",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT COUNT(record) AS unique_addresses INTO wallet_stats FROM (SELECT SUM(eth_amount) AS record FROM trades GROUP BY user_addr) GROUP BY wallet_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, uniqueAddrCqs)
	volCqs, err := cq.NewContinuousQuery(
		"summary_volume",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume, COUNT(eth_amount) AS total_trade, "+
			"MEAN(usd_amount) AS usd_per_trade, MEAN(eth_amount) AS eth_per_trade INTO wallet_stats "+
			"FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades "+
			"WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') "+
			"OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY wallet_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, volCqs)
	return result, nil
}
