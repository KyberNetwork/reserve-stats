package broadcast

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	geoURLFlag         = "broadcast-url"
	geoURLDefaultValue = "https://broadcast.kyber.network"

	readAccessKeyFlag = "read-access-key"
	readSecretKeyFlag = "read-secret-key"
)

// NewCliFlags returns cli flags to configure a broadcast client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   geoURLFlag,
			Usage:  "Fetch trade logs to block",
			Value:  geoURLDefaultValue,
			EnvVar: "GEO_URL",
		},
		cli.StringFlag{
			Name:   readAccessKeyFlag,
			Usage:  "key for access GET api",
			EnvVar: "READ_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   readSecretKeyFlag,
			Usage:  "seceret key for GET api",
			EnvVar: "READ_SECRET_KEY",
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

	return NewClient(sugar, geoURL, c.String(readAccessKeyFlag), c.String(readSecretKeyFlag))
}
