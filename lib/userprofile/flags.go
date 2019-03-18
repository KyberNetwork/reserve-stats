package userprofile

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	userprofileURLFlag        = "user-profile-url"
	userprofileSigningKeyFlag = "user-profile-signing-key"
	maxUserCacheFlag          = "max-user-profile-cache"
	maxUserCacheDefault       = 1000
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   userprofileURLFlag,
			Usage:  "user profile API URL",
			EnvVar: "USER_PROFILE_URL",
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
			EnvVar: "USER_PROFILE_SIGNING_KEY",
		},
	}
}

// NewClientFromContext returns new core client from cli flags.
func NewClientFromContext(sugar *zap.SugaredLogger, c *cli.Context) (*Client, error) {
	userURL := c.String(userprofileURLFlag)
	if userURL == "" {
		return nil, nil
	}
	err := validation.Validate(userURL,
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
		return nil, fmt.Errorf("user signing key: %s", err.Error())
	}
	return NewClient(sugar, userURL, signingKey)
}

// NewCachedClientFromContext return new cached client from cli flags
func NewCachedClientFromContext(client *Client, c *cli.Context) *CachedClient {
	if client == nil {
		return nil
	}
	maxCacheSize := c.Int64(maxUserCacheFlag)
	return NewCachedClient(client, maxCacheSize)
}
