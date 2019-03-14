package main

import (
	"log"
	"os"

	appNames "github.com/KyberNetwork/reserve-stats/app-names"
	"github.com/KyberNetwork/reserve-stats/app-names/http"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/urfave/cli"
)

const (
	dataFilePathFlag = "data-path"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Address to Intergration App Name"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AppName)...)
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   dataFilePathFlag,
			Usage:  "file path to address to app name json",
			EnvVar: "DATA_PATH",
		},
	)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	addrToAppname := appNames.NewMapAddrAppName(appNames.WithDataFile(c.String(dataFilePathFlag)))

	server, err := http.NewServer(httputil.NewHTTPAddressFromContext(c), addrToAppname, sugar)
	if err != nil {
		return err
	}

	sugar.Info("Run Addr to Appname module")
	return server.Run()

}
