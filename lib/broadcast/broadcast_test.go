package broadcast

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

const (
	readKeyID     = "readKeyID"
	readSecretKey = "xxx123xxx"
)

func newTestGeoInfo(server *httptest.Server) *Client {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	return NewClient(sugar, server.URL, WithAuth(readKeyID, readSecretKey))
}

func TestGetValidResponse(t *testing.T) {
	const (
		tx              = "0x18b7985314631687b09350698d6f8428ab003fa3abc1ce20b8cccfc48cb0700f"
		ipResponse      = "31.166.85.223"
		countryResponse = "SA"
	)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var res tradeLogGeoInfoResp
		res.IP = ipResponse
		res.Country = countryResponse
		js, err := json.Marshal(res)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if req.URL.String() != fmt.Sprintf("/get-tx-info/%s", tx) {
			t.Error("Request to wrong endpoint", "result", req.URL.String())
		}
		rw.Header().Set("Content-Type", "application/json")
		if _, err = rw.Write(js); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	g := newTestGeoInfo(server)
	_, ip, country, err := g.GetTxInfo(tx)
	if err != nil {
		t.Error("Could not get ipInfo")
	}
	if ip != ipResponse {
		t.Error("Get invalid IP", "result", ip, "expected", ipResponse)
	}
	if country != countryResponse {
		t.Error("Get invalid country", "result", country, "expected", countryResponse)
	}
}

func TestInvalidResponse(t *testing.T) {
	const (
		tx = "0x18b7985314631687b09350698d6f8428ab003fa3abc1ce20b8cccfc48cb0700"
	)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "not found", http.StatusNotFound)
	}))

	g := newTestGeoInfo(server)
	_, _, _, err := g.GetTxInfo(tx)
	if err != nil {
		t.Errorf("Get unexpected error: %s", err.Error())
	}
}
