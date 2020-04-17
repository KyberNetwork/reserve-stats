package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	huobiFetcher "github.com/KyberNetwork/reserve-stats/accounting/huobi/fetcher"
	"github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	retryDelayFlag       = "retry-delay"
	maxAttemptFlag       = "max-attempts"
	batchDurationFlag    = "batch-duration"
	defaultMaxAttempt    = 3
	defaultRetryDelay    = time.Second
	defaultBatchDuration = 48 * time.Hour
	// as https://huobiapi.github.io/docs/spot/v1/en/#search-past-orders, farthest day is 180 days from now
	// then it should be safe to set to 179 days
	farthestDate = 179 * 24 * time.Hour
	symbolsFlag  = "symbols"
)

var sugar *zap.SugaredLogger

func main() {
	app := libapp.NewApp()
	app.Name = "Huobi Fetcher"
	app.Usage = "Huobi Fetcher for trade logs"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags,
		cli.IntFlag{
			Name:   maxAttemptFlag,
			Usage:  "The maximum number of attempts to retry fetching data",
			EnvVar: "MAX_ATTEMPTS",
			Value:  defaultMaxAttempt,
		},
		cli.DurationFlag{
			Name:   retryDelayFlag,
			Usage:  "The duration to put fetcher job to sleep after each fail attempt",
			EnvVar: "RETRY_DELAY",
			Value:  defaultRetryDelay,
		},
		cli.DurationFlag{
			Name:   batchDurationFlag,
			Usage:  "The duration for a batch query. If the duration is too big, the query will require a lot of memory to store. Default is 90 days each batch",
			EnvVar: "BATCH_DURATION",
			Value:  defaultBatchDuration,
		},
		cli.StringSliceFlag{
			Name:   symbolsFlag,
			Usage:  "symbol to crawl trades",
			EnvVar: "SYMBOLS",
		},
	)
	app.Flags = append(app.Flags, huobi.NewCliFlags()...)
	app.Flags = append(app.Flags, timeutil.NewMilliTimeRangeCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexTradesDB)...)
	if err := app.Run(os.Args); err != nil {
		// log.Fatal does not signal to sentry, using sugar instead
		sugar.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		flush func()
		err   error
	)
	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flush, err = libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	huobiClient, err := huobi.NewClientFromContext(c, sugar)
	if err != nil {
		return err
	}

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot create db from flags: %v", err)
	}

	hdb, err := postgres.NewDB(sugar, db)
	if err != nil {
		return fmt.Errorf("cannot create huobi database instance: %v", err)

	}

	from, err := timeutil.FromTimeMillisFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot get from time: %v", err)
	}
	if from.IsZero() {
		sugar.Info("from timestamp is not provided, get from farthest time")
		from = time.Now().Add(-1 * farthestDate)
	}
	sugar.Infow("fetch trade from time", "time", from)

	to, err := timeutil.ToTimeMillisFromContext(c)
	if err != nil {
		return fmt.Errorf("cannot get to time: %v", err)
	}
	if to.IsZero() {
		to = time.Now()
	}

	retryDelay := c.Duration(retryDelayFlag)
	maxAttempts := c.Int(maxAttemptFlag)
	symbols := c.StringSlice(symbolsFlag)

	huobiSymbols := []huobi.Symbol{}
	for _, symbol := range symbols {
		huobiSymbols = append(huobiSymbols,
			huobi.Symbol{
				SymBol: symbol,
			})
	}

	batchDuration := c.Duration(batchDurationFlag)
	fetcher := huobiFetcher.NewFetcher(sugar, huobiClient, retryDelay, maxAttempts, batchDuration, hdb)
	return fetcher.GetTradeHistory(from, to, huobiSymbols)
}

// func updateTradeRecord(trades map[int64]huobi.TradeHistory, trade huobi.TradeHistory) error {
// 	//if trade does not exist, add it to the list
// 	if _, exist := trades[trade.ID]; !exist {
// 		trades[trade.ID] = trade
// 	} else if trade.State == "partial-filled" || (trade.State == "filled" && trades[trade.ID].State != "filled") {
// 		// if trade is partial-filled multiple time, the amount is the sum of all times
// 		tempTrade := trades[trade.ID]
// 		amount, err := strconv.ParseFloat(tempTrade.Amount, 64)
// 		if err != nil {
// 			return err
// 		}
// 		orderAmount, err := strconv.ParseFloat(trade.FieldAmount, 64)
// 		if err != nil {
// 			return err
// 		}
// 		amount += orderAmount
// 		fee, err := strconv.ParseFloat(tempTrade.FieldFees, 64)
// 		if err != nil {
// 			return err
// 		}
// 		orderFee, err := strconv.ParseFloat(trade.FieldFees, 64)
// 		if err != nil {
// 			return err
// 		}
// 		fee += orderFee
// 		tempTrade.Amount = strconv.FormatFloat(amount, 'f', -1, 64)
// 		tempTrade.FieldFees = strconv.FormatFloat(fee, 'f', -1, 64)
// 		tempTrade.State = trade.State
// 		trades[trade.ID] = tempTrade
// 	}
// 	return nil
// }
