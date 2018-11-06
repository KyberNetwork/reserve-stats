package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
)

func aggregateTradeSummary(is *InfluxStorage) error {
	const (
		endpoint = "http://127.0.0.1:8086/"
	)
	cqs, err := tradelogcq.CreateSummaryCqs(is.dbName)
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

func TestTradeSummary(t *testing.T) {
	const (
		dbName = "test_trade_summary"
		// These params are expected to be change when export.dat changes.
		fromTime  = 1539216000000
		toTime    = 1539254666000
		ethAmount = 17.390905490542348
		timeStamp = "2018-10-11T00:00:00Z"
	)

	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, is.tearDown())
	}()
	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregateTradeSummary(is))
	summary, err := is.GetTradeSummary(fromTime, toTime)
	assert.NoError(t, err)

	t.Logf("Volume result %v", summary)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := summary[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	if result.ETHVolume != ethAmount {
		t.Fatal(fmt.Errorf("Expect USD amount to be %.18f, got %.18f", ethAmount, result.ETHVolume))
	}
}
