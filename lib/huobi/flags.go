package huobi

import (
	"fmt"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	huobiAPIKeyFlag    = "huobi-api-key"
	huobiSecretKeyFlag = "huobi-secret-key"
)

//NewCliFlags return cli flags to configure cex client
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
	}
}

//NewHuobiClientFromContext return huobi client
func NewHuobiClientFromContext(c *cli.Context, sugar *zap.SugaredLogger) (*Client, error) {
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

	return NewHuobiClient(apiKey, secretKey, sugar), nil
}
