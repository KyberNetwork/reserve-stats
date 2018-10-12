package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/timeutils"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
)

// ReserveRateResponse is the struct to marshall
type ReserveRateResponse struct {
	Data    map[uint64]common.ReserveRates
	Success bool
}

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
	rateQueried, ok := decoded.Data[testRsvRate.BlockNumber]
	if !ok {
		t.Fatalf("Rate queried doesn't containt the desire block number")
	}

	// Since DB's precision is in ms, compare the two timestamp in ms. s
	if (timeutils.TimeToTimestampMs(testRsvRate.Timestamp)) != (timeutils.TimeToTimestampMs(rateQueried.Timestamp)) {
		t.Fatalf("Result from http server is different with dabatase")
	}
	rateQueried.Timestamp = testRsvRate.Timestamp

	if !reflect.DeepEqual(rateQueried, testRsvRate) {
		t.Fatalf("Result from http server is different with dabatase")
	}
}
