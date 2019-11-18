package coingecko

import (
	"time"

	"github.com/urfave/cli"
)

const (
	reqTimeWaitingFlag    = "coingecko-req-waiting-time"
	defaultReqTimeWaiting = time.Second
)

// NewFlags return cli config for coingecko
func NewFlags() []cli.Flag {
	return []cli.Flag{
		cli.DurationFlag{
			Name:   reqTimeWaitingFlag,
			Usage:  "waiting time after each request to avoid rate limit",
			EnvVar: "COINGECKO_REQ_WAITING_TIME",
			Value:  defaultReqTimeWaiting,
		},
	}
}

// NewCoinGeckoFromContext return coingecko provider
func NewCoinGeckoFromContext(c *cli.Context) *CoinGecko {
	return New(c.Duration(reqTimeWaitingFlag))
}
