package httputil

import (
	"fmt"

	"github.com/urfave/cli"
)

const (
	//PortFlag string for flag port
	PortFlag = "port"
)

// NewHTTPFlags returns cli flags to configure a http server.
func NewHTTPFlags(prefix string, defaultPort HTTPPort) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   PortFlag,
			Usage:  "Define http server port",
			Value:  fmt.Sprint(defaultPort),
			EnvVar: prefix + "PORT",
		},
	}
}
