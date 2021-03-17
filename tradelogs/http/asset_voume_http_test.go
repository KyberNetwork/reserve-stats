package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	testUSDAmount = 0.222
	testVolAmount = 0.333
)

func (s *mockStorage) GetReserveVolume(rsvAddr ethereum.Address, token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}

func (s *mockStorage) GetMonthlyVolume(reserveAddr ethereum.Address, fromTime, toTime time.Time) (map[uint64]*common.VolumeStats, error) {
	return nil, nil
}

func (s *mockStorage) GetAssetVolume(token ethereum.Address, fromTime, toTime time.Time, frequency string) (map[uint64]*common.VolumeStats, error) {
	from := timeutil.TimeToTimestampMs(fromTime)
	to := timeutil.TimeToTimestampMs(toTime)
	var (
		mockVolumeStat = common.VolumeStats{
			USDAmount: testUSDAmount,
			Volume:    testVolAmount,
		}
		mockResult = map[uint64]*common.VolumeStats{
			from: &mockVolumeStat,
			to:   &mockVolumeStat,
		}
	)

	return mockResult, nil
}

func TestAssetVolumeHttp(t *testing.T) {
	var (
		endpoint    = "/asset-volume"
		freq        = ""
		validFrom   = 1539129600000
		invalidFrom = "xxxx"
		validTo     = 1539302400000
		// mock core only return ETH, KNC is not in the list of mock core's clients
		validAsset               = blockchain.USDTAddr.Hex()
		invalidAsset             = "KNC"
		invalidFromInputEndpoint = fmt.Sprintf("%s?from=%s&to=%d&asset=%s&freq=%s", endpoint, invalidFrom, validTo, validAsset, freq)
		invalidAssetEndpoint     = fmt.Sprintf("%s?from=%s&to=%d&asset=%s&freq=%s", endpoint, invalidFrom, validTo, invalidAsset, freq)
		validEndpoint            = fmt.Sprintf("%s?from=%d&to=%d&asset=%s", endpoint, validFrom, validTo, validAsset)
	)
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test invalid from Input",
			Endpoint: invalidFromInputEndpoint,
			Method:   http.MethodGet,
			Assert:   expectInvalidInput,
		},
		{
			Msg:      "Test invalid asset Input",
			Endpoint: invalidAssetEndpoint,
			Method:   http.MethodGet,
			Assert:   expectInvalidInput,
		},
		{
			Msg:      "Test valid Input",
			Endpoint: validEndpoint,
			Method:   http.MethodGet,
			Assert:   expectCorrectVolume,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}

func expectInvalidInput(t *testing.T, resp *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func expectCorrectVolume(t *testing.T, resp *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, resp.Code)
	var result map[uint64]common.VolumeStats
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Error("Could not decode result", "err", err)
	}
	correctVolume := common.VolumeStats{
		USDAmount: testUSDAmount,
		Volume:    testVolAmount,
	}
	for _, vol := range result {
		if !reflect.DeepEqual(vol, correctVolume) {
			t.Error("Wrong volume", "expected Volume", correctVolume, "returned volumes", vol)
		}
	}
}
