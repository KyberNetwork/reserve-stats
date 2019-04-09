package huobi

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

const (
	huobiAPIKeyFlag       = "huobi-api-key"
	huobiSecretKeyFlag    = "huobi-secret-key"
	huobiRequestPerSecond = "huobi-requests-per-second"
	huobiClientValidation = "huobi-client-validation"
)

//NewCliFlags return cli flags to configure cex-trade client
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
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
		cli.Float64Flag{
			Name:   huobiRequestPerSecond,
			Usage:  "huobi request limit per second, default to 8 which huobi's tested rate limit",
			EnvVar: "HUOBI_REQUESTS_PER_SECOND",
			Value:  8,
		},
		cli.BoolTFlag{
			Name:   huobiClientValidation,
			Usage:  "if set to true, the client is validate by calling GetAccounts with its API key",
			EnvVar: "HUOBI_CLIENT_VALIDATION",
		},
	}
}

// NewClientFromContext return huobi client
func NewClientFromContext(c *cli.Context, sugar *zap.SugaredLogger) (*Client, error) {
	var (
		apiKey, secretKey string
		options           []Option
	)
	if c.String(huobiAPIKeyFlag) == "" {
		return nil, fmt.Errorf("cannot create huobi client, lack of api key")
	}
	apiKey = c.String(huobiAPIKeyFlag)

	if c.String(huobiSecretKeyFlag) == "" {
		return nil, fmt.Errorf("cannot create huobi client, lack of secret key")
	}
	secretKey = c.String(huobiSecretKeyFlag)

	rps := c.Float64(huobiRequestPerSecond)
	if rps <= 0 {
		return nil, errors.New("request per second must be greater than 0")
	}

	options = append(options, WithRateLimiter(rate.NewLimiter(rate.Limit(rps), 1)))
	if validateRequire := c.BoolT(huobiClientValidation); validateRequire {
		options = append(options, WithValidation())
	}
	return NewClient(apiKey, secretKey, sugar, options...)
}
