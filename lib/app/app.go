package app

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"
)

const (
	modeFlag = "mode"

	developmentMode = "develop"
	productionMode  = "production"
)

var (
	validRunningModes = map[string]struct{}{
		developmentMode: {},
		productionMode:  {},
	}
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
		cli.StringFlag{
			Name:  deployment.Flag,
			Usage: "Kyber Network deployment name",
			Value: productionMode,
		},
	}
	return app
}

// Validate validates common application configuration flags.
func Validate(c *cli.Context) error {
	mode := c.GlobalString(modeFlag)
	_, ok := validRunningModes[mode]
	if !ok {
		return fmt.Errorf("invalid running mode: %q", c.GlobalString(modeFlag))
	}
	return nil
}
