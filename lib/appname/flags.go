package appname

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	appNameURLFlag = "app-name-url"
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   appNameURLFlag,
			Usage:  "Core API URL",
			EnvVar: "APP_NAME_URL",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	appNameURL := c.String(appNameURLFlag)
	err := validation.Validate(appNameURL,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("appName URL: %s", err.Error())
	}

	return NewClient(sugar, appNameURL)
}
