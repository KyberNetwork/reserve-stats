package blockchain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"

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

// NewEthereumClientFromFlag returns Ethereum client from flag variable, or error if occurs
func NewEthereumClientFromFlag(c *cli.Context) (*ethclient.Client, error) {
	ethereumNodeURL := c.GlobalString(ethereumNodeFlag)
	return ethclient.Dial(ethereumNodeURL)
}
