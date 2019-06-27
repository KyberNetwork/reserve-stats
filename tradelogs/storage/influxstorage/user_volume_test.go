package influxstorage

import (
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/cq"
	ethereum "github.com/ethereum/go-ethereum/common"
)

func getUserVolumeTestData(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateUserVolumeCqs(is.dbName)
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

func TestInfluxStorage_GetUserVolume(t *testing.T) {
	const (
		dbName = "test_burn_fee"
		// These params are expected to be change when export.dat changes.
		fromTime          = 1539000000000
		toTime            = 1539250666000
		expectedEthAmount = 0.05
		expectedUsdAmount = 11.283100808873083
		freq              = "h"
		timeStamp         = "2018-10-11T09:00:00Z"
		userAddressHex    = "0x0826601F28B691CEEa2Be05EC1c922Ea0eC2d82D"
	)

	is, err := newTestInfluxStorage(dbName)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, is.tearDown())
	}()
	require.NoError(t, loadTestData(dbName))
	require.NoError(t, getUserVolumeTestData(is))

	var (
		userAddress = ethereum.HexToAddress(userAddressHex)
		from        = timeutil.TimestampMsToTime(uint64(fromTime))
		to          = timeutil.TimestampMsToTime(uint64(toTime))
	)

	userVolume, err := is.GetUserVolume(userAddress, from, to, freq)
	require.NoError(t, err)
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	require.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	require.Contains(t, userVolume, timeUint)
	result := userVolume[timeUint]
	require.Equal(t, expectedEthAmount, result.ETHAmount)
	require.Equal(t, expectedUsdAmount, result.USDAmount)
}
