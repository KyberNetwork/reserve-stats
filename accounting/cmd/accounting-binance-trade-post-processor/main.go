package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	fetcher "github.com/KyberNetwork/reserve-stats/accounting/binance/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/binance/storage/tradestorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/marketdata"
)

const (
	retryDelayFlag    = "retry-delay"
	attemptFlag       = "attempt"
	defaultRetryDelay = 2 * time.Minute
	defaultAttempt    = 4

	marketDataBaseURL = "https://staging-market-data.knstats.com"
)

var sugar *zap.SugaredLogger

func main() {
	app := libapp.NewApp()
	app.Name = "Accounting binance trades fetcher"
	app.Usage = "Fetch and store trades history from binance"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.DurationFlag{
			Name:   retryDelayFlag,
			Usage:  "delay time when do a retry",
			EnvVar: "RETRY_DELAY",
			Value:  defaultRetryDelay,
		},
		cli.IntFlag{
			Name:   attemptFlag,
			Usage:  "number of time doing retry",
			EnvVar: "ATTEMPT",
			Value:  defaultAttempt,
		},
	)

	app.Flags = append(app.Flags, binance.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexTradesDB)...)

	if err := app.Run(os.Args); err != nil {
		sugar.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		flusher func()
		err     error
	)
	sugar, flusher, err = libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}

	defer flusher()

	sugar.Info("initiate fetcher")

	storage, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	binanceStorage, err := tradestorage.NewDB(sugar, storage)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := binanceStorage.Close(); cErr != nil {
			sugar.Errorf("Close database error", "error", cErr)
		}
	}()

	binanceClient, err := binance.NewBinance("", "", sugar) // this is public client to get exchange info
	if err != nil {
		return err
	}

	marketDataClient := marketdata.NewMarketDataClient(marketDataBaseURL, sugar)

	var (
		tokenPairs  []binance.Symbol
		quotes      = make(map[string]string)
		quoteString []string
	)
	exchangeInfo, err := binanceClient.GetExchangeInfo()
	if err != nil {
		return err
	}
	tokenPairs = exchangeInfo.Symbols
	for _, pair := range tokenPairs {
		if _, exist := quotes[pair.QuoteAsset]; !exist {
			quotes[pair.QuoteAsset] = pair.QuoteAsset
			quoteString = append(quoteString, pair.QuoteAsset)
		}
	}
	regexpString := fmt.Sprintf(".*(%s)$", strings.Join(quoteString, "|"))
	re := regexp.MustCompile(regexpString)

	sugar.Infow("quotes", "list", quotes)

	retryDelay := c.Duration(retryDelayFlag)
	attempt := c.Int(attemptFlag)
	batchSize := 5 // dummy batch size to init fetcher

	binanceFetcher := fetcher.NewFetcher(sugar, binanceClient, retryDelay, attempt, batchSize, binanceStorage, "", marketDataClient, nil)
	notEthTrades, err := binanceStorage.GetNotETHTrades()
	if err != nil {
		return err
	}
	for originalSymbol, trades := range notEthTrades {
		quote := quoteFromOriginalSymbol(re, originalSymbol)
		symbol := "ETH" + quote
		if err := binanceFetcher.UpdateTradeNotETH(originalSymbol, symbol, trades); err != nil {
			return err
		}
	}
	return nil
}

func quoteFromOriginalSymbol(re *regexp.Regexp, symbol string) string {
	res := re.FindAllStringSubmatch(symbol, -1)
	return res[0][1]
}
