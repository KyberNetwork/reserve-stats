package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/tokenprice/server"
	"github.com/KyberNetwork/reserve-stats/tokenprice/storage"
)

const (
	hostFlag = "host"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Token rate crawler"
	app.Usage = "Crawl token rate from some EX"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   hostFlag,
			Usage:  "provide host for token price api",
			EnvVar: "HOST",
		},
	)

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(storage.DefaultDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()
	s, err := storage.NewStorageFromContext(sugar, c)
	if err != nil {
		sugar.Errorw("failed to init storage", "error", err)
		return err
	}
	sv := server.NewServer(sugar, c.String(hostFlag), s)
	return sv.Start()
}
