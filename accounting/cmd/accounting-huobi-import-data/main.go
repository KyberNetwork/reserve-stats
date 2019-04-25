package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

const (
	tradeHistoryFileFlag            = "trade-history-file"
	withdrawHistoryFileFlag         = "withdraw-history-file"
	defaultTradeHistoryFileValue    = "huobi_trade_history.csv"
	defaultWithdrawHistoryFileValue = "huobi_withdraw_history.csv"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Huobi Fetcher"
	app.Usage = "Huobi Fetcher for trade logs"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   tradeHistoryFileFlag,
			Usage:  "huobi trade history file",
			EnvVar: "TRADE_HISTORY_FILE",
			Value:  filepath.Join(currentLocation(), defaultTradeHistoryFileValue),
		},
		cli.StringFlag{
			Name:   withdrawHistoryFileFlag,
			Usage:  "huobi withdraw history file",
			EnvVar: "WITHDRAW_HISTORY_FILE",
			Value:  filepath.Join(currentLocation(), defaultWithdrawHistoryFileValue),
		},
	)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexTradesDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func currentLocation() string {
	_, fileName, _, _ := runtime.Caller(0)
	return filepath.Dir(fileName)
}

func importTradeHistory(sugar *zap.SugaredLogger, historyFile string, hdb *postgres.HuobiStorage) error {
	csvFile, err := os.Open(historyFile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	var tradeHistories []huobi.TradeHistory
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, line := range lines {
		accountID, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			return err
		}
		userID, err := strconv.ParseInt(line[2], 10, 64)
		if err != nil {
			return err
		}
		tradeHistories = append(tradeHistories, huobi.TradeHistory{
			//TODO: match trade history from csv file to go object
			AccountID: accountID,
			UserID:    userID,
			Symbol:    line[3],
		})
	}
	return hdb.UpdateTradeHistory(tradeHistories)
}

func run(c *cli.Context) error {
	if err := libapp.Validate(c); err != nil {
		return err
	}

	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	historyFile := c.String(tradeHistoryFileFlag)

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	hdb, err := postgres.NewDB(sugar, db)
	if err != nil {
		return err
	}
	if historyFile != "" {
		if err := importTradeHistory(sugar, historyFile, hdb); err != nil {
			return err
		}
	}
	return nil
}
