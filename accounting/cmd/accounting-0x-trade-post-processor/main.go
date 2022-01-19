package main

import (
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/zerox"
	"github.com/KyberNetwork/reserve-stats/accounting/zerox/storage"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"github.com/KyberNetwork/reserve-stats/lib/marketdata"
)

const (
	tradelogsBaseURLFlag    = "tradelogs-base-url"
	defaultTradelogsBaseURL = "https://0x-exchange-proxy.knstats.com"

	fromTimeFlag    = "from-time"
	defaultFromTime = 1633089600 //
	toTimeFlag      = "to-time"
)

func main() {
	app := libapp.NewApp()
	app.Name = "0x convert tradelogs to eth price"
	app.Usage = "trades history from binance"
	app.Action = run

	// app.Flags = append(app.Flags, binance.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexTradesDB)...)
	app.Flags = append(app.Flags, marketdata.NewMarketDataFlags()...)
	app.Flags = append(app.Flags, httputil.NewHTTPCliFlags(httputil.Accounting0xTradesPort)...)
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   tradelogsBaseURLFlag,
			Usage:  "default url to get tradelogs",
			EnvVar: "tradelogs-base-url",
			Value:  defaultTradelogsBaseURL,
		},
		cli.Int64Flag{
			Name:   fromTimeFlag,
			Usage:  "from time to get tradelogs in unix second",
			EnvVar: "FROM_TIME",
		},
		cli.Int64Flag{
			Name:   toTimeFlag,
			Usage:  "to time to get tradelogs in unix second",
			EnvVar: "TO_TIME",
		},
	)

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()
	zap.ReplaceGlobals(logger)
	l := logger.Sugar()

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	st, err := storage.NewZeroxStorage(db, l)
	if err != nil {
		return err
	}

	baseURL := c.String(tradelogsBaseURLFlag)
	marketDataBaseURL := marketdata.GetMarketDataBaseURLFromContext(c)
	marketDataClient := marketdata.NewMarketDataClient(marketDataBaseURL, l)
	binanceClient := binance.NewClient("", "")
	client := zerox.NewZeroXTradelogClient(baseURL, binanceClient, marketDataClient, l)

	ticker := time.NewTicker(time.Minute)
	for ; true; <-ticker.C {
		startTime := c.Int64(fromTimeFlag)
		if startTime == 0 {
			// get latest saved trade timestamp from db
			startTime, err = st.GetLastTradeTimestamp()
			if err != nil {
				return err
			}
			if startTime == 0 {
				startTime = defaultFromTime
			}
		}
		endTime := c.Int64(toTimeFlag)
		if endTime == 0 {
			endTime = time.Now().Unix() // default timestamp is now
		}
		if err := findConvertTrade(startTime, endTime, client, st); err != nil {
			return err
		}
		l.Info("Done. Wait for next loop")
	}
	return nil
}

func findNonETHTrade(trades []zerox.Tradelog) []zerox.Tradelog {
	var (
		tradelogs []zerox.Tradelog
	)
	for _, trade := range trades {
		if trade.InputToken.Symbol != "WETH" && trade.OutputToken.Symbol != "WETH" {
			tradelogs = append(tradelogs, trade)
		}
	}
	return tradelogs
}

func findConvertTrade(startTime, endTime int64, client *zerox.TradelogClient, st *storage.ZeroxStorage) error {
	var (
		nonETHTradelogs []zerox.Tradelog
	)
	fromTime := time.Unix(startTime, 0)
	for {
		toTime := fromTime.Add(24 * time.Hour)
		if toTime.Unix() > endTime {
			toTime = time.Unix(endTime, 0)
		}
		tempTrades, err := client.GetTradelogsFromHTTP(fromTime, toTime)
		if err != nil {
			return err
		}
		nonETHTradelogs = findNonETHTrade(tempTrades)
		convertTrades, err := client.ConvertTrades(nonETHTradelogs)
		if err != nil {
			return err
		}

		if err := st.InsertConvertTrades(convertTrades); err != nil {
			return err
		}

		if err := st.InsertTradelogs(tempTrades); err != nil {
			return err
		}

		if toTime.Unix() >= endTime {
			break
		}

		fromTime = toTime
	}
	return nil
}
