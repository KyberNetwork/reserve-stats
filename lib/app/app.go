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

	validDeployments = map[deployment.Deployment]struct{}{
		deployment.Production: {},
		deployment.Staging:    {},
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

func stringToDeploymentMode(mode string) deployment.Deployment {
	switch mode {
	case deployment.Staging.String():
		return deployment.Staging
	case deployment.Production.String():
		return deployment.Production
	}
	return deployment.Production
}

// Validate validates common application configuration flags.
func Validate(c *cli.Context) error {
	mode := c.GlobalString(modeFlag)
	_, ok := validRunningModes[mode]
	if !ok {
		return fmt.Errorf("invalid running mode: %q", c.GlobalString(modeFlag))
	}

	dpl := c.GlobalString(Flag)
	deploymentMode := stringToDeploymentMode(dpl)
	if _, ok = validDeployments[deploymentMode]; !ok {
		return fmt.Errorf("invalid dpl: %q", c.GlobalString(Flag))
	}

	return nil
}
