package main

import (
	"fmt"
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/ipinfo"
	"github.com/KyberNetwork/reserve-stats/lib/app"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/urfave/cli"
)

const (
	ipFlag      = "ip"
	dataDirFlag = "data-dir"
)

func main() {
	app := app.NewApp()
	app.Name = "ip locator checker"
	app.Usage = "get countery of given IP address"
	app.Version = "0.0.1"

	app.Action = locateIP

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:  ipFlag,
			Usage: "IP want to check",
		},
		cli.StringFlag{
			Name:  dataDirFlag,
			Usage: "directory to store the GeoLite2-Country.mmdb file",
			Value: ".",
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func locateIP(c *cli.Context) error {
	err := validation.Validate(
		c.String(ipFlag),
		validation.Required,
	)
	if err != nil {
		return fmt.Errorf("--ip flags is required")
	}

	logger, err := app.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	f, err := ipinfo.NewLocator(sugar, c.String(dataDirFlag))
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
