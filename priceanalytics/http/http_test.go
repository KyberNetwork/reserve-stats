package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/priceanalytics/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHTTPPriceAnalyticServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	assert.Nil(t, err, "logger should be initiated successfully")

	sugar := logger.Sugar()
	db, teardown := testutil.MustNewDevelopmentDB()
	priceStorage, err := storage.NewPriceStorage(sugar, db)
	assert.Nil(t, err, "price storage should be initiated successfully")

	defer func() {
		assert.NoError(t, teardown())
	}()

	s := NewHTTPServer(sugar, "", priceStorage)
	s.register()

	const (
		requestEnpoint = "/price-analytics"
		queryTimestamp = 1541478232000
	)

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "missing timestamp update",
			Endpoint: requestEnpoint,
			Method:   http.MethodPost,
			Body: []byte(`
			{
				"block_expiration": true,
				"triggering_tokens_list": [
					{
						"token": "KNC",
						"ask_price": 0.123,
						"bid_price": 0.125,
						"mid_afp_price": 0.124,
						"mid_afp_old_price": 0.12,
						"min_spread": 0.002,
						"trigger_update": true
					},
					{
						"token": "OMG",
						"ask_price": 0.123,
						"bid_price": 0.125,
						"mid_afp_price": 0.124,
						"mid_afp_old_price": 0.12,
						"min_spread": 0.002,
						"trigger_update": false
					}
				]
			}
			`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, resp.Code)
			},
		},
		{
			Msg:      "success update",
			Endpoint: requestEnpoint,
			Method:   http.MethodPost,
			Body: []byte(`
			{
				"timestamp": 1541478232000,
				"block_expiration": true,
				"triggering_tokens_list": [
					{
						"token": "KNC",
						"ask_price": 0.123,
						"bid_price": 0.125,
						"mid_afp_price": 0.124,
						"mid_afp_old_price": 0.12,
						"min_spread": 0.002,
						"trigger_update": true
					},
					{
						"token": "OMG",
						"ask_price": 0.123,
						"bid_price": 0.125,
						"mid_afp_price": 0.124,
						"mid_afp_old_price": 0.12,
						"min_spread": 0.002,
						"trigger_update": false
					}
				]
			}
			`),
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
			},
		},
		{
			Msg:      "success get",
			Endpoint: fmt.Sprintf("%s?from=%d", requestEnpoint, queryTimestamp),
			Method:   http.MethodGet,
			Assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, resp.Code)
				body, _ := ioutil.ReadAll(resp.Body)
				t.Log(string(body))
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, s.r) })
	}
}
