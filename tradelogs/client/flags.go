package client

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	tradeLogURLFlag    = "trade-log-url"
	defaultTradeLogURL = "https://stats-gateway.knstats.com"
)

// NewTradeLogCliFlags returns cli with app prefix flags to configure read keypair.
func NewTradeLogCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   tradeLogURLFlag,
			Usage:  "fetch trade log",
			EnvVar: "TRADE_LOG_URL",
			Value:  defaultTradeLogURL,
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context, options ...TradeLogClientOption) (*TradeLogClient, error) {
	tradeLogURL := c.String(tradeLogURLFlag)
	err := validation.Validate(tradeLogURL,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid trade log url: %q, error: %s", tradeLogURL, err)
	}

	return NewTradeLogClient(sugar, tradeLogURL, options...), nil
}
