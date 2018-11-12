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
	userListEndpoint = "/user-list"
	userListFromTime = 1541548800000
	userListToTime   = 1541635200000
)

func TestUserListQuery(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test invalid request params",
			Endpoint: fmt.Sprintf("%s?from=%d&to=%d", userListEndpoint, userListToTime, userListFromTime),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg:      "Test valid request",
			Endpoint: fmt.Sprintf("%s?from=%d&to=%d", userListEndpoint, userListFromTime, userListToTime),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}
