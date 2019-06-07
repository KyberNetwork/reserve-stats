package chainalysis

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	chainAlysisURLFlag    = "chain-alysis-url"
	defaultChainAlysisURL = "https://api.chainalysis.com/api/kyt/v1"
	chainAlysisAPIKeyFlag = "chain-alysis-api-key"
)

// NewChainAlysisCliFlags returns cli with app prefix flags to configure read keypair.
func NewChainAlysisCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   chainAlysisURLFlag,
			Usage:  "fetch trade log",
			EnvVar: "CHAIN_ALYSIS_URL",
			Value:  defaultChainAlysisURL,
		},
		cli.StringFlag{
			Name:   chainAlysisAPIKeyFlag,
			Usage:  "api key to interate with chain alysis api",
			EnvVar: "CHAIN_ALYSIS_API_KEY",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	chainAlysisURL := c.String(chainAlysisURLFlag)
	err := validation.Validate(chainAlysisURL,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid trade log url: %q, error: %s", chainAlysisURL, err)
	}

	return NewChainAlysisClient(sugar, chainAlysisURL, c.String(chainAlysisAPIKeyFlag)), nil
}
