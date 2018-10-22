package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type mockStorage struct {
}

func (s *mockStorage) SaveTradeLogs(logs []common.TradeLog, rates []tokenrate.ETHUSDRate) error {
	return nil
}

func (s *mockStorage) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	return nil, nil
}

func newTestServer() *Server {
	return &Server{storage: &mockStorage{}}
}

func TestTradeLogsRoute(t *testing.T) {
	s := newTestServer()
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test valid request",
			Endpoint: "/trade-logs",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result []common.TradeLog
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
			},
		},
		{
			Msg:      "Test invalid time range",
			Endpoint: fmt.Sprintf("/trade-logs?from=0&to=%d", time.Hour/time.Millisecond*25),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)

				var result struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}

				assert.Contains(t, result.Error, "time range is too broad")
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}
