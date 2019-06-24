package postgrestorage

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func TestTradeLogDB_GetAggregatedBurnFee(t *testing.T) {
	const (
		dbName = "test_burn_fee"
		// These params are expected to be change when export.dat changes.
		fromTime       = 1539000000000
		toTime         = 1539250666000
		expectedAmount = 1.3960673166743627
		freq           = "h"
		timeStamp      = "2018-10-11T09:00:00Z"
		reserveAddrHex = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	)

	tldb, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, tldb.tearDown(dbName))
	}()
	err = loadTestData(tldb.db, testDataFile)
	require.NoError(t, err)

	var (
		rsvAddr  = ethereum.HexToAddress(reserveAddrHex)
		rsvAddrs = []ethereum.Address{
			rsvAddr,
		}
		from = timeutil.TimestampMsToTime(uint64(fromTime))
		to   = timeutil.TimestampMsToTime(uint64(toTime))
	)

	burnFee, err := tldb.GetAggregatedBurnFee(from, to, freq, rsvAddrs)
	require.NoError(t, err)
	require.Equal(t, 1, len(burnFee))

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	require.NoError(t, err)
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
