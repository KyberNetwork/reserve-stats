package externalclient

import (
	"fmt"

	"github.com/adshao/go-binance"
	"github.com/nanmu42/etherscan-api"
	"github.com/urfave/cli"
)

const (
	binanceAPIKeyFlag    = "binance-api-key"
	binanceSecretKeyFlag = "binance-secret-key"
	etherscanAPIKeyFlag  = "etherscan-api-key"
	etherscanNetworkFlag = "etherscan-network"
)

//NewCliFlags return cli flags to configure cex client
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   binanceAPIKeyFlag,
			Usage:  "API key for binance client",
			EnvVar: "BINANCE_API_KEY",
		},
		cli.StringFlag{
			Name:   binanceSecretKeyFlag,
			Usage:  "secret key for binance client",
			EnvVar: "BINANCE_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   etherscanAPIKeyFlag,
			Usage:  "api key for etherscan client",
			EnvVar: "ETHERSCAN_API_KEY",
		},
		cli.StringFlag{
			Name:   etherscanNetworkFlag,
			Usage:  "network for etherscan client",
			EnvVar: "ETHERSCAN_NETWORK",
			Value:  string(etherscan.Mainnet),
		},
	}
}

//NewBinanceClientFromContext return binance client
func NewBinanceClientFromContext(c *cli.Context) (*binance.Client, error) {
	var (
		apiKey, secretKey string
	)
	if c.String(binanceAPIKeyFlag) == "" {
		return nil, fmt.Errorf("cannot create binance client, lack of api key")
	}
	apiKey = c.String(binanceAPIKeyFlag)

	if c.String(binanceSecretKeyFlag) == "" {
		return nil, fmt.Errorf("cannot create binance client, lack of secret key")
	}

	secretKey = c.String(binanceSecretKeyFlag)
	return binance.NewClient(apiKey, secretKey), nil
}

//NewHuobiClientFromContext return huobi client
//TODO: write our own huobi client because we will need to get withdraw history from
// unofficial api
func NewHuobiClientFromContext(c *cli.Context) {

}

//NewEtherscanClientFromContext return etherscan client
func NewEtherscanClientFromContext(c *cli.Context) (*etherscan.Client, error) {
	apiKey := c.String(etherscanAPIKeyFlag)
	network := c.String(etherscanNetworkFlag)
	return etherscan.New(etherscan.Network(network), apiKey), nil
}
