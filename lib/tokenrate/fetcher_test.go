package tokenrate

import (
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/tokenrate/coingecko"
	"github.com/influxdata/influxdb/client/v2"
)

func TestSaveTokenRate(t *testing.T) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		t.Fatal(err)
	}

	fetcher, err := NewETHUSDRateFetcher(
		sugar,
		"test_db",
		influxClient,
		coingecko.New(),
	)
	if err != nil {
		t.Fatal(err)
	}
	rate := ETHUSDRate{
		BlockNumber: 12314,
		Timestamp:   time.Now(),
		Rate:        123.234,
		Provider:    "testProvider",
	}
	if err := fetcher.SaveTokenRate(rate); err != nil {
		t.Fatal(err)
	}
}
