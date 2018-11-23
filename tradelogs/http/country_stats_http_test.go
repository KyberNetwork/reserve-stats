package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/stretchr/testify/assert"
)

const (
	testCountryETHAmount = 0.111
	testCountryUSDAmount = 0.222
)

func (s *mockStorage) GetCountryStats(country string, fromTime, toTime time.Time, timezone int8) (map[uint64]*common.CountryStats, error) {
	from := timeutil.TimeToTimestampMs(fromTime)
	to := timeutil.TimeToTimestampMs(toTime)
	var (
		mockCountryStat = common.CountryStats{
			TotalETHVolume: testCountryETHAmount,
			TotalUSDVolume: testCountryUSDAmount,
		}
		mockResult = map[uint64]*common.CountryStats{
			from: &mockCountryStat,
			to:   &mockCountryStat,
		}
	)

	return mockResult, nil
}

func TestCountryStatsHttp(t *testing.T) {
	var (
		endpoint       = "/country-stats"
		country        = "VN"
		invalidCountry = "xqdc"
		unknownCountry = "UNKNOWN"
		validFrom      = 1539129600000
		invalidFrom    = "xxxx"
		validTo        = 1539216000000
		timezone       = 0
		// mock core only return ETH, KNC is not in the list of mock core's clients
		invalidFromInputEndpoint = fmt.Sprintf("%s?from=%s&to=%d&country=%s", endpoint, invalidFrom, validTo, country)
		invalidCountryEndpoint   = fmt.Sprintf("%s?from=%d&to=%d&country=%s", endpoint, validFrom, validTo, invalidCountry)
		validEndpoint            = fmt.Sprintf("%s?from=%d&to=%d&country=%s&timezone=%d", endpoint, validFrom, validTo, country, timezone)
		unknownCountryEndpoint   = fmt.Sprintf("%s?from=%d&to=%d&country=%s", endpoint, validFrom, validTo, unknownCountry)
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
			Msg:      "Test invalid country Input",
			Endpoint: invalidCountryEndpoint,
			Method:   http.MethodGet,
			Assert:   expectInvalidInput,
		},
		{
			Msg:      "Test valid Input",
			Endpoint: validEndpoint,
			Method:   http.MethodGet,
			Assert:   expectCorrectCountryStats,
		},
		{
			Msg:      "Test unkown country input",
			Endpoint: unknownCountryEndpoint,
			Method:   http.MethodGet,
			Assert:   expectCorrectCountryStats,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}

func expectCorrectCountryStats(t *testing.T, resp *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, resp.Code)
	var result map[uint64]*common.CountryStats
	log.Printf("response body: %s", resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Error("Could not decode result", "err", err)
	}
	correctStat := common.CountryStats{
		TotalETHVolume: testCountryETHAmount,
		TotalUSDVolume: testCountryUSDAmount,
	}
	t.Logf("result %v", result)
	for _, stat := range result {
		if !reflect.DeepEqual(*stat, correctStat) {
			t.Error("Wrong stat", "expected stat", correctStat, "returned stat", *stat)
		}
	}
}
