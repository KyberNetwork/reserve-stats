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

	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/tradestorage"
	huobistorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/postgres"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

var (
	ts          *Server
	huobiTrades = map[int64]huobi.TradeHistory{
		59378: {
			ID:              59378,
			Symbol:          "ethusdt",
			AccountID:       100009,
			Amount:          "10.1000000000",
			Price:           "100.1000000000",
			CreatedAt:       1526428800000,
			Type:            "buy-limit",
			FieldAmount:     "10.1000000000",
			FieldCashAmount: "1011.0100000000",
			FieldFees:       "0.0202000000",
			FinishedAt:      1526428800000,
			UserID:          1000,
			Source:          "api",
			State:           "filled",
			CanceledAt:      0,
			Exchange:        "huobi",
			Batch:           "",
		},
	}
	expectedHuobiTrades = []huobi.TradeHistory{
		{
			ID:              59378,
			Symbol:          "ethusdt",
			AccountID:       100009,
			Amount:          "10.1000000000",
			Price:           "100.1000000000",
			CreatedAt:       1526428800000,
			Type:            "buy-limit",
			FieldAmount:     "10.1000000000",
			FieldCashAmount: "1011.0100000000",
			FieldFees:       "0.0202000000",
			FinishedAt:      1526428800000,
			UserID:          1000,
			Source:          "api",
			State:           "filled",
			CanceledAt:      0,
			Exchange:        "huobi",
			Batch:           "",
		},
	}
	binanceTrades = map[string][]binance.TradeHistory{
		"binance_account_1": []binance.TradeHistory{
			{
				Symbol:          "BNBBTC",
				ID:              28457,
				OrderID:         100234,
				Price:           "4.00000100",
				Quantity:        "12.00000000",
				QuoteQuantity:   "48.000012",
				Commission:      "10.10000000",
				CommissionAsset: "BNB",
				Time:            1528675200000,
				IsBuyer:         true,
				IsMaker:         false,
				IsBestMatch:     false,
			},
		},
	}
)

func TestTrades(t *testing.T) {
	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "get trades from all exchanges",
			Endpoint: "/trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1526428800000",
				"to":   "1528675200001",
			},
			Body: nil,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				require.Equal(t, http.StatusOK, resp.Code)

				var trades getTradesResponse
				err := json.NewDecoder(resp.Body).Decode(&trades)
				require.NoError(t, err)
				assert.Equal(t, getTradesResponse{
					Huobi:   expectedHuobiTrades,
					Binance: binanceTrades,
				}, trades)
			},
		},
		{
			Msg:      "get empty trades with newer time range",
			Endpoint: "/trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1531440000000",
				"to":   "1531440000000",
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
			Endpoint: "/trades",
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
			Endpoint: "/trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1526428700000",
				"to":   "1526428900000",
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
					Huobi: expectedHuobiTrades,
				}, trades)
			},
		},
		{
			Msg:      "get getTrades from binance",
			Endpoint: "/trades",
			Method:   http.MethodGet,
			Params: map[string]string{
				"from": "1528675100000",
				"to":   "1528675300000",
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
			Endpoint: "/trades",
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
	db, teardown := testutil.MustNewDevelopmentDB()

	hs, err := huobistorage.NewDB(sugar, db)
	if err != nil {
		log.Fatal(err)
	}

	if err = hs.UpdateTradeHistory(huobiTrades); err != nil {
		log.Fatal(err)
	}

	bs, err := tradestorage.NewDB(sugar, db)
	if err != nil {
		log.Fatal(err)
	}

	if err = bs.UpdateTradeHistory(binanceTrades["binance_account_1"], "binance_1"); err != nil {
		log.Fatal(err)
	}

	ts = NewServer(sugar, "", hs, bs)
	ts.register()

	ret := m.Run()

	if err = teardown(); err != nil {
		log.Fatal(err)
	}

	os.Exit(ret)
}
