package etherscan

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/nanmu42/etherscan-api"
	"github.com/urfave/cli"
	"golang.org/x/time/rate"
)

const (
	etherscanAPIKeyFlag       = "etherscan-api-key"
	etherscanNetworkFlag      = "etherscan-network"
	etherscanRequestPerSecond = "etherscan-requests-per-second"
)

// NewCliFlags returns cli flags for Etherscan client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   etherscanAPIKeyFlag,
			Usage:  "Etherscan API Key",
			EnvVar: "ETHERSCAN_API_KEY",
		},
		cli.StringFlag{
			Name:   etherscanNetworkFlag,
			Usage:  "Etherscan Network to operate on, valid choices: api, api-ropsten, api-kovan, api-rinkeby, api-tobalaba",
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

func limitRate(limiter *rate.Limiter, semaphore chan struct{}) func(module, action string, param map[string]interface{}) error {
	return func(module, action string, param map[string]interface{}) error {
		semaphore <- struct{}{}
		return limiter.Wait(context.Background())
	}
}

// NewEtherscanClientFromContext returns Ethereum client from flag variable, or error if occurs
func NewEtherscanClientFromContext(c *cli.Context) (*etherscan.Client, error) {
	apiKey := c.String(etherscanAPIKeyFlag)
	network := c.String(etherscanNetworkFlag)
	rps := c.Float64(etherscanRequestPerSecond)
	if rps <= 0 {
		return nil, errors.New("rate limit must be more than 0")
	}
	//Etherscan doesn't  allow burst,  i.e: 5 request per second  really mean 1 request per 0.2  second
	limiter := rate.NewLimiter(rate.Limit(rps), 1)

	rpsInt := int(math.Round(rps))
	// semaphore by buffer channel
	semaphore := make(chan struct{}, rpsInt)

	switch etherscan.Network(network) {
	case etherscan.Mainnet, etherscan.Ropsten, etherscan.Kovan, etherscan.Tobalaba:
		client := etherscan.New(etherscan.Network(network), apiKey)
		client.BeforeRequest = limitRate(limiter, semaphore)
		client.AfterRequest = func(module, action string, param map[string]interface{}, outcome interface{}, requestErr error) {
			<-semaphore
		}
		return client, nil
	default:
		return nil, fmt.Errorf("unknown network: %s", network)
	}
}
