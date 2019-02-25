package influxstorage

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
<<<<<<< HEAD:tradelogs/storage/burnfee_test.go
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
=======
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/cq"
	ethereum "github.com/ethereum/go-ethereum/common"
>>>>>>> interface for tradelog storage:tradelogs/storage/influxstorage/burnfee_test.go
)

func aggregationBurnFeeTestData(is *InfluxStorage) error {
	cqs, err := tradelogcq.CreateBurnFeeCqs(is.dbName)
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

func TestGetBurnFee(t *testing.T) {
	const (
		dbName = "test_burnfee"
		// These params are expected to be change when export.dat changes.
		fromTime       = 1539000000000
		toTime         = 1539250666000
		expectedAmount = 1.3960673166743627
		freq           = "h"
		timeStamp      = "2018-10-11T09:00:00Z"
		reserveAddrHex = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	)

	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, is.tearDown())
	}()
	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregationBurnFeeTestData(is))

	var (
		rsvAddr  = ethereum.HexToAddress(reserveAddrHex)
		rsvAddrs = []ethereum.Address{
			rsvAddr,
		}
		from = timeutil.TimestampMsToTime(uint64(fromTime))
		to   = timeutil.TimestampMsToTime(uint64(toTime))
	)

	burnFee, err := is.GetAggregatedBurnFee(from, to, freq, rsvAddrs)
	assert.NoError(t, err)

	t.Logf("Burnfee result %v", burnFee)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := burnFee[rsvAddr]
	if !ok {
		t.Fatalf("expect to find result at rsv %s, yet there is none", rsvAddr.Hex())
	}
	amount, ok := result[strconv.FormatUint(timeUint, 10)]
	if !ok {
		t.Fatalf("expect to find result at rsv %s timestamp %s, yet there is none", rsvAddr.Hex(), timeStamp)
	}

	if amount != expectedAmount {
		t.Fatal(fmt.Errorf("Expect burnFee amount to be %.18f, got %.18f", expectedAmount, amount))
	}
}
