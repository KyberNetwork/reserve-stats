package broadcast

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	broadcastURLFlag = "broadcast-url"

	broadcastAccessKeyIDFlag     = "broadcast-access-key-id"
	broadcastSecretAccessKeyFlag = "broadcast-secret-access-key"
)

// NewCliFlags returns cli flags to configure a broadcast client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   broadcastURLFlag,
			Usage:  "URL of broadcast service",
			EnvVar: "BROADCAST_URL",
		},
		cli.StringFlag{
			Name:   broadcastAccessKeyIDFlag,
			Usage:  "access key id of broadcast service",
			EnvVar: "BROADCAST_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   broadcastSecretAccessKeyFlag,
			Usage:  "secret access key of broadcast service",
			EnvVar: "BROADCAST_SECRET_ACCESS_KEY",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (Interface, error) {
	var (
		url             = c.String(broadcastURLFlag)
		accessKeyID     = c.String(broadcastAccessKeyIDFlag)
		secretAccessKey = c.String(broadcastSecretAccessKeyFlag)

		options []ClientOption
	)

	if len(url) == 0 {
		sugar.Warnw("no broadcast service configured, using noop service")
		return NewNoop(), nil
	}

	err := validation.Validate(url,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid broadcast url: %q, error: %s", url, err)
	}

	if len(accessKeyID) != 0 && len(broadcastAccessKeyIDFlag) != 0 {
		options = append(options, WithAuth(accessKeyID, secretAccessKey))
	}

	return NewClient(sugar, url, options...), nil
}
