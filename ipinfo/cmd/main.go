package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/ipinfo"
	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/urfave/cli"
)

const (
	dataDirFlag = "data-dir"
)

func main() {
	app := app.NewApp()
	app.Name = "ip locator checker"
	app.Usage = "get countery of given IP address"
	app.Version = "0.0.1"

	app.Action = iplocatorServer

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   dataDirFlag,
			Usage:  "directory to store the GeoLite2-Country.mmdb file",
			Value:  ".",
			EnvVar: "IP_INFO_DATA_DIR",
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func iplocatorServer(c *cli.Context) error {
	logger, err := app.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()
	server, err := ipinfo.NewHTTPServer(sugar, c.String(dataDirFlag))
	if err != nil {
		return err
	}
	return server.Run()
}
