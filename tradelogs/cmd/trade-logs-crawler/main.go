package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	nodeURLFlag         = "node"
	nodeURLDefaultValue = "https://mainnet.infura.io"
	fromBlockFlag       = "from-block"
	toBlockFlag         = "to-block"

	envVarPrefix = "TRADE_LOGS_CRAWLER_"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs Fetcher"
	app.Usage = "Fetch trade logs on KyberNetwork"
	app.Version = "0.0.1"
	app.Action = getTradeLogs

	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   "node",
			Usage:  "Ethereum node provider URL",
			Value:  nodeURLDefaultValue,
			EnvVar: envVarPrefix + "NODE",
		},
		cli.StringFlag{
			Name:   "from-block",
			Usage:  "Fetch trade logs from block",
			EnvVar: envVarPrefix + "FROM_BLOCK",
		},
		cli.StringFlag{
			Name:   "to-block",
			Usage:  "Fetch trade logs to block",
			EnvVar: envVarPrefix + "TO_BLOCK",
		},
	)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseBigIntFlag(c *cli.Context, flag string) (*big.Int, error) {
	val := c.String(flag)
	if err := validation.Validate(val, validation.Required, is.Digit); err != nil {
		return nil, err
	}

	result, ok := big.NewInt(0).SetString(val, 0)
	if !ok {
		return nil, fmt.Errorf("invalid number %s", c.String(flag))
	}
	return result, nil
}

func getTradeLogs(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	fromBlock, err := parseBigIntFlag(c, fromBlockFlag)
	if err != nil {
		return fmt.Errorf("invalid from block: %q, error: %s", c.String(fromBlockFlag), err)
	}

	toBlock, err := parseBigIntFlag(c, toBlockFlag)
	if err != nil {
		return fmt.Errorf("invalid to block: %q, error: %s", c.String(toBlockFlag), err)
	}

	nodeURL := c.String(nodeURLFlag)
	if err = validation.Validate(nodeURL, validation.Required, is.URL); err != nil {
		return fmt.Errorf("invalid node url: %q, error: %s", nodeURL, err)
	}

	crawler, err := tradelogs.NewTradeLogCrawler(
		sugar,
		nodeURL,
		tradelogs.NewCMCEthUSDRate(),
	)
	if err != nil {
		return err
	}

	tradeLogs, err := crawler.GetTradeLogs(fromBlock, toBlock, time.Second*5)
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(tradeLogs)
}
