package storage

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/core"
	// "github.com/KyberNetwork/reserve-stats/lib/timeutil"
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

func aggregationTestData(db string) error {
	const (
		DstAmountHourCmd = "SELECT SUM(dst_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT dst_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY \"dst_addr\", time(1h)"
		DstAmountDayCmd  = "SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE dst_addr!='' GROUP BY \"dst_addr\", time(1d)"
		SrcAmountHourCmd = "SELECT SUM(src_amount) AS token_volume, SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume INTO volume_hour FROM (SELECT src_amount, eth_amount, eth_amount*eth_usd_rate AS usd_amount FROM trades WHERE (src_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' AND dst_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2') OR (src_addr!='0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' AND dst_addr!='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE')) GROUP BY \"src_addr\", time(1h)"
		SrcAmountDayCmd  = "SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, SUM(usd_volume) AS usd_volume INTO volume_day FROM volume_hour WHERE src_addr!=''GROUP BY \"src_addr\", time(1d)"
		endpoint         = "http://127.0.0.1:8086/"
	)
	client := http.Client{Timeout: time.Second * 5}
	if err := doInfluxHTTPReq(client, DstAmountHourCmd, endpoint+"query", db); err != nil {
		log.Printf("fuck")
		return err
	}
	if err := doInfluxHTTPReq(client, DstAmountDayCmd, endpoint+"query", db); err != nil {
		return err
	}
	if err := doInfluxHTTPReq(client, SrcAmountHourCmd, endpoint+"query", db); err != nil {
		return err
	}
	if err := doInfluxHTTPReq(client, SrcAmountDayCmd, endpoint+"query", db); err != nil {
		return err
	}
	return nil
}

func TestGetAssetVolume(t *testing.T) {
	const dbName = "test_volume"
	is, err := newTestInfluxStorage(dbName)
	// defer func() {
	// 	if err := is.tearDown(); err != nil {
	// 		t.Fatal(err)
	// 	}
	// }()
	if err := loadTestData(dbName); err != nil {
		t.Fatal(err)
	}

	if err := aggregationTestData(dbName); err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	// 2018-10-11T16:00:00+07:00
	// These number is expected to be change when export.dat changes.
	volume, err := is.GetAssetVolume(core.ETHToken, 1539248043000, 1539248666000, "h")
	for k, _ := range volume {
		log.Printf("%s", k.Format(time.RFC3339))
	}
	log.Printf("%v", volume)
	if err != nil {
		t.Fatal(err)
	}
	timeStamp, err := time.Parse(time.RFC3339, "2018-10-11T09:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	result, ok := volume[timeStamp]
	if !ok {
		t.Fatalf("expect to find result at timestamp %s, yet there is none", timeStamp.Format(time.RFC3339))
	}

	if result.USDAmount != 238.33849929550047 {
		t.Fatal(fmt.Errorf("Expect USD amount to be 238.33849929550047, got %.f", result.USDAmount))
	}
}

// SELECT SUM(token_volume) as token_volume, SUM(eth_volume) as eth_volume, sum(usd_volume) as usd_volume FROM volume_hour WHERE (dst_addr='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE' OR src_addr='0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE') GROUP BY time(1h) fill(none)
// name: volume_hour
// time                 token_volume       eth_volume         usd_volume
// ----                 ------------       ----------         ----------
// 2018-10-11T08:00:00Z 16.333613369193138 16.333613369193138 3685.8761243552653
// 2018-10-11T09:00:00Z 1.0561746426481893 1.0561746426481893 238.33849929550047
