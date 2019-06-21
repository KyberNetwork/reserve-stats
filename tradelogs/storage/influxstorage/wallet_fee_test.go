package influxstorage

import (
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"

	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/cq"
)

func aggregateWalletFee(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateWalletFeeCqs(is.dbName)
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

func TestInfluxStorage_GetAggregatedWalletFee(t *testing.T) {
	const (
		dbName = "test_aggregated_wallet_fee"
		// These params are expected to be change when export.dat changes.
		timeStamp               = "2018-10-11T00:00:00Z"
		walletFeeExpectedAmount = float64(6.66)
	)

	var (
		fromTime    = timeutil.TimestampMsToTime(1539216000000)
		toTime      = timeutil.TimestampMsToTime(1539254666000)
		reserveAddr = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
		walletAddr  = "0xDECAF9CD2367cdbb726E904cD6397eDFcAe6068D"
	)

	is, err := newTestInfluxStorage(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, is.tearDown())
	}()

	require.NoError(t, loadTestData(dbName))
	require.NoError(t, aggregateWalletFee(is))
	integrationVol, err := is.GetAggregatedWalletFee(reserveAddr, walletAddr, "d", fromTime, toTime, 0)

	require.NoError(t, err)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)

	require.Contains(t, integrationVol, timeUint)
	result := integrationVol[timeUint]
	assert.Equal(t, walletFeeExpectedAmount, result)
}
