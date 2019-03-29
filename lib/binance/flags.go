package binance

import (
	"errors"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	binanceAPIKeyFlag       = "binance-api-key"
	binanceSecretKeyFlag    = "binance-secret-key"
	binanceRequestPerSecond = "binance-requests-per-second"
)

//NewCliFlags return cli flags to configure cex-trade client
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
		cli.Float64Flag{
			Name:   binanceRequestPerSecond,
			Usage:  "binance request limit per second, default to 20 which etherscan's normal rate limit",
			EnvVar: "BINANCE_REQUESTS_PER_SECOND",
			Value:  10,
		},
	}
}

//NewClientFromContext return binance client
func NewClientFromContext(c *cli.Context, sugar *zap.SugaredLogger) (*Client, error) {
	var (
		apiKey, secretKey string
	)
	if c.String(binanceAPIKeyFlag) == "" {
		return nil, errors.New("cannot create binance client, lack of api key")
	}
	apiKey = c.String(binanceAPIKeyFlag)

	if c.String(binanceSecretKeyFlag) == "" {
		return nil, errors.New("cannot create binance client, lack of secret key")
	}
	secretKey = c.String(binanceSecretKeyFlag)
	rps := c.Float64(binanceRequestPerSecond)
	if rps <= 0 {
		return nil, errors.New("rate limit must be greater than 0")
	}

	limiter := NewRateLimiter(rps)
	return NewBinance(apiKey, secretKey, sugar, WithRateLimiter(limiter)), nil
}
