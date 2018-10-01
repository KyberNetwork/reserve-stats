package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/ip"
	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/urfave/cli"
)

const (
	ipFlag         = "ip"
	ipDefaultValue = "8.8.8.8"
)

func main() {
	app := app.NewApp()
	app.Name = "ip locator checker"
	app.Usage = "get countery of given IP address"
	app.Version = "0.0.1"

	app.Action = iplocator

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:  ipFlag,
			Usage: "IP want to check",
			Value: ipDefaultValue,
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func iplocator(c *cli.Context) error {
	logger, err := app.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	f, err := ip.NewLocator(
		sugar,
	)
	if err != nil {
		return err
	}

	result, err := f.IPToCountry(c.String(ipFlag))
	if err != nil {
		return err
	}

	sugar.Infow(result)

	return nil
}
