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
	bindAddressFlag = "bindAddress"

	defaultBindAddress = "127.0.0.1:8000"
)

func main() {
	app := libapp.NewAppWithMode()
	app.Name = "Token price API"
	app.Usage = "Serve api for token price"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   bindAddressFlag,
			Usage:  "Address to serve ETH price endpoint",
			Value:  defaultBindAddress,
			EnvVar: "BIND_ADDRESS",
		},
	)

	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(storage.DefaultDB)...)
	app.Flags = append(app.Flags, libapp.NewSentryFlags()...)
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
	sv := server.NewServer(sugar, c.String(bindAddressFlag), s)
	return sv.Start()
}
