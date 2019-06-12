package influxstorage

import (
	"log"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/cq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func aggregateIntegrationVolume(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateIntegrationVolumeCq(is.dbName)
	if err != nil {
		return err
	}
	log.Printf("%v", cqs)
	log.Println("ready to run")
	for _, cq := range cqs {
		err = cq.Execute(is.influxClient, is.sugar)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestIntegrationVolume(t *testing.T) {
	const (
		dbName = "test_integration_volume"
		// These params are expected to be change when export.dat changes.

		timeStamp = "2018-10-11T00:00:00Z"
	)

	var (
		fromTime = timeutil.TimestampMsToTime(1539216000000)
		toTime   = timeutil.TimestampMsToTime(1539254666000)
	)

	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, is.tearDown())
	}()

	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregateIntegrationVolume(is))
	integrationVol, err := is.GetIntegrationVolume(fromTime, toTime)

	require.NoError(t, err)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := integrationVol[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}
	assert.NotZero(t, result.KyberSwapVolume)
	assert.NotZero(t, result.NonKyberSwapVolume)
}
