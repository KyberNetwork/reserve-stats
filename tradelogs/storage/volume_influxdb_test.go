package storage

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
	ethereum "github.com/ethereum/go-ethereum/common"
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
		endpoint = "http://127.0.0.1:8086/"
	)
	cqs, err := tradelogcq.CreateAssetVolumeCqs(is.dbName)
	if err != nil {
		return err
	}
	for _, cq := range cqs {
		err = cq.Execute(is.influxClient, is.sugar)
		if err != nil {
			return err
		}
	}
	return nil
}

func aggregationVolumeTestData(is *InfluxStorage) error {
	var (
		hourCmds = []string{
			"SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND src_rsv_addr!='') GROUP BY src_addr,src_rsv_addr,time(1h)",
			"SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND dst_rsv_addr!='') GROUP BY src_addr,dst_rsv_addr,time(1h)",
			"SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND src_rsv_addr!='') GROUP BY dst_addr,src_rsv_addr,time(1h)",
			"SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO rsv_volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') AND dst_rsv_addr!='') GROUP BY dst_addr,dst_rsv_addr,time(1h)",
		}
		dayCmds = []string{
			"SELECT SUM(token_volume) as token_volume, SUM(eth_volume) AS eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE src_addr!='' AND src_rsv_addr!='' GROUP BY src_addr,src_rsv_addr,time(1d)",
			"SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE src_addr!='' AND dst_rsv_addr!='' GROUP BY src_addr,dst_rsv_addr,time(1d)",
			"SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE dst_addr!='' and src_rsv_addr!='' GROUP BY dst_addr,src_rsv_addr,time(1d)",
			"SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO rsv_volume_day FROM rsv_volume_hour WHERE dst_addr!='' and dst_rsv_addr!='' GROUP BY dst_addr,dst_rsv_addr,time(1d)",
		}
	)
	for _, cmd := range hourCmds {
		if _, err := is.queryDB(is.influxClient, cmd); err != nil {
			return err
		}
	}
	for _, cmd := range dayCmds {
		if _, err := is.queryDB(is.influxClient, cmd); err != nil {
			return err
		}
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

	t.Logf("Volume result %v", volume)

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

func TestGetReserveVolume(t *testing.T) {
	const (
		dbName = "test_volume"

		// These params are expected to be change when export.dat changes.
		fromTime   = 1539248043000
		toTime     = 1539248666000
		ethAmount  = 227.05539848662738
		freq       = "h"
		timeStamp  = "2018-10-11T09:00:00Z"
		rsvAddrStr = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"
	)

	is, err := newTestInfluxStorage(dbName)
	defer func() {
		if err := is.tearDown(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := loadTestData(dbName); err != nil {
		t.Fatal(err)
	}
	if err := aggregationVolumeTestData(is); err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	volume, err := is.GetReserveVolume(ethereum.HexToAddress(rsvAddrStr), core.ETHToken, fromTime, toTime, freq)
	t.Logf("%v", volume)
	if err != nil {
		t.Fatal(err)
	}
	timeUnix, err := time.Parse(time.RFC3339, timeStamp)
	if err != nil {
		t.Fatal(err)
	}
	result, ok := volume[timeUnix]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeUnix.Format(time.RFC3339))
	}

	if result.USDAmount != ethAmount {
		t.Fatal(fmt.Errorf("Expect USD amount to be %.18f, got %.18f", ethAmount, result.USDAmount))
	}
}
