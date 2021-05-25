package coreclient

import (
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	coreEndpointFlag        = "core-endpoint"
	coreClientAPIKeyFlag    = "core-api-key"
	coreClientAPISecretFlag = "core-secret-key"
)

// NewCoreClientFlags return core client flags
func NewCoreClientFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   coreEndpointFlag,
			Usage:  "endpoint for core",
			EnvVar: "CORE_ENDPOINT",
		},
		cli.StringFlag{
			Name:   coreClientAPIKeyFlag,
			Usage:  "api key for core",
			EnvVar: "CORE_API_KEY",
		},
		cli.StringFlag{
			Name:   coreClientAPISecretFlag,
			Usage:  "secret key for core",
			EnvVar: "CORE_SECRET_KEY",
		},
	}
}

// NewCoreClientFromContext return new core client object
func NewCoreClientFromContext(c *cli.Context, sugar *zap.SugaredLogger) *CoreClient {
	endpoint := c.String(coreEndpointFlag)
	if endpoint == "" {
		sugar.Error("core endpoint is not provided")
	}
	apiKey := c.String(coreClientAPIKeyFlag)
	if apiKey == "" {
		sugar.Error("core api key is not provided")
	}
	secretKey := c.String(coreClientAPISecretFlag)
	if secretKey == "" {
		sugar.Error("core secret key is not provided")
	}
	return NewCoreClient(endpoint, apiKey, secretKey, sugar)
}
