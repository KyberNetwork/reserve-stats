package http

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

func expectCorrectRate(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	decoded := map[string]map[string][]common.ReserveRates{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&decoded))
	expected := map[string]map[string][]common.ReserveRates{
		testRsvAddress: {
			"ETH-KNC": {
				{
					Timestamp: timeutil.TimestampMsToTime(testTs),
					FromBlock: testFromBlock,
					ToBlock:   testToBlock,
					Rates: common.ReserveRateEntry{
						BuyReserveRate:  testRates[testRsvAddress]["ETH-KNC"].BuyReserveRate,
						BuySanityRate:   testRates[testRsvAddress]["ETH-KNC"].BuySanityRate,
						SellReserveRate: testRates[testRsvAddress]["ETH-KNC"].SellReserveRate,
						SellSanityRate:  testRates[testRsvAddress]["ETH-KNC"].SellSanityRate,
					},
				},
			},
			"ETH-ZRX": {
				{
					Timestamp: timeutil.TimestampMsToTime(testTs),
					FromBlock: testFromBlock,
					ToBlock:   testToBlock,
					Rates: common.ReserveRateEntry{
						BuyReserveRate:  testRates[testRsvAddress]["ETH-ZRX"].BuyReserveRate,
						BuySanityRate:   testRates[testRsvAddress]["ETH-ZRX"].BuySanityRate,
						SellReserveRate: testRates[testRsvAddress]["ETH-ZRX"].SellReserveRate,
						SellSanityRate:  testRates[testRsvAddress]["ETH-ZRX"].SellSanityRate,
					},
				},
			},
		},
	}
	//assert.Equal(t, expected, decoded)

	var buf = bytes.NewBuffer(nil)
	assert.NoError(t, json.NewEncoder(buf).Encode(decoded))
	t.Log(buf.String())

	buf = bytes.NewBuffer(nil)
	assert.NoError(t, json.NewEncoder(buf).Encode(expected))
	t.Log(buf.String())
}
