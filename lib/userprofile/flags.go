package userprofile

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	userprofileURLFlag        = "user-url"
	userprofileSigningKeyFlag = "user-signing-key"
	maxUserCacheFlag          = "max-user-cache"
	maxUserCacheDefault       = 1000
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   userprofileURLFlag,
			Usage:  "user profile API URL",
			EnvVar: "USER_URL",
		},
		cli.Int64Flag{
			Name:   maxUserCacheFlag,
			Usage:  "max Cache Size for user profile client. Default to 1000",
			EnvVar: "MAX_USER_CACHE",
			Value:  maxUserCacheDefault,
		},
		cli.StringFlag{
			Name:   userprofileSigningKeyFlag,
			Usage:  "user profile Signing Key",
			EnvVar: "USER_SIGNING_KEY",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	userURL := c.String(userprofileURLFlag)
	err := validation.Validate(userURL,
		validation.Required,
		is.URL,
	)
	if err != nil {
		return nil, fmt.Errorf("user profile url: %s", err.Error())
	}
	signingKey := c.String(userprofileSigningKeyFlag)
	err = validation.Validate(signingKey,
		validation.Required,
	)
	if err != nil {
		return nil, fmt.Errorf("core signing key: %s", err.Error())
	}
	return NewClient(sugar, userURL, signingKey)
}

// NewCachedClientFromContext return new cached client from cli flags
func NewCachedClientFromContext(client *Client, c *cli.Context) *CachedClient {
	maxCacheSize := c.Int64(maxUserCacheFlag)
	return NewCachedClient(client, maxCacheSize)
}
