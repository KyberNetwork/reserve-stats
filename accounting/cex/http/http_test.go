package http

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tradestorage "github.com/KyberNetwork/reserve-stats/accounting/binance-storage/trade-storage"
	huobistorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/postgres"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

var (
	ts          *Server
	huobiTrades = []huobi.TradeHistory{
		{
			ID:              59378,
			Symbol:          "ethusdt",
			AccountID:       100009,
			Amount:          "10.1000000000",
			Price:           "100.1000000000",
			CreatedAt:       1494901162595,
			Type:            "buy-limit",
			FieldAmount:     "10.1000000000",
			FieldCashAmount: "1011.0100000000",
			FieldFees:       "0.0202000000",
			FinishedAt:      1494901400468,
			UserID:          1000,
			Source:          "api",
			State:           "filled",
			CanceledAt:      0,
			Exchange:        "huobi",
			Batch:           "",
		},
	}
	binanceTrades = []binance.TradeHistory{
		{
			Symbol:   "BNBBTC",
			ID:       28457,
			OrderID:  100234,
			Price:    "4.00000100",
			Quantity: "12.00000000",
			// TODO: missing QuoteQty
			Commission:      "10.10000000",
			CommissionAsset: "BNB",
			Time:            1499865549590,
			IsBuyer:         true,
			IsMaker:         false,
			IsBestMatch:     false,
		},
	}
)

func TestTrades(t *testing.T) {
	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "get getTrades from all exchanges",
			Endpoint: "/cex_trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1494901162000",
				"to":   "1499865549600",
			},
			Body: nil,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, resp.Code)

				var trades getTradesResponse
				err := json.NewDecoder(resp.Body).Decode(&trades)
				require.NoError(t, err)
				assert.Equal(t, getTradesResponse{
					Huobi:   huobiTrades,
					Binance: binanceTrades,
				}, trades)
			},
		},
		{
			Msg:      "get empty trades with newer time range",
			Endpoint: "/cex_trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1499865549600",
				"to":   "1499865559600",
			},
			Body: nil,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, resp.Code)

				var trades getTradesResponse
				err := json.NewDecoder(resp.Body).Decode(&trades)
				require.NoError(t, err)
				assert.Len(t, trades.Huobi, 0)
				assert.Len(t, trades.Binance, 0)
			},
		},
		{
			Msg:      "get empty trades with older time range",
			Endpoint: "/cex_trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1494901160000",
				"to":   "1494901162000",
			},
			Body: nil,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, resp.Code)

				var trades getTradesResponse
				err := json.NewDecoder(resp.Body).Decode(&trades)
				require.NoError(t, err)
				assert.Len(t, trades.Huobi, 0)
				assert.Len(t, trades.Binance, 0)
			},
		},
		{
			Msg:      "get getTrades from huobi",
			Endpoint: "/cex_trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1494901162000",
				"to":   "1499865549600",
				"cex":  "huobi",
			},
			Body: nil,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, resp.Code)

				var trades getTradesResponse
				err := json.NewDecoder(resp.Body).Decode(&trades)
				require.NoError(t, err)
				assert.Equal(t, getTradesResponse{
					Huobi: huobiTrades,
				}, trades)
			},
		},
		{
			Msg:      "get getTrades from binance",
			Endpoint: "/cex_trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1494901162000",
				"to":   "1499865549600",
				"cex":  "binance",
			},
			Body: nil,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, resp.Code)

				var trades getTradesResponse
				err := json.NewDecoder(resp.Body).Decode(&trades)
				require.NoError(t, err)
				assert.Equal(t, getTradesResponse{
					Binance: binanceTrades,
				}, trades)
			},
		},
		{
			Msg:      "get getTrades with invalid exchange",
			Endpoint: "/cex_trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1494901162000",
				"to":   "1499865549600",
				"cex":  "world_bank",
			},
			Body: nil,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, ts.r) })
	}
}

func TestMain(m *testing.M) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	_, db := testutil.MustNewDevelopmentDB()

	hs, err := huobistorage.NewDB(sugar, db, huobistorage.WithTradeTableName("huobi_cex_http_test"))
	if err != nil {
		log.Fatal(err)
	}

	if err = hs.UpdateTradeHistory(huobiTrades); err != nil {
		log.Fatal(err)
	}

	bs, err := tradestorage.NewDB(sugar, db, "binance_cex_http_test")
	if err != nil {
		log.Fatal(err)
	}

	if err = bs.UpdateTradeHistory(binanceTrades); err != nil {
		log.Fatal(err)
	}

	ts = NewServer(sugar, "", hs, bs)
	ts.register()

	ret := m.Run()

	if err = hs.TearDown(); err != nil {
		log.Fatal(err)
	}

	if err = bs.DeleteTable(); err != nil {
		log.Fatal(err)
	}

	os.Exit(ret)
}
