package testutil

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/KyberNetwork/reserve-stats/lib/node"
)

const ethereumNodeEnv = "ETHEREUM_NODE"

// MustNewDevelopmentwEthereumClient creates a new Ethereum client to use in tests.
func MustNewDevelopmentwEthereumClient() *ethclient.Client {
	endpoint, ok := os.LookupEnv(ethereumNodeEnv)
	if !ok {
		endpoint = node.InfuraEndpoint()
	}
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	return client
}
