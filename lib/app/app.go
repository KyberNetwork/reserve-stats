package app

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	modeFlag         = "mode"
	developmentMode  = "development"
	productionMode   = "production"
	ethereumNodeFlag = "ethereum-node"
)

// NewApp creates a new cli App instance with common flags pre-loaded.
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  modeFlag,
			Usage: "app running mode",
			Value: developmentMode,
		},
	}
	return app
}

// NewEthereumNodeFlags returns cli flag for ethereum node url input
func NewEthereumNodeFlags(prefix string) cli.Flag {
	return cli.StringFlag{
		Name:   ethereumNodeFlag,
		Usage:  "Ethereum Node URL",
		EnvVar: prefix + "ETHEREUM_NODE",
	}
}

// NewEthereumClientFromFlag returns Ethereum client from flag variable, or error if occurs
func NewEthereumClientFromFlag(c *cli.Context) (*ethclient.Client, error) {
	ethereumNodeURL := c.GlobalString(ethereumNodeFlag)
	return ethclient.Dial(ethereumNodeURL)
}

// NewLogger creates a new logger instance.
// The type of logger instance will be different with different application running modes.
func NewLogger(c *cli.Context) (*zap.Logger, error) {
	mode := c.GlobalString(modeFlag)
	switch mode {
	case productionMode:
		return zap.NewProduction()
	default:
		return zap.NewDevelopment()
	}
}
