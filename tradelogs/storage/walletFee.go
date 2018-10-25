package storage

import "time"

//GetAggregatedWalletFee return wallet fee follow by
//provided reserve address, wallet address, from time, to time
//frequency (minute, hour, day) and timezone
func (is *InfluxStorage) GetAggregatedWalletFee(reserveAddr, walletAddr, freq string,
	fromTime, toTime time.Time, timezone int64) (map[string]float64, error) {
	var (
		result map[string]float64
		err    error
	)
	return result, err
}
