package etherscanclient

import (
	"github.com/nanmu42/etherscan-api"
	"github.com/urfave/cli"
)

const (
	etherscanAPIKeyFlag  = "etherscan-api-key"
	etherscanNetworkFlag = "etherscan-network"
)

//NewCliFlags return cli flags to configure cex client
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   etherscanAPIKeyFlag,
			Usage:  "api key for etherscan client",
			EnvVar: "ETHERSCAN_API_KEY",
		},
		cli.StringFlag{
			Name:   etherscanNetworkFlag,
			Usage:  "network for etherscan client, value: api-ropsten, api-rinkedby, etc",
			EnvVar: "ETHERSCAN_NETWORK",
			Value:  string(etherscan.Mainnet),
		},
	}
}

//NewEtherscanClientFromContext return etherscan client
func NewEtherscanClientFromContext(c *cli.Context) (*etherscan.Client, error) {
	apiKey := c.String(etherscanAPIKeyFlag)
	network := c.String(etherscanNetworkFlag)
	return etherscan.New(etherscan.Network(network), apiKey), nil
}
