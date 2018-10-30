package storage

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

func doInfluxHTTPReq(client http.Client, cmd, endpoint, db string) error {
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("q", cmd)
	q.Add("db", db)
	req.URL.RawQuery = q.Encode()
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code, expected: %d, got: %d", http.StatusOK, rsp.StatusCode)
	}
	return nil
}

func aggregationTestData(is *InfluxStorage) error {
	const (
		DstAmountHourCmd = "SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY \"dst_addr\", time(1h)"
		DstAmountDayCmd  = "SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE dst_addr!='' GROUP BY \"dst_addr\", time(1d)"
		SrcAmountHourCmd = "SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY \"src_addr\", time(1h)"
		SrcAmountDayCmd  = "SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE src_addr!=''GROUP BY \"src_addr\", time(1d)"
		endpoint         = "http://127.0.0.1:8086/"
	)
	if _, err := is.queryDB(is.influxClient, DstAmountHourCmd); err != nil {
		return err
	}
	if _, err := is.queryDB(is.influxClient, DstAmountDayCmd); err != nil {
		return err
	}
	if _, err := is.queryDB(is.influxClient, SrcAmountHourCmd); err != nil {
		return err
	}
	if _, err := is.queryDB(is.influxClient, SrcAmountDayCmd); err != nil {
		return err
	}
	return nil
}

func TestGetAssetVolume(t *testing.T) {
	const (
		dbName = "test_volume"
		// These params are expected to be change when export.dat changes.
		fromTime  = 1539248043000
		toTime    = 1539248666000
		ethAmount = 238.33849929550047
		freq      = "h"
		timeStamp = "2018-10-11T09:00:00Z"
	)

	is, err := newTestInfluxStorage(dbName)
	assert.NoError(t, err)

	defer func() {
		assert.NoError(t, is.tearDown())
	}()
	assert.NoError(t, loadTestData(dbName))
	assert.NoError(t, aggregationTestData(is))
	volume, err := is.GetAssetVolume(core.ETHToken, fromTime, toTime, freq)
	assert.NoError(t, err)

	t.Logf("Voume result %v", volume)

	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	assert.NoError(t, err)
	timeUint := timeutil.TimeToTimestampMs(timeUnix)
	result, ok := volume[timeUint]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	if result.USDAmount != ethAmount {
		t.Fatal(fmt.Errorf("Expect USD amount to be %.18f, got %.18f", ethAmount, result.USDAmount))
	}
}
