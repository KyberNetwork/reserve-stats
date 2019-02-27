package binance

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	binanceAPIKeyFlag    = "binance-api-key"
	binanceSecretKeyFlag = "binance-secret-key"
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
	}
}

//NewBinanceClientFromContext return binance client
func NewBinanceClientFromContext(c *cli.Context) (*Client, error) {
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
