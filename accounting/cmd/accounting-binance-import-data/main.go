package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	withdrawstorage "github.com/KyberNetwork/reserve-stats/accounting/binance/storage/withdrawalstorage"
	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	withdrawHistoryFileFlag = "withdraw-history-file"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Binance importer"
	app.Usage = "Binance importer for withdraw fee"
	app.Action = run
	app.Version = "0.0.1"
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   withdrawHistoryFileFlag,
			Usage:  "binance withdraw history file",
			EnvVar: "WITHDRAW_HISTORY_FILE",
		},
	)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexWithdrawalsDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func importWithdrawHistory(sugar *zap.SugaredLogger, historyFile string, hdb *withdrawstorage.BinanceStorage) error {
	var (
		logger            = sugar.With("func", caller.GetCurrentFunctionName())
		withdrawHistories []binance.WithdrawHistory
	)
	logger.Infow("import withdraw history from file", "file", historyFile)

	csvFile, err := os.Open(historyFile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for id, line := range lines {
		if id == 0 {
			continue
		}

		amount, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			return err
		}

		fee, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			return err
		}

		applyTime, err := time.Parse("2006-01-02 15:04:05", line[0])
		if err != nil {
			fmt.Println(err)
		}
		logger.Infow("apply time", "time", applyTime)
		applyTimeMs := timeutil.TimeToTimestampMs(applyTime)

		status := int64(common.WithdrawStatuses[line[8]])

		withdrawHistories = append(withdrawHistories, binance.WithdrawHistory{
			Asset:     line[1],
			Amount:    amount,
			TxFee:     fee,
			Address:   line[4],
			TxID:      line[5],
			ApplyTime: applyTimeMs,
			Status:    status,
		})
	}
	return hdb.UpdateWithdrawHistoryWithFee(withdrawHistories, "binance_v1_main")
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

	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}

	withdrawFile := c.String(withdrawHistoryFileFlag)

	wdb, err := withdrawstorage.NewDB(sugar, db)
	if err != nil {
		return err
	}

	if withdrawFile != "" {
		if err := importWithdrawHistory(sugar, withdrawFile, wdb); err != nil {
			return err
		}
	} else {
		sugar.Info("No withdraw history file provided. Skip")
	}

	return nil
}
