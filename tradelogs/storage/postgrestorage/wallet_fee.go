package postgrestorage

import "time"

//TODO implement this
func (tldb *TradeLogDB) GetAggregatedWalletFee(reserveAddr, walletAddr, freq string,
	fromTime, toTime time.Time, timezone int8) (map[uint64]float64, error) {
	return nil, nil
}
