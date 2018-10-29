package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/stretchr/testify/assert"
)

const (
	endpoint        = "/wallet-fee"
	reserveAddr     = "0x63825c174ab367968EC60f061753D3bbD36A0D8F" // Kyber reserve
	walletAddr      = "0x8a654566edd646283c920e3225873fca5370f489" // a wallet address have rate at jun 1st
	fromTime        = 1527811200000                                // june 1st 0:00:00
	toTime          = 1527897599000                                // june 1st 23:59:59
	freq            = "D"                                          // Day freqency
	timezone        = 0                                            // UTC
	invalidAddress  = "dah1oshfsaoh"                               // address is not valid
	invalidFreq     = "q"                                          // not exist frequency
	invalidTimezone = 15                                           // not supported timezone
)

func TestWalletFeeQuery(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test invalid request, lack of params",
			Endpoint: endpoint,
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg: "Test address is not valid",
			Endpoint: fmt.Sprintf("%s?walletAddr=%s&reserve=%s&from=%d&to=%d&freq=%s&timezone=%d", endpoint,
				invalidAddress, reserveAddr, fromTime, toTime, freq, timezone),
			Method: http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg: "Test freq is not valid",
			Endpoint: fmt.Sprintf("%s?walletAddr=%s&reserve=%s&from=%d&to=%d&freq=%s&timezone=%d", endpoint,
				walletAddr, reserveAddr, fromTime, toTime, invalidFreq, timezone),
			Method: http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg: "Test timezone is not supported",
			Endpoint: fmt.Sprintf("%s?walletAddr=%s&reserve=%s&from=%d&to=%d&freq=%s&timezone=%d", endpoint,
				walletAddr, reserveAddr, fromTime, toTime, freq, invalidTimezone),
			Method: http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg: "Test valid query",
			Endpoint: fmt.Sprintf("%s?walletAddr=%s&reserve=%s&from=%d&to=%d&freq=%s&timezone=%d", endpoint,
				walletAddr, reserveAddr, fromTime, toTime, freq, timezone),
			Method: http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}
