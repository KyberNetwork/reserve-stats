package coinbase

import (
	"time"

	"github.com/urfave/cli"
)

const (
	reqTimeWaitingFlag    = "coinbase-req-waiting-time"
	defaultReqTimeWaiting = time.Second
)

// NewFlags return cli config for coinbase
func NewFlags() []cli.Flag {
	return []cli.Flag{
		cli.DurationFlag{
			Name:   reqTimeWaitingFlag,
			Usage:  "waiting time after each request to avoid rate limit",
			EnvVar: "COINBASE_REQ_WAITING_TIME",
			Value:  defaultReqTimeWaiting,
		},
	}
}

// NewCoinBaseFromContext return coibase provider
func NewCoinBaseFromContext(c *cli.Context) *CoinBase {
	return New(c.Duration(reqTimeWaitingFlag))
}
