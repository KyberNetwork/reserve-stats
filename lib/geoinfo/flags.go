package geoinfo

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	geoURLFlag         = "geo-url"
	geoURLDefaultValue = "https://broadcast.kyber.network"
)

// NewCliFlags returns cli flags to configure a geoinfo client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   geoURLFlag,
			Usage:  "Fetch trade logs to block",
			Value:  geoURLDefaultValue,
			EnvVar: "GEO_URL",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	geoURL := c.String(geoURLFlag)
	err := validation.Validate(geoURL,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid geo url: %q, error: %s", geoURL, err)
	}

	return NewClient(sugar, geoURL)
}
