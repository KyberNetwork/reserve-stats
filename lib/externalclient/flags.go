package externalclient

import (
	"fmt"
	"log"

	"github.com/nanmu42/etherscan-api"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	binanceAPIKeyFlag    = "binance-api-key"
	binanceSecretKeyFlag = "binance-secret-key"
	huobiAPIKeyFlag      = "huobi-api-key"
	huobiSecretKeyFlag   = "huobi-secret-key"
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
			Name:   huobiAPIKeyFlag,
			Usage:  "API key for huobi client",
			EnvVar: "HUOBI_API_KEY",
		},
		cli.StringFlag{
			Name:   huobiSecretKeyFlag,
			Usage:  "secret key for huobi client",
			EnvVar: "HUOBI_SECRET_KEY",
		},
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

//NewBinanceClientFromContext return binance client
func NewBinanceClientFromContext(c *cli.Context) (*BinanceClient, error) {
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

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return NewBinanceClient(apiKey, secretKey, sugar), nil
}

//NewHuobiClientFromContext return huobi client
func NewHuobiClientFromContext(c *cli.Context) (*HuobiClient, error) {
	var (
		apiKey, secretKey string
	)
	if c.String(huobiAPIKeyFlag) == "" {
		return nil, fmt.Errorf("cannot create binance client, lack of api key")
	}
	apiKey = c.String(huobiAPIKeyFlag)

	if c.String(huobiSecretKeyFlag) == "" {
		return nil, fmt.Errorf("cannot create binance client, lack of secret key")
	}
	secretKey = c.String(huobiSecretKeyFlag)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return NewHuobiClient(apiKey, secretKey, sugar), nil
}

//NewEtherscanClientFromContext return etherscan client
func NewEtherscanClientFromContext(c *cli.Context) (*etherscan.Client, error) {
	apiKey := c.String(etherscanAPIKeyFlag)
	network := c.String(etherscanNetworkFlag)
	return etherscan.New(etherscan.Network(network), apiKey), nil
}
