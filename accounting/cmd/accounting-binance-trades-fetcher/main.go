package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
)

const (
	fromFlag = "from"
	toFlag   = "to"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting binance trades fetcher"
	app.Usage = "Fetch and store trades history from binance"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.Int64Flag{
			Name:   fromFlag,
			Usage:  "From timestamp(millisecond) to get trade history from",
			EnvVar: "FROM",
		},
		cli.Int64Flag{
			Name:   toFlag,
			Usage:  "To timestamp(millisecond) to get trade history to",
			EnvVar: "TO",
		})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	return nil
}
