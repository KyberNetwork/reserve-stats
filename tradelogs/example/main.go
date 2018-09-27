package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
)

const (
	nodeURLFlag         = "node"
	nodeURLDefaultValue = "https://mainnet.infura.io"
	fromBlockFlag       = "from-block"
	toBlockFlag         = "to-block"
)

func main() {
	var fromBlock, toBlock string

	app := app.NewApp()
	app.Name = "Trade Logs Fetcher"
	app.Usage = "Fetch trade logs on KyberNetwork"
	app.Version = "0.0.1"
	app.Action = getTradeLogs

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:  "node",
			Usage: "Ethereum node provider URL",
			Value: nodeURLDefaultValue,
		},
		cli.StringFlag{
			Name:        "from-block",
			Usage:       "Fetch trade logs from block",
			Destination: &fromBlock,
		},
		cli.StringFlag{
			Name:        "to-block",
			Usage:       "Fetch trade logs to block",
			Destination: &toBlock,
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getTradeLogs(c *cli.Context) error {
	logger, err := app.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	fromBlock, err := strconv.ParseUint(c.GlobalString(fromBlockFlag), 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid from-block value: %s (should be integer)", c.GlobalString(fromBlockFlag))
	}

	toBlock, err := strconv.ParseUint(c.GlobalString(toBlockFlag), 10, 64)
	if err != nil {
		return fmt.Errorf("Invalid to-block value: %s (should be integer)", c.GlobalString(fromBlockFlag))
	}

	crawler, err := tradelogs.NewTradeLogCrawler(
		c.GlobalString(nodeURLFlag),
		tradelogs.NewCMCEthUSDRate(),
	)
	if err != nil {
		return err
	}

	tradeLogs, err := crawler.GetTradeLogs(fromBlock, toBlock)
	if err == nil {
		for _, logItem := range tradeLogs {
			fmt.Printf("%+v\n", logItem)
		}
	}

	return nil
}
