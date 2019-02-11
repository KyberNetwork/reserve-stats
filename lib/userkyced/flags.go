package userkyced

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	userKycedURL            = "user-kyc-url"
	userKycedSigningKeyFlag = "user-kyc-signing-key"
	userKycedKeyIDFlag      = "user-kyc-key-id"
)

// NewCliFlags returns cli flags to configure a user kyc status client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   userKycedURL,
			Usage:  "user kyced API URL. If this doesn't support, fallback to default PosGres DB for kyced checking",
			EnvVar: "USER_KYCED_URL",
			Value:  "",
		},
		cli.StringFlag{
			Name:   userKycedSigningKeyFlag,
			Usage:  "user profile Signing Key",
			EnvVar: "USER_KYCED_SIGNING_KEY",
			Value:  "",
		},
		cli.StringFlag{
			Name:   userKycedKeyIDFlag,
			Usage:  "user profile Signing Key ID",
			EnvVar: "USER_KYCED_KEY_ID",
			Value:  "",
		},
	}
}

// NewClientFromContext returns new user kyc client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	userURL := c.String(userKycedURL)
	if userURL == "" {
		return nil, nil
	}
	err := validation.Validate(userURL,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("user kyced url: %s", err.Error())
	}
	signingKey := c.String(userKycedSigningKeyFlag)

	keyID := c.String(userKycedKeyIDFlag)

	return NewClient(sugar, userURL, signingKey, keyID)
}
