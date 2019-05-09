package chainalysis

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	address = "0x472dbf5a1b070f9efc2491cb3b98445e06599e21"
	txHash  = "0x2d0709e63572bf1ece091048b3fa15eb62c18427f109b3e540c69cf9b0da4cc1"

	apiKey = "xx232425xx"
)

func newTestChainAlysis(server *httptest.Server, apiKey string) *Client {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	return NewChainAlysisClient(sugar, server.URL, apiKey)
}

func TestValidRegisterWithdrawalAddress(t *testing.T) {
	var rw = []registerWithdrawal{
		{
			Asset:   ethSymbol,
			Address: address,
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.URL.String()) != fmt.Sprintf("/users/%s/withdrawaladdresses", address) {
			t.Error("Request to wrong endpoint", "result", req.URL.String())
		}
		if req.Header.Get("Token") != apiKey {
			t.Error("Invalid api key", "api key", req.Header.Get("Token"))
		}
		rw.WriteHeader(http.StatusOK)
		if _, err := rw.Write([]byte{}); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	ca := newTestChainAlysis(server, apiKey)
	err := ca.registerWithdrawalAddress(ethereum.HexToAddress(address), rw)
	if err != nil {
		t.Error("Cannot register withdrawal")
	}
}

func TestValidRegisterSentTransfer(t *testing.T) {
	var rst = []registerSentTransfer{
		{
			Asset:             ethSymbol,
			TransferReference: txHash,
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.ToLower(req.URL.String()) != fmt.Sprintf("/users/%s/transfers/sent", address) {
			t.Error("Request to wrong endpoint", "result", req.URL.String())
		}
		if req.Header.Get("Token") != apiKey {
			t.Error("Invalid api key", "api key", req.Header["Token"])
		}
		rw.WriteHeader(http.StatusOK)
		if _, err := rw.Write([]byte{}); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	ca := newTestChainAlysis(server, apiKey)
	err := ca.registerSentTransfer(ethereum.HexToAddress(address), rst)
	if err != nil {
		t.Error("Cannot register withdrawal")
	}
}
