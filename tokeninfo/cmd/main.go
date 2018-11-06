package main

import (
	"encoding/json"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/tokeninfo"
	"github.com/urfave/cli"
)

const (
	nodeURLFlag         = "node"
	nodeURLDefaultValue = "https://mainnet.infura.io"
	outputFlag          = "output"
)

func main() {
	app := libapp.NewApp()
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
			Name:  outputFlag,
			Usage: "output file location",
			Value: "./output.json",
		},
	)
	app.Flags = append(app.Flags, libapp.NewEthereumNodeFlags())

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func reserve(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	client, err := libapp.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	internalNetworkClient, err := contracts.NewInternalNetwork(
		contracts.InternalNetworkContractAddress().MustGetFromContext(c),
		client,
	)
	if err != nil {
		return err
	}

	f, err := tokeninfo.NewReserveCrawler(
		logger.Sugar(),
		internalNetworkClient,
	)
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
