package http

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
	"github.com/KyberNetwork/reserve-stats/reserverates/storage"
	influxRateStorage "github.com/KyberNetwork/reserve-stats/reserverates/storage/influx"
)

const (
	host           = "http://127.0.0.1:9001"
	testInfluxURL  = "http://127.0.0.1:8086"
	testRsvAddress = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	dbName         = "test_reserve_rate"
	testFromBlock  = 123
	testToBlock    = 124
)

var (
	testRates = map[string]map[string]common.ReserveRateEntry{
		testRsvAddress: {
			"ETH-KNC": {
				BuyReserveRate:  1,
				BuySanityRate:   2,
				SellReserveRate: 3,
				SellSanityRate:  4,
			},
			"ETH-ZRX": {
				BuyReserveRate:  5,
				BuySanityRate:   6,
				SellReserveRate: 7,
				SellSanityRate:  8,
			},
		},
	}
	testTS = timeutil.TimeToTimestampMs(time.Now())
)

func newTestServer(sugar *zap.SugaredLogger, dbInstance storage.ReserveRatesStorage) (*Server, error) {
	if err := dbInstance.UpdateRatesRecords(123, testRates); err != nil {
		return nil, err
	}
	return NewServer(host, dbInstance, sugar)
}

func tearDownTestDB(t *testing.T, influxClient client.Client) {
	cmd := fmt.Sprintf("DROP DATABASE %s", dbName)
	_, err := influxClient.Query(client.Query{
		Command: cmd,
	})
	assert.Nil(t, err, "rate storage test db should be teardown")
}

func TestHTTPRateServer(t *testing.T) {
	const requestEndpoint = "reserve-rates"

	sugar := testutil.MustNewDevelopmentSugaredLogger()

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: testInfluxURL,
	})
	assert.Nil(t, err, "influx client should be created successfully")

	dbInstance, err := influxRateStorage.NewRateInfluxDBStorage(sugar, influxClient, dbName, blockchain.NewMockBlockTimeResolve(timeutil.TimestampMsToTime(testTS)))
	assert.Nil(t, err, "Rate storage should be created successfully")

	defer tearDownTestDB(t, influxClient)

	server, err := newTestServer(sugar, dbInstance)
	assert.Nil(t, err, "test server should be created succesfully")

	server.register()

	var tests = []httputil.HTTPTestCase{
		{
			Msg:      "success query",
			Endpoint: fmt.Sprintf("%s/%s?from=%d&to=%d&reserve=%s", host, requestEndpoint, testTS, testTS, testRsvAddress),
			Method:   http.MethodGet,
			Assert:   expectCorrectRate,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.Msg, func(t *testing.T) { httputil.RunHTTPTestCase(t, tc, server.r) })
	}
}
