package blockchain

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/node"
)

const (
	ethereumNodeFlag = "ethereum-node"
)

// NewEthereumNodeFlags returns cli flag for ethereum node url input
func NewEthereumNodeFlags() cli.Flag {
	return cli.StringFlag{
		Name:   ethereumNodeFlag,
		Usage:  "Ethereum Node URL",
		EnvVar: "ETHEREUM_NODE",
		Value:  node.InfuraEndpoint(),
	}
}

type reqMessage struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      int               `json:"id"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
}

var (
	reserveRateHelper = json.RawMessage(`{
  "0x221CA93C327eDe7d8f9296E2a790905CD7021105": {
    "nonce": "0x112233",
    "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063cd85afd414610030575b600080fd5b61011c6004803603606081101561004657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019064010000000081111561008357600080fd5b82018360208201111561009557600080fd5b803590602001918460208302840111640100000000831117156100b757600080fd5b9091929391929390803590602001906401000000008111156100d857600080fd5b8201836020820111156100ea57600080fd5b8035906020019184602083028401116401000000008311171561010c57600080fd5b90919293919293905050506101bb565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b83811015610163578082015181840152602081019050610148565b50505050905001838103825284818151815260200191508051906020019060200280838360005b838110156101a557808201518184015260208101905061018a565b5050505090500194505050505060405180910390f35b606080838390508686905014610239576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f73726373206c656e67746820213d20646573747300000000000000000000000081525060200191505060405180910390fd5b8585905067ffffffffffffffff8111801561025357600080fd5b506040519080825280602002602001820160405280156102825781602001602082028036833780820191505090505b5091508585905067ffffffffffffffff8111801561029f57600080fd5b506040519080825280602002602001820160405280156102ce5781602001602082028036833780820191505090505b50905060008090505b8686905081101561082357600073ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff166347e6924f6040518163ffffffff1660e01b815260040160206040518083038186803b15801561034057600080fd5b505afa158015610354573d6000803e3d6000fd5b505050506040513d602081101561036a57600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff161461056b578773ffffffffffffffffffffffffffffffffffffffff166347e6924f6040518163ffffffff1660e01b815260040160206040518083038186803b1580156103dc57600080fd5b505afa1580156103f0573d6000803e3d6000fd5b505050506040513d602081101561040657600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff1663a58092b788888481811061043e57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1687878581811061046757fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff166040518363ffffffff1660e01b8152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019250505060206040518083038186803b15801561051757600080fd5b505afa15801561052b573d6000803e3d6000fd5b505050506040513d602081101561054157600080fd5b810190808051906020019092919050505082828151811061055e57fe5b6020026020010181815250505b8773ffffffffffffffffffffffffffffffffffffffff1663d5847d336040518163ffffffff1660e01b815260040160206040518083038186803b1580156105b157600080fd5b505afa1580156105c5573d6000803e3d6000fd5b505050506040513d60208110156105db57600080fd5b810190808051906020019092919050505073ffffffffffffffffffffffffffffffffffffffff1663b8e9c22e73eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee73ffffffffffffffffffffffffffffffffffffffff1689898581811061063e57fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146106a45788888481811061068257fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff166106ce565b8686848181106106b057fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff165b4373eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee73ffffffffffffffffffffffffffffffffffffffff168b8b8781811061070657fe5b9050602002013573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610745576000610748565b60015b60006040518563ffffffff1660e01b8152600401808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018315151515815260200182815260200194505050505060206040518083038186803b1580156107c357600080fd5b505afa1580156107d7573d6000803e3d6000fd5b505050506040513d60208110156107ed57600080fd5b810190808051906020019092919050505083828151811061080a57fe5b60200260200101818152505080806001019150506102d7565b50818191509150955095935050505056fea26469706673582212200863abd4e9f1bed9375f34dc062b9ee1d21ab443fb10f8a736bfeb1e20f1ed1464736f6c63430006080033"
  }
}`)
)

type roundTripperExt struct {
	c *http.Client
}

func (r roundTripperExt) RoundTrip(request *http.Request) (*http.Response, error) {
	rt := request.Clone(context.Background())
	body, _ := ioutil.ReadAll(request.Body)
	// log.Printf("%s \n\n\n\n", body)
	_ = request.Body.Close()
	if len(body) > 0 {
		rt.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}
	var req reqMessage
	if err := json.Unmarshal(body, &req); err == nil {
		if req.Method == "eth_call" /* && (bytes.Contains(req.Params[0], []byte(`0x99f9fbd2`)) || bytes.Contains(req.Params[0], []byte(`0x110bb26c`)))*/ {
			req.Params = append(req.Params, reserveRateHelper)
		}
		d2, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}
		rt.ContentLength = int64(len(d2))
		rt.Body = ioutil.NopCloser(bytes.NewBuffer(d2))
	}
	return r.c.Do(rt)
}

// NewEthereumClientFromFlag returns Ethereum client from flag variable, or error if occurs
func NewEthereumClientFromFlag(c *cli.Context) (*ethclient.Client, error) {
	ethereumNodeURL := c.GlobalString(ethereumNodeFlag)
	cc := &http.Client{Transport: roundTripperExt{c: &http.Client{}}}
	r, err := rpc.DialHTTPWithClient(ethereumNodeURL, cc)
	if err != nil {
		zap.S().Panicw("init custom ethclient failed", "err", err)
	}
	client := ethclient.NewClient(r)
	return client, nil
}

// NodeURLFromFlag ...
func NodeURLFromFlag(c *cli.Context) string {
	return c.GlobalString(ethereumNodeFlag)
}
