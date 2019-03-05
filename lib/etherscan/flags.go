package etherscanclient

import (
	"context"
	"errors"

	"github.com/nanmu42/etherscan-api"
	"github.com/urfave/cli"
	"golang.org/x/time/rate"
)

const (
	etherscanAPIKeyFlag       = "etherscan-api-key"
	etherscanNetworkFlag      = "etherscan-network"
	etherscanRequestPerSecond = "etherscan-requests-per-second"
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
		cli.Float64Flag{
			Name:   etherscanRequestPerSecond,
			Usage:  "etherscan request limit per second, default to 5 which etherscan's normal rate limit",
			EnvVar: "ETHERSCAN_REQUEST_PER_SEC",
			Value:  5,
		},
	}
}

func limitRate(limiter *rate.Limiter) func(module, action string, param map[string]interface{}) error {
	return func(module, action string, param map[string]interface{}) error {
		return limiter.Wait(context.Background())
	}
}

//NewEtherscanClientFromContext return etherscan client
func NewEtherscanClientFromContext(c *cli.Context) (*etherscan.Client, error) {
	apiKey := c.String(etherscanAPIKeyFlag)
	network := c.String(etherscanNetworkFlag)
	client := etherscan.New(etherscan.Network(network), apiKey)
	rps := c.Float64(etherscanRequestPerSecond)
	if rps <= 0 {
		return nil, errors.New("rate limit must be more than 0")
	}
	//Etherscan doesn't  allow burst,  i.e: 5 request per second  really mean 1 request per 0.2  second
	limiter := rate.NewLimiter(rate.Limit(rps), 1)
	client.BeforeRequest = limitRate(limiter)
	return client, nil
}
