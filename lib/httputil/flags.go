package httputil

import (
	"fmt"

	"github.com/urfave/cli"
)

const (
	//PortFlag string for flag port
	PortFlag = "port"
	// httpAddressFlag tells which network address the HTTP server will listen to.
	// Example: 127.0.0.1:8000
	httpAddressFlag = "listen"

	//HTTPKeyFlag is secret key for signing a request
	HTTPKeyFlag = "secret-key"
)

// NewHTTPCliFlags creates new cli flags for HTTP Server.
func NewHTTPCliFlags(defaultPort HTTPPort) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   httpAddressFlag,
			Usage:  "HTTP server address",
			EnvVar: "HTTP_ADDRESS",
			Value:  fmt.Sprintf("127.0.0.1:%d", defaultPort),
		},
	}
}

// NewHTTPAddressFromContext returns the configured address to listen to from cli flags configuration.
func NewHTTPAddressFromContext(c *cli.Context) string {
	return c.String(httpAddressFlag)
}
