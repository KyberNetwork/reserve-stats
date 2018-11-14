package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
	"github.com/stretchr/testify/assert"
)

func aggregateCountryStat(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateCountryCqs(is.dbName)
	if err != nil {
		return err
	}
	for _, cq := range cqs {
		err = cq.Execute(is.influxClient, is.sugar)
		if err != nil {
			return err
		}
	}
	return nil
}
func TestCountryStats(t *testing.T) {
	const (
		dbName = "test_country_stats"
		// These params are expected to be change when export.dat changes.
		fromTime  = 1539216000000
		toTime    = 1539254666000
		ethAmount = 0.0032361552369382478
		timeStamp = "2018-10-11T00:00:00Z"
		country   = "VN"
	)
	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, is.tearDown())
	}()

	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregateCountryStat(is))

	stats, err := is.GetCountryStats(country, fromTime, toTime)
	assert.NoError(t, err)
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := stats[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}
	if result.TotalETHVolume != ethAmount {
		t.Fatal(fmt.Errorf("Expect USD amount to be %.18f, got %.18f", ethAmount, result.TotalETHVolume))
	}
}
