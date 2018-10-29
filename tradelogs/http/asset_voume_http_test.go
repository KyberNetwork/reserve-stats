package http

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	testETHAmount = 0.111
	testUSDAmount = 0.222
	testVolAmount = 0.333
)

type mockCore struct {
}

func (c *mockCore) Tokens() ([]core.Token, error) {
	return []core.Token{
		core.ETHToken,
	}, nil
}

func (c *mockCore) FromWei(ethereum.Address, *big.Int) (float64, error) {
	return 0, nil
}

func (c *mockCore) ToWei(ethereum.Address, float64) (*big.Int, error) {
	return nil, nil
}

func (s *mockStorage) GetAssetVolume(token core.Token, fromTime, toTime uint64, frequency string) (map[time.Time]*common.VolumeStats, error) {
	var (
		from           = timeutil.TimestampMsToTime(fromTime)
		to             = timeutil.TimestampMsToTime(fromTime)
		mockVolumeStat = common.VolumeStats{
			ETHAmount: testETHAmount,
			USDAmount: testUSDAmount,
			Volume:    testVolAmount,
		}
		mockResult = map[time.Time]*common.VolumeStats{
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
		validAsset               = "ETH"
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
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}

func expectInvalidInput(t *testing.T, resp *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func expectCorrectVolume(t *testing.T, resp *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, resp.Code)
	var result map[time.Time]common.VolumeStats
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Error("Could not decode result", "err", err)
	}
	correctVolume := common.VolumeStats{
		ETHAmount: testETHAmount,
		USDAmount: testUSDAmount,
		Volume:    testVolAmount,
	}
	for _, vol := range result {
		if !reflect.DeepEqual(vol, correctVolume) {
			t.Error("Wrong volume", "expected Volume", correctVolume, "returned volumes", vol)
		}
	}
}
