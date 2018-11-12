package cq

import "github.com/KyberNetwork/reserve-stats/lib/cq"

//CreateWalletStatsCqs return a new set of cqs required for wallet stats aggregation
func CreateWalletStatsCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	walletStatDayCqs, err := cq.NewContinuousQuery(
		"wallet_stats_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, walletStatDayCqs)
	return result, nil
}
