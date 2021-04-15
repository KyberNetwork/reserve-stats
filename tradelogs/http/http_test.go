package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

func (s *mockStorage) GetIntegrationVolume(fromTime, toTime time.Time) (map[uint64]*common.IntegrationVolume, error) {
	return nil, nil
}

type mockStorage struct {
}

func (s *mockStorage) GetTokenSymbol(address string) (string, error) {
	return "", nil
}

func (s *mockStorage) UpdateTokens(addresses, symbols []string) error {
	return nil
}

func (s *mockStorage) LastBlock() (int64, error) {
	return 0, nil
}
func (s *mockStorage) SaveTradeLogs(logs *common.CrawlResult) error {
	return nil
}

func (s *mockStorage) LoadTradeLogs(from, to time.Time) ([]common.Tradelog, error) {
	return nil, nil
}

func (s *mockStorage) GetUserVolume(userAddr ethereum.Address, fromTime, toTime time.Time, freq string) (map[uint64]common.UserVolume, error) {
	return nil, nil
}

func (s *mockStorage) GetUserList(fromTime, toTime time.Time) ([]common.UserInfo, error) {
	return nil, nil
}

func (s *mockStorage) LoadTradeLogsByTxHash(txHash ethereum.Hash) ([]common.Tradelog, error) {
	return nil, nil
}

func (s *mockStorage) GetStats(from, to time.Time) (common.StatsResponse, error) {
	return common.StatsResponse{}, nil
}

func (s *mockStorage) GetTopTokens(from, to time.Time, limit uint64) (common.TopTokens, error) {
	return common.TopTokens{}, nil
}

func (s *mockStorage) GetTopIntegrations(from, to time.Time, limit uint64) (common.TopIntegrations, error) {
	return common.TopIntegrations{}, nil
}

func (s *mockStorage) GetTopReserves(from, to time.Time, limit uint64) (common.TopReserves, error) {
	return common.TopReserves{}, nil
}

func newTestServer() (*Server, error) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	return NewServer(
		&mockStorage{},
		"",
		sugar,
		nil,
	), nil

}

func TestTradeLogsRoute(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test valid request",
			Endpoint: "/trade-logs",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result []common.Tradelog
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

				assert.Contains(t, result.Error, "max time frame exceed")
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}
