package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/lib/application"
	"github.com/KyberNetwork/reserve-stats/token"
	"github.com/urfave/cli"
)

const (
	nodeURLFlag         = "node"
	nodeURLDefaultValue = "https://mainnet.infura.io"
	outputFlag          = "output"
)

func main() {
	app := application.NewApp()
	app.Name = "token reserve fetcher"
	app.Usage = "fetching token reserve mapping information"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "reserve",
			Aliases: []string{"r"},
			Usage:   "report which reserves provides which token",
			Action:  reserve,
		},
	}

	app.Flags = append(app.Flags,
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
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func reserve(c *cli.Context) error {
	logger, err := application.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	f, err := token.NewReserveCrawler(
		sugar,
		c.GlobalString(nodeURLFlag))
	if err != nil {
		return err
	}

	output, err := os.Create(c.GlobalString(outputFlag))
	if err != nil {
		return err
	}

	result, err := f.Fetch()
	if err != nil {
		return err
	}

	return json.NewDecoder(output).Decode(result)
}
