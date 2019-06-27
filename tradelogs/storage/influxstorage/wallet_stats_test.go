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

func aggregateWalletStats(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateWalletStatsCqs(is.dbName)
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

func TestInfluxStorage_GetWalletStats(t *testing.T) {
	const (
		dbName = "test_wallet_stats_volume"
		// These params are expected to be change when export.dat changes.

		timeStamp     = "2018-10-11T00:00:00Z"
		walletAddress = "0xDECAF9CD2367cdbb726E904cD6397eDFcAe6068D"
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

	require.NoError(t, loadTestData(dbName))
	require.NoError(t, aggregateWalletStats(is))
	integrationVol, err := is.GetWalletStats(fromTime, toTime, walletAddress, 0)

	require.NoError(t, err)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := integrationVol[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	t.Logf("%+v", result)
}
