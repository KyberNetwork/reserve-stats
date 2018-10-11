package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/utils"
	"github.com/KyberNetwork/reserve-stats/reserve-rates-crawler/common"
)

// ReserveRateResponse is the struct to marshall
type ReserveRateResponse struct {
	Data    []common.ReserveRates
	Success bool
}

func expectCorrectRate(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()
	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
	decoded := ReserveRateResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatal(err)
	}
	rateQueried := decoded.Data[0]

	testRsvRate := common.ReserveRates{}
	if err := json.Unmarshal([]byte(testRsvRateJSON), &testRsvRate); err != nil {
		t.Error(err)
	}
	// Since DB's precision is in ms, compare the two timestamp in ms. s
	if (utils.TimeToTimestampMs(testRsvRate.Timestamp)) != (utils.TimeToTimestampMs(rateQueried.Timestamp)) {
		t.Fatalf("Result from http server is different with dabatase")
	}
	rateQueried.Timestamp = testRsvRate.Timestamp

	if !reflect.DeepEqual(rateQueried, testRsvRate) {
		t.Fatalf("Result from http server is different with dabatase")
	}
}
