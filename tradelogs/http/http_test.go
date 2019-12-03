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
func (s *mockStorage) SaveTradeLogs(logs []common.TradeLog) error {
	return nil
}

func (s *mockStorage) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	return nil, nil
}

func (s *mockStorage) GetAggregatedBurnFee(from, to time.Time, freq string, reserveAddrs []ethereum.Address) (map[ethereum.Address]map[string]float64, error) {
	return nil, nil
}

func (s *mockStorage) GetAggregatedWalletFee(reserveAddr, walletAddr, freq string, fromTime, toTime time.Time, timezone int8) (map[uint64]float64, error) {
	return nil, nil
}

func (s *mockStorage) GetTradeSummary(from, to time.Time, timezone int8) (map[uint64]*common.TradeSummary, error) {
	return nil, nil
}

func (s *mockStorage) GetUserVolume(userAddr ethereum.Address, fromTime, toTime time.Time, freq string) (map[uint64]common.UserVolume, error) {
	return nil, nil
}

func (s *mockStorage) GetUserList(fromTime, toTime time.Time) ([]common.UserInfo, error) {
	return nil, nil
}

func (s *mockStorage) GetWalletStats(fromTime, toTime time.Time, walletAddr string, timezone int8) (map[uint64]common.WalletStats, error) {
	return nil, nil
}

func (s *mockStorage) GetTokenHeatmap(token ethereum.Address, from, to time.Time, timezone int8) (map[string]common.Heatmap, error) {
	return nil, nil
}

func (s *mockStorage) LoadTradeLogsByTxHash(txHash ethereum.Hash) ([]common.TradeLog, error) {
	return nil, nil
}

func (s *mockStorage) GetStats(from, to time.Time) (common.StatsResponse, error) {
	return common.StatsResponse{}, nil
}

func (s *mockStorage) GetTopTokens(from, to time.Time) (common.TopTokens, error) {
	return common.TopTokens{}, nil
}

func (s *mockStorage) GetTopIntegrations(from, to time.Time) (common.TopIntegrations, error) {
	return common.TopIntegrations{}, nil
}

func (s *mockStorage) GetTopReserves(from, to time.Time) (common.TopReserves, error) {
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

				assert.Contains(t, result.Error, "max time frame exceed")
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, router) })
	}
}

func TestBurnFeeRoute(t *testing.T) {
	s, err := newTestServer()
	if err != nil {
		t.Fatal(err)
	}
	router := s.setupRouter()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "Test valid request",
			Endpoint: "/burn-fee?freq=h&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result map[ethereum.Address]map[string]float64
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
			},
		},
		{
			Msg:      "Test missing reserve address",
			Endpoint: "/burn-fee?freq=h",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)

				var result map[ethereum.Address]map[string]float64
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}
			},
		},
		{
			Msg:      "Test invalid reserve address",
			Endpoint: "/burn-fee?freq=h&reserve=invalidAddress",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)

				var result struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}

				assert.Contains(t, result.Error, "Field validation for 'ReserveAddrs[0]' failed on the 'isAddress' tag")
			},
		},
		{
			Msg:      "Test invalid frequency",
			Endpoint: "/burn-fee?freq=invalid&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F",
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)

				var result struct {
					Error string `json:"error"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Error("Could not decode result", "err", err)
				}

				assert.Contains(t, result.Error, "invalid frequency")
			},
		},
		{
			Msg:      "Test time range too broad",
			Endpoint: fmt.Sprintf("/burn-fee?from=0&to=%d&freq=h&reserve=0x63825c174ab367968EC60f061753D3bbD36A0D8F", hourlyBurnFeeMaxDuration/time.Millisecond+1),
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
