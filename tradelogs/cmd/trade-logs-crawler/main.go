package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/tokenrate"
	"github.com/KyberNetwork/reserve-stats/tradelogs"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	"github.com/KyberNetwork/tokenrate/coingecko"
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
			Name:   nodeURLFlag,
			Usage:  "Ethereum node provider URL",
			Value:  nodeURLDefaultValue,
			EnvVar: envVarPrefix + "NODE",
		},
		cli.StringFlag{
			Name:   fromBlockFlag,
			Usage:  "Fetch trade logs from block",
			EnvVar: envVarPrefix + "FROM_BLOCK",
		},
		cli.StringFlag{
			Name:   toBlockFlag,
			Usage:  "Fetch trade logs to block",
			EnvVar: envVarPrefix + "TO_BLOCK",
		},
	)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, core.NewCliFlags()...)
	app.Flags = append(app.Flags, broadcast.NewCliFlags()...)

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
	const dbName = "trade_logs"
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	coreClient, err := core.NewClientFromContext(sugar, c)
	if err != nil {
		return err
	}

	// Crawl trade logs from blockchain
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

	geoClient, err := broadcast.NewClientFromContext(sugar, c)
	if err != nil {
		return err
	}

	crawler, err := tradelogs.NewTradeLogCrawler(
		sugar,
		nodeURL,
		geoClient,
	)
	if err != nil {
		return err
	}

	tradeLogs, err := crawler.GetTradeLogs(fromBlock, toBlock, time.Second*5)
	if err != nil {
		return err
	}

	// Store trade logs into influx DB
	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	influxStorage, err := storage.NewInfluxStorage(
		sugar,
		dbName,
		influxClient,
		core.NewCachedClient(coreClient),
	)
	if err != nil {
		return err
	}

	err = influxStorage.SaveTradeLogs(tradeLogs)
	if err != nil {
		return err
	}

	// fetch eth usd rate
	ethUSDRateFetcher, err := tokenrate.NewETHUSDRateFetcher(sugar, dbName, influxClient, coingecko.New())
	if err != nil {
		return err
	}

	for _, tradelog := range tradeLogs {
		if _, err := ethUSDRateFetcher.FetchRates(tradelog.BlockNumber, tradelog.Timestamp); err != nil {
			return err
		}
	}

	return json.NewEncoder(os.Stdout).Encode(tradeLogs)
}
