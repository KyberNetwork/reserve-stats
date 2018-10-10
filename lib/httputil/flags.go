package httputil

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/urfave/cli"
)

const (
	//PortFlag string for flag port
	PortFlag = "port"
	// httpAddressFlag tells which network address the HTTP server will listen to.
	// Example: 127.0.0.1:8000
	httpAddressFlag = "listen"
)

// NewHTTPCliFlags creates new cli flags for HTTP Server.
func NewHTTPCliFlags(defaultPort HTTPPort) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   httpAddressFlag,
			Usage:  "HTTP server address",
			EnvVar: app.JoinEnvVar(app.JoinEnvVar("HTTP_ADDRESS")), //
			Value:  fmt.Sprintf("127.0.0.1:%d", defaultPort),
		},
	}
}

// NewHTTPAddressFromContext returns the configured address to listen to from cli flags configuration.
func NewHTTPAddressFromContext(c *cli.Context) string {
	return c.String(httpAddressFlag)
}

// NewHostFromContext return a host string for http server to bind.
func NewHostFromContext(c *cli.Context) string {
	portStr := c.String(PortFlag)
	return fmt.Sprintf(":%s", portStr)
}
