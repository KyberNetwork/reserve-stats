package storage

import (
	"testing"
	"time"
)

//BenchmarkWriteReadInflux can only be ran with pre-crawled database --from-block=6690000 --to-block=6700000
func BenchmarkWriteReadInflux(b *testing.B) {
	const (
		dbName   = "trade_logs"
		fromTime = 1518566400
		toTime   = 1542240000
	)
	is, err := newTestInfluxStorage(dbName)
	if err != nil {
		b.Fatal(err)
	}
	tradeLogs, err := is.LoadTradeLogs(time.Unix(fromTime, 0), time.Unix(toTime, 0))
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		if err = is.tearDown(); err != nil {
			b.Fatal(err)
		}
		if err = is.createDB(); err != nil {
			b.Fatal(err)
		}
		err = is.SaveTradeLogs(tradeLogs)
		if err != nil {
			b.Fatal(err)
		}

	}
}
