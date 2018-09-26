package main

import (
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/report/token"
	"github.com/urfave/cli"
)

const (
	nodeURLFlag         = "node"
	nodeURLDefaultValue = "https://mainnet.infura.io"
	outputFlag          = "output"
)

func main() {
	app := cli.NewApp()
	app.Name = "token reserve fetcher"
	app.Usage = "fetching token reserve mapping information"
	app.Action = run
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  nodeURLFlag,
			Usage: "Ethereum node provider URL",
			Value: nodeURLDefaultValue,
		},
		cli.StringFlag{
			Name:  outputFlag,
			Usage: "output file location",
			Value: "./output.json",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	f, err := token.NewFetcher(c.String(nodeURLFlag), c.String(outputFlag))
	if err != nil {
		return err
	}
	return f.Fetch()
}
