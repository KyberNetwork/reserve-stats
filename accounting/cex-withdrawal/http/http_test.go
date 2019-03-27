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

	huobiPostgres "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history/postgres"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

var (
	ts  *Server
	hdb *huobiPostgres.HuobiStorage
	tdb *sqlx.DB
)

func TestGetHuobiWithdrawal(t *testing.T) {
	t.Log("creating a test get huobi withdrawal test")
	var (
		createdAt uint64 = 1525754125590
		testData         = []huobi.WithdrawHistory{
			{
				ID:         2272335,
				CreatedAt:  createdAt,
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
		tests = []httputil.HTTPTestCase{
			{
				Msg:      "get an existing test record",
				Endpoint: "/withdrawals",
				Params: map[string]string{
					"from": strconv.FormatUint(createdAt-10, 10),
					"to":   strconv.FormatUint(createdAt+10, 10),
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
						Huobi: testData,
					}, result)
				},
			},
			{
				Msg:      "get an invalid exchange name",
				Endpoint: "/withdrawals",
				Params: map[string]string{
					"from": strconv.FormatUint(createdAt-10, 10),
					"to":   strconv.FormatUint(createdAt+10, 10),
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
	err = hdb.UpdateWithdrawHistory(testData)
	require.NoError(t, err)
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, ts.r) })
	}

}
func TestMain(m *testing.M) {
	var err error
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	_, tdb = testutil.MustNewDevelopmentDB()

	hdb, err = huobiPostgres.NewDB(sugar, tdb, huobiPostgres.WithWithdrawalTableName("test_cex_withdrawal_api"))
	if err != nil {
		log.Fatal(err)
	}
	ts, err = NewServer("", hdb, nil, sugar)
	if err != nil {
		log.Fatal(err)
	}
	ts.register()

	ret := m.Run()
	err = hdb.TearDown()
	if err != nil {
		log.Fatal(err)
	}
	err = hdb.Close()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(ret)
}
