package main

import (
	"log"
	"os"

	appNames "github.com/KyberNetwork/reserve-stats/integration-app-names"
	"github.com/KyberNetwork/reserve-stats/integration-app-names/http"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/urfave/cli"
)

const (
	dataFilePathFlag = "data-path"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Addre to Intergration App Name"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.AddrToAppName)...)
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   dataFilePathFlag,
			Usage:  "file path to address to app name json",
			EnvVar: "APPNAME_DATA_PATH",
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
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}

	defer logger.Sync()
	sugar := logger.Sugar()

	addrToAppname := appNames.NewMapAddrAppName()
	if err := addrToAppname.LoadFromFile(c.String(dataFilePathFlag)); err != nil {
		return err
	}
	server, err := http.NewServer(httputil.NewHTTPAddressFromContext(c), addrToAppname, sugar)
	if err != nil {
		return err
	}

	sugar.Info("Run Addr to Appname module")
	return server.Run()

}
