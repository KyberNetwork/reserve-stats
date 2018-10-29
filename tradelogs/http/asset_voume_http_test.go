package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

const (
	testETHAmount = 0.111
	testUSDAmount = 0.222
	testVolAmount = 0.333
)

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
		endpoint             = "/asset-volume"
		invalidAddress       = "0ABC"
		validFrom            = 1539129600000
		validTo              = 1539302400000
		validAsset           = "ETH"
		invalidInputEndpoint = fmt.Sprintf("%s?from=%d&to=%d&asset=%s&reserve=%s", endpoint, validFrom, validTo, validAsset, invalidAddress)
	)
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test invalid request",
			Endpoint: invalidInputEndpoint,
			Method:   http.MethodGet,
			Assert:   expectInvalidInput,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}

func expectInvalidInput(t *testing.T, resp *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
