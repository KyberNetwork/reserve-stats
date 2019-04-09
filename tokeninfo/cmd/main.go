package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/tokeninfo"
)

const outputFlag = "output"

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
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func reserve(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flusher, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flusher()

	client, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	internalNetworkClient, err := contracts.NewInternalNetwork(
		contracts.InternalNetworkContractAddress().MustGetOneFromContext(c),
		client,
	)
	if err != nil {
		return err
	}

	f, err := tokeninfo.NewReserveCrawler(
		sugar,
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

	return json.NewEncoder(output).Encode(result)
}
