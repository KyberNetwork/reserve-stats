package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
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
		cli.Uint64Flag{
			Name:   fromFlag,
			Usage:  "From timestamp(millisecond) to get trade history from",
			EnvVar: "FROM",
		},
		cli.Uint64Flag{
			Name:   toFlag,
			Usage:  "To timestamp(millisecond) to get trade history to",
			EnvVar: "TO",
		},
	)

	app.Flags = append(app.Flags, binance.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getTradeHistoryForOneSymBol(binanceClient *binance.Client, fromTime, toTime time.Time, symbol string,
	tradeHistories *sync.Map, wg *sync.WaitGroup) error {
	var (
		logger = binanceClient.Sugar.With("func", "accounting/cmd/accounting-binance-trade-fetcher")
	)
	result := []binance.TradeHistory{}
	startTime := fromTime
	endTime := toTime
	defer wg.Done()
	for {
		tradeHistoriesResponse, err := binanceClient.GetTradeHistory(symbol, -1, startTime, endTime)
		if err != nil {
			return err
		}
		// while result != empty, get trades latest time to toTime
		if len(tradeHistoriesResponse) == 0 {
			break
		}
		logger.Debugw("trade history for", "symbol", symbol, "history", tradeHistoriesResponse)
		result = append(result, tradeHistoriesResponse...)
		lastTrade := tradeHistoriesResponse[len(tradeHistoriesResponse)-1]
		startTime = timeutil.TimestampMsToTime(lastTrade.Time + 1)
	}
	if len(result) != 0 {
		tradeHistories.Store(symbol, result)
	}
	return nil
}

func getTradeHistory(binanceClient *binance.Client, fromTime, toTime time.Time) error {
	var (
		tradeHistories sync.Map
		logger         = binanceClient.Sugar.With("func", "accounting/cmd/accounting-binance-trade-fetcher")
	)
	// get list token
	exchangeInfo, err := binanceClient.GetExchangeInfo()
	if err != nil {
		return err
	}
	tokenPairs := exchangeInfo.Symbols
	// get fromTime to toTime
	wg := sync.WaitGroup{}
	for _, pairs := range tokenPairs {
		wg.Add(1)
		go getTradeHistoryForOneSymBol(binanceClient, fromTime, toTime, pairs.Symbol, &tradeHistories, &wg)
	}
	wg.Wait()
	// log here for test get trade history without persistence storage
	tradeHistories.Range(func(key, value interface{}) bool {
		logger.Info("symbol", key, "history", value)
		return true
	})
	// TODO: save to storage
	return nil
}

func run(c *cli.Context) error {
	var (
		fromTime, toTime time.Time
	)

	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}

	defer logger.Sync()

	sugar := logger.Sugar()
	sugar.Info("initiate fetcher")

	binanceClient, err := binance.NewClientFromContext(c, sugar)
	if err != nil {
		return err
	}

	if c.Uint64(fromFlag) != 0 {
		fromTime = timeutil.TimestampMsToTime(c.Uint64(fromFlag))
	}

	if c.Uint64(toFlag) != 0 {
		toTime = timeutil.TimestampMsToTime(c.Uint64(toFlag))
	}

	err = getTradeHistory(binanceClient, fromTime, toTime)
	if err != nil {
		return err
	}

	return nil
}
