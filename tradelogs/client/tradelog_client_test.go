package client

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"

	ethereum "github.com/ethereum/go-ethereum/common"
)

const (
	tradeLogAccessKeyID     = "read-key-id"
	tradeLogSecretAccessKey = "xx232425xx"

	userAddress     = "0x472dbf5a1b070f9efc2491cb3b98445e06599e21"
	receiverAddress = "0x63825c174ab367968ec60f061753d3bbd36a0d8f"

	fromTime = 123
	toTime   = 234
)

func newTestTradeLog(server *httptest.Server) *Client {
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	return NewTradeLogClient(sugar, server.URL, WithAuth(tradeLogAccessKeyID, tradeLogSecretAccessKey))
}

func TestValidGetTradeLog(t *testing.T) {
	var log = []common.TradeLog{
		{
			EthAmount:       new(big.Int),
			UserAddress:     ethereum.HexToAddress(userAddress),
			ReceiverAddress: ethereum.HexToAddress(receiverAddress),
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		js, err := json.Marshal(log)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if req.URL.String() != fmt.Sprintf("/trade-logs?from=%d&to=%d", fromTime, toTime) {
			t.Error("Request to wrong endpoint", "result", req.URL.String())
		}

		rw.WriteHeader(http.StatusOK)
		if _, err := rw.Write(js); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	tl := newTestTradeLog(server)
	tradeLogs, err := tl.GetTradeLogs(fromTime, toTime)
	if err != nil {
		t.Error("Could not get trade logs")
	}
	if len(tradeLogs) == 0 {
		t.Error("tradeLogs should not be empty")
	}
	l := tradeLogs[0]
	if strings.ToLower(l.UserAddress.Hex()) != userAddress {
		t.Error("Get invalid user address", "result", strings.ToLower(l.UserAddress.Hex()), "expected", userAddress)
	}
	if strings.ToLower(l.ReceiverAddress.Hex()) != receiverAddress {
		t.Error("Get invalid receiver address", "result", strings.ToLower(l.ReceiverAddress.Hex()), "expected", receiverAddress)
	}
	if l.EthAmount.Cmp(new(big.Int)) != 0 {
		t.Error("Get invalid eth amount", "result", l.EthAmount, "expected", new(big.Int))
	}
}
