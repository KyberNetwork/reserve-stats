package client

import (
	"fmt"

	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	addressServerFlag = "address-server-url"
)

//NewClientFlags return flags to init addresses client
func NewClientFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   addressServerFlag,
			Usage:  "The address of Reserve Addresses server with port number",
			EnvVar: "ADDRESS_SERVER_URL",
		}}
}

//NewClientFromContext return a client from address flag input
func NewClientFromContext(c *cli.Context, sugar *zap.SugaredLogger) (*Client, error) {
	clientAddress := c.String(addressServerFlag)
	if clientAddress == "" {
		return nil, fmt.Errorf("no client address is set. Set it via %s flag", addressServerFlag)
	}
	return NewClient(sugar, clientAddress)
}
