package app

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/urfave/cli"
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
			Name:  Flag,
			Usage: "Kyber Network deployment name",
			Value: productionMode,
		},
	}
	return app
}

func stringToDeploymentMode(mode string) (deployment.Deployment, error) {
	switch mode {
	case deployment.Staging.String():
		return deployment.Staging, nil
	case deployment.Production.String():
		return deployment.Production, nil
	}
	return 0, fmt.Errorf("deployment mode is not valid: %s", mode)
}

// Validate validates common application configuration flags.
func Validate(c *cli.Context) error {
	mode := c.GlobalString(modeFlag)
	_, ok := validRunningModes[mode]
	if !ok {
		return fmt.Errorf("invalid running mode: %q", c.GlobalString(modeFlag))
	}

	dpl := c.GlobalString(Flag)
	_, err := stringToDeploymentMode(dpl)
	if err != nil {
		return err
	}
	return nil
}
