package main

import (
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/http"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

const (
	addrFlag = "addr"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs HTTP Api"
	app.Usage = "Serve trade logs data"
	app.Version = "0.0.1"
	app.Action = func(c *cli.Context) error {
		logger, err := libapp.NewLogger(c)
		if err != nil {
			return err
		}
		defer logger.Sync()

		sugar := logger.Sugar()

		coreClient, err := core.NewClientFromContext(sugar, c)
		if err != nil {
			return err
		}

		influxClient, err := influxdb.NewClientFromContext(c)
		if err != nil {
			return err
		}

		serverAddr := c.String(addrFlag)

		influxStorage, err := storage.NewInfluxStorage(
			sugar,
			"trade_logs",
			influxClient,
			core.NewCachedClient(coreClient),
		)
		if err != nil {
			return err
		}

		api := http.NewServer(influxStorage, serverAddr)
		api.Start()

		return nil
	}

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   addrFlag,
			Usage:  "Trade logs server address",
			EnvVar: "TRADE_LOGS_SERVER_ADDR",
		},
	)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, core.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
