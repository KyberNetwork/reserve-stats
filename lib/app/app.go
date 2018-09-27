package app

import (
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	modeFlag        = "mode"
	developmentMode = "development"
	productionMode  = "production"
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
