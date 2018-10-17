package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

// ReserveRateResponse is the struct to marshall
type ReserveRateResponse map[string]map[uint64]common.ReserveRates

func expectCorrectRate(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	testRsvRate := common.ReserveRates{}
	if err := json.Unmarshal([]byte(testRsvRateJSON), &testRsvRate); err != nil {
		t.Error(err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
	decoded := ReserveRateResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatal(err)
	}
	rates, ok := decoded[testRsvAddress]
	if !ok {
		t.Fatalf("response data doesn't contain expected reserve address: %s", testRsvAddress)
	}

	rate, ok := rates[testRsvRate.BlockNumber]
	if !ok {
		t.Fatalf("response data doesn't contain expected block number: %d", testRsvRate.BlockNumber)
	}

	// Since DB's precision is in ms, compare the two timestamp in ms. s
	if (timeutil.TimeToTimestampMs(testRsvRate.Timestamp)) != (timeutil.TimeToTimestampMs(rate.Timestamp)) {
		t.Fatalf("wrong timestamp, expected: %s, got: %s", testRsvRate.Timestamp, rate.Timestamp)
	}
}
