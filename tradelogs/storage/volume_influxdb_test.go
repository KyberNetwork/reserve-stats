package storage

import (
	"fmt"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
)

func aggregationTestData(is *InfluxStorage) error {

	cqs, err := tradelogcq.CreateAssetVolumeCqs(is.dbName)
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

func aggregationVolumeTestData(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateReserveVolumeCqs(is.dbName)
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

func TestGetAssetVolume(t *testing.T) {
	const (
		dbName = "test_volume"
		// These params are expected to be change when export.dat changes.
		fromTime   = 1539248043000
		toTime     = 1539248666000
		ethAmount  = 238.33849929550047
		freq       = "h"
		timeStamp  = "2018-10-11T09:00:00Z"
		ethAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	)

	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)

	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)

	defer func() {
		assert.NoError(t, is.tearDown())
	}()
	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregationTestData(is))
	volume, err := is.GetAssetVolume(ethereum.HexToAddress(ethAddress), from, to, freq)
	assert.NoError(t, err)

	t.Logf("Volume result %v", volume)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := volume[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	if result.USDAmount != ethAmount {
		t.Fatal(fmt.Errorf("Expect USD amount to be %.18f, got %.18f", ethAmount, result.USDAmount))
	}
}

func TestGetReserveVolume(t *testing.T) {
	const (
		dbName = "test_rsv_volume"

		// These params are expected to be change when export.dat changes.
		fromTime   = 1539248043000
		toTime     = 1539248666000
		ethAmount  = 227.05539848662738
		freq       = "h"
		timeStamp  = "2018-10-11T09:00:00Z"
		rsvAddrStr = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
		ethAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	)

	is, err := newTestInfluxStorage(dbName)
	defer func() {
		if err := is.tearDown(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := loadTestData(dbName); err != nil {
		t.Fatal(err)
	}
	if err := aggregationVolumeTestData(is); err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	from := timeutil.TimestampMsToTime(fromTime)
	to := timeutil.TimestampMsToTime(toTime)

	volume, err := is.GetReserveVolume(ethereum.HexToAddress(rsvAddrStr), ethereum.HexToAddress(ethAddress), from, to, freq)
	t.Logf("Volume result %v", volume)
	if err != nil {
		t.Fatal(err)
	}
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	if err != nil {
		t.Fatal(err)
	}
	result, ok := volume[timeutil.TimeToTimestampMs(timeUnix)]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	if result.USDAmount != ethAmount {
		t.Fatal(fmt.Errorf("Expect USD amount to be %.18f, got %.18f", ethAmount, result.USDAmount))
	}
}
