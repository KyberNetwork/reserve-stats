package userprofile

import (
	"fmt"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	userprofileURLFlag  = "user-url"
	maxUserCacheFlag    = "max-user-cache"
	maxUserCacheDefault = 1000
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
	return NewClient(sugar, userURL)
}

// NewCachedClientFromContext return new cached client from cli flags
func NewCachedClientFromContext(client *Client, c *cli.Context) *CachedClient {
	maxCacheSize := c.Int64(maxUserCacheFlag)
	return NewCachedClient(client, maxCacheSize)
}
