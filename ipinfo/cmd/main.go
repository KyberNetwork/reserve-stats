package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/ipinfo"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"
)

const dataDirFlag = "data-dir"

func main() {
	app := libapp.NewApp()
	app.Name = "ip locator checker"
	app.Usage = "get countery of given IP address"
	app.Version = "0.0.1"

	app.Action = runHTTPServer

	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.IPLocatorPort)...)

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   dataDirFlag,
			Usage:  "directory to store the GeoLite2-Country.mmdb file",
			Value:  ".",
			EnvVar: "DATA_DIR",
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runHTTPServer(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	err = validation.Validate(
		c.String(httputil.PortFlag),
		is.Int,
	)
	if err != nil {
		sugar.Errorw("Get error while validate --port flag", "err", err.Error())
		return err
	}

	server, err := ipinfo.NewHTTPServer(sugar, c.String(dataDirFlag), httputil.NewHTTPAddressFromContext(c))
	if err != nil {
		return err
	}
	return server.Run()
}
