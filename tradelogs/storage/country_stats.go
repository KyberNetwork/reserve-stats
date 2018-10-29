package storage

import "time"

//GetCountryStats return stats of a country from time to time by daily fred in provided timezone
func (is *InfluxStorage) GetCountryStats(countryCode string, timezone int64, fromTime, toTime time.Time) (map[uint64]float64, error) {
	var (
		result map[uint64]float64
		err    error
	)
	return result, err
}
