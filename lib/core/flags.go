package core

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	coreURLFlag        = "core-url"
	coreSigningKeyFlag = "core-signing-key"
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags() ([]cli.Flag) {
	return []cli.Flag{
		cli.StringFlag{
			Name:   coreURLFlag,
			Usage:  "Core API URL",
			EnvVar: "CORE_URL",
		},
		cli.StringFlag{
			Name:   coreSigningKeyFlag,
			Usage:  "Core API Signing Key",
			EnvVar: "CORE_SIGNING_KEY",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	coreURL := c.String(coreURLFlag)
	err := validation.Validate(coreURL,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("core url: %s", err.Error())
	}

	coreSigningKey := c.String(coreSigningKeyFlag)
	err = validation.Validate(coreSigningKeyFlag,
		validation.Required,
	)
	if err != nil {
		return nil, fmt.Errorf("core signing key: %s", err.Error())
	}

	return NewClient(sugar, coreURL, coreSigningKey)
}
