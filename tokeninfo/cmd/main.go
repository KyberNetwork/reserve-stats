package main

import (
	"encoding/json"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"log"
	"math/big"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/tokeninfo"
	"github.com/urfave/cli"
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

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	client, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}

	proxyAddress := contracts.ProxyContractAddress().MustGetOneFromContext(c)
	internalNetworkAddress, err := contracts.InternalNetworkContractAddress(proxyAddress, client, big.NewInt(1))
	if err != nil {
		return err
	}

	internalNetworkClient, err := contracts.NewInternalNetwork(
		internalNetworkAddress,
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
