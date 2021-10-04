package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/withdrawalstorage"
	huobiPostgres "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history/postgres"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

var (
	ts       *Server
	hdb      *huobiPostgres.HuobiStorage
	bdb      *withdrawalstorage.BinanceStorage
	tdb      *sqlx.DB
	teardown func() error
)

func TestGetHuobiWithdrawal(t *testing.T) {
	t.Log("creating a test get huobi withdrawal test")
	var (
		testTimestamp uint64 = 1525754125590
		huobiTestData        = []huobi.WithdrawHistory{
			{
				ID:         2272335,
				CreatedAt:  testTimestamp,
				UpdatedAt:  1525754753403,
				Currency:   "ETH",
				Type:       "withdraw",
				Amount:     0.48957444,
				State:      "confirmed",
				Fee:        0.01,
				Address:    "f6a605cdd9b2471ffdff706f8b7665a12b862158",
				AddressTag: "",
				TxHash:     "cdef3adad017d9564e62282f5e0f0d87d72b995759f1f7f4e473137cc1b96e56",
			},
		}

		binanceTestData = []binance.WithdrawHistory{
			{
				ID:        "7213fea8e94b4a5593d507237e5a555b",
				Amount:    "1",
				Address:   "0x6915f16f8791d0a1cc2bf47c13a6b2a92000504b",
				Asset:     "ETH",
				TxID:      "0xdf33b22bdb2b28b1f75ccd201a4a4m6e7g83jy5fc5d5a9d1340961598cfcb0a1",
				ApplyTime: "2018-05-08 4:35:26",
				TxFee:     "12131",
				Status:    4,
			},
			{
				ID:        "7213fea8e94b4a5534ggsd237e5a555b",
				Amount:    "1000",
				Address:   "463tWEBn5XZJSxLU34r6g7h8jtxuNcDbjLSjkn3XAXHCbLrTTErJrBWYgHJQyrCwkNgYvyV3z8zctJLPCZy24jvb3NiTcTJ",
				Asset:     "XMR",
				TxID:      "b3c6219639c8ae3f9cf010cdc24fw7f7yt8j1e063f9b4bd1a05cb44c4b6e2509",
				ApplyTime: "2018-05-08 4:35:27",
				Status:    4,
				TxFee:     "31231",
			},
		}

		binanceExpectedResponse = []BinanceWithdrawalResponse{
			{
				ID:        "7213fea8e94b4a5593d507237e5a555b",
				Amount:    1,
				Address:   "0x6915f16f8791d0a1cc2bf47c13a6b2a92000504b",
				Asset:     "ETH",
				TxID:      "0xdf33b22bdb2b28b1f75ccd201a4a4m6e7g83jy5fc5d5a9d1340961598cfcb0a1",
				ApplyTime: 1525754126000,
				Status:    4,
				TxFee:     12131,
			},
			{
				ID:        "7213fea8e94b4a5534ggsd237e5a555b",
				Amount:    1000,
				Address:   "463tWEBn5XZJSxLU34r6g7h8jtxuNcDbjLSjkn3XAXHCbLrTTErJrBWYgHJQyrCwkNgYvyV3z8zctJLPCZy24jvb3NiTcTJ",
				Asset:     "XMR",
				TxID:      "b3c6219639c8ae3f9cf010cdc24fw7f7yt8j1e063f9b4bd1a05cb44c4b6e2509",
				ApplyTime: 1525754127000,
				Status:    4,
				TxFee:     31231,
			},
		}

		tests = []httputil.HTTPTestCase{
			{
				Msg:      "get an existing test record for huobi",
				Endpoint: "/withdrawals",
				Params: map[string]string{
					"from": strconv.FormatUint(testTimestamp-10000, 10),
					"to":   strconv.FormatUint(testTimestamp+10000, 10),
					"cex":  "huobi",
				},
				Method: http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					t.Helper()
					require.Equal(t, http.StatusOK, resp.Code)
					var result response
					err := json.NewDecoder(resp.Body).Decode(&result)
					require.NoError(t, err)
					assert.Equal(t, response{
						Huobi: map[string][]huobi.WithdrawHistory{
							"huobi_v1_main": huobiTestData,
						},
					}, result)
				},
			},
			{
				Msg:      "get an existing test record for all exchanges",
				Endpoint: "/withdrawals",
				Params: map[string]string{
					"from": strconv.FormatUint(testTimestamp-10000, 10),
					"to":   strconv.FormatUint(testTimestamp+10000, 10),
				},
				Method: http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					t.Helper()
					require.Equal(t, http.StatusOK, resp.Code)
					var result response
					err := json.NewDecoder(resp.Body).Decode(&result)
					require.NoError(t, err)
					assert.Equal(t, response{
						Huobi: map[string][]huobi.WithdrawHistory{
							"huobi_v1_main": huobiTestData,
						},
						Binance: map[string][]BinanceWithdrawalResponse{
							"binance_1": binanceExpectedResponse,
						},
					}, result)
				},
			},
			{
				Msg:      "get an invalid exchange name",
				Endpoint: "/withdrawals",
				Params: map[string]string{
					"from": strconv.FormatUint(testTimestamp-10000, 10),
					"to":   strconv.FormatUint(testTimestamp+10000, 10),
					"cex":  "huoxbxix",
				},
				Method: http.MethodGet,
				Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
					t.Helper()
					require.Equal(t, http.StatusBadRequest, resp.Code)
				},
			},
		}
		err error
	)

	err = hdb.UpdateWithdrawHistory(huobiTestData, "huobi_v1_main")
	require.NoError(t, err)

	err = bdb.UpdateWithdrawHistory(binanceTestData, "binance_1")
	require.NoError(t, err)

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, ts.r) })
	}

}
func TestMain(m *testing.M) {
	var err error
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	tdb, teardown = testutil.MustNewDevelopmentDB()

	hdb, err = huobiPostgres.NewDB(sugar, tdb)
	if err != nil {
		log.Fatal(err)
	}

	bdb, err = withdrawalstorage.NewDB(sugar, tdb)
	if err != nil {
		log.Fatal(err)
	}

	ts, err = NewServer("", hdb, bdb, sugar)
	if err != nil {
		log.Fatal(err)
	}
	ts.register()

	ret := m.Run()
	if err = teardown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(ret)
}
