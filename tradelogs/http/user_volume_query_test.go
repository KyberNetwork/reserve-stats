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
	userVolumeEndpoint    = "/user-volume"
	userVolumeAddr        = "0x804aDa8c08A2E8ecff1a6535bf28DC4f1EfF4f8e"
	userVolumeInvalidAddr = "dah1oshfsaoh"
	userVolumeFromTime    = 1541548800000
	userVolumeToTime      = 1541635200000
)

func TestUserVolume(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test invalid time range",
			Endpoint: fmt.Sprintf("%s?from=%d&to=%d&userAddr=%s", userVolumeEndpoint, userVolumeToTime, userVolumeFromTime, userVolumeAddr),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg:      "Test invalid user address",
			Endpoint: fmt.Sprintf("%s?from=%d&to=%d&userAddr=%s", userVolumeEndpoint, userVolumeFromTime, userVolumeToTime, userVolumeInvalidAddr),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg:      "Test valid request",
			Endpoint: fmt.Sprintf("%s?from=%d&to=%d&userAddr=%s", userVolumeEndpoint, userVolumeFromTime, userVolumeToTime, userVolumeAddr),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}
