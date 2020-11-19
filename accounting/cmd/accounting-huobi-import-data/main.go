package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/postgres"
	withdrawstorage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history/postgres"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"
)

const (
	tradeHistoryFileFlag    = "trade-history-file"
	withdrawHistoryFileFlag = "withdraw-history-file"
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
		},
		cli.StringFlag{
			Name:   withdrawHistoryFileFlag,
			Usage:  "huobi withdraw history file",
			EnvVar: "WITHDRAW_HISTORY_FILE",
		},
	)
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(common.DefaultCexTradesDB)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func importTradeHistory(sugar *zap.SugaredLogger, historyFile string, hdb *postgres.HuobiStorage) error {
	var (
		logger         = sugar.With("func", caller.GetCurrentFunctionName())
		types          = []string{"", "buy-market", "sell-market", "buy-limit", "sell-limit"}
		tradeHistories = make(map[int64]huobi.TradeHistory)
	)
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
		orderID, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			logger.Debugw("get order id error", "error", err, "id", id)
			return err
		}
		logger.Infow("order id", "id", orderID)

		updatedAt, err := strconv.ParseUint(line[8], 10, 64)
		if err != nil {
			return err
		}
		logger.Infow("updated at", "time", updatedAt)

		orderType, err := strconv.ParseInt(line[2], 10, 64)
		if err != nil {
			logger.Debugw("get order type error", "error", err, "id", id)
			return err
		}
		logger.Infow("order type", "type", orderType)

		orderState, err := strconv.ParseInt(line[4], 10, 64)
		if err != nil {
			logger.Debugw("get order state error", "error", err, "id", id)
			return err
		}

		orderAmount, err := strconv.ParseFloat(line[6], 64)
		if err != nil {
			logger.Debugw("get order amount error", "error", err, "id", id)
			return err
		}

		orderFee, err := strconv.ParseFloat(line[7], 64)
		if err != nil {
			logger.Debugw("get order fee error", "error", err, "id", id)
			return err
		}

		if _, existed := tradeHistories[orderID]; !existed {
			order := huobi.TradeHistory{
				ID:        orderID,
				Symbol:    line[1],
				Source:    "api",
				Type:      types[orderType],
				Price:     line[5],
				State:     common.OrderState(orderState).String(),
				FieldFees: line[7],
				Amount:    strconv.FormatFloat(orderAmount, 'f', -1, 64),
				// because we use created-at as timestamp to detect the latest time stored
				// however, the data huobi sent us only have one field updated-at related to time then
				// we use this field as created timestamp also
				CreatedAt: updatedAt,
			}
			switch common.OrderState(orderState) {
			case common.PartialCanceled, common.Filled:
				order.FinishedAt = updatedAt
			case common.Canceled:
				order.CanceledAt = updatedAt
			}
			tradeHistories[orderID] = order
		} else if orderState == int64(common.PartialFilled) || (orderState == int64(common.Filled) && tradeHistories[orderID].State != common.Filled.String()) {
			// if order is partial-filled or filled, then the amount should be the sum of all
			order := tradeHistories[orderID]
			amount, err := strconv.ParseFloat(order.Amount, 64)
			if err != nil {
				return err
			}
			amount += orderAmount
			fee, err := strconv.ParseFloat(order.FieldFees, 64)
			if err != nil {
				return err
			}
			fee += orderFee
			order.Amount = strconv.FormatFloat(amount, 'f', -1, 64)
			order.FieldFees = strconv.FormatFloat(fee, 'f', -1, 64)
			order.State = common.OrderState(orderState).String()
			tradeHistories[orderID] = order
		}
	}
	return hdb.UpdateTradeHistory(tradeHistories, "huobi_v1_main")
}

func importWithdrawHistory(sugar *zap.SugaredLogger, historyFile string, hdb *withdrawstorage.HuobiStorage) error {
	var (
		logger            = sugar.With("func", caller.GetCurrentFunctionName())
		withdrawHistories []huobi.WithdrawHistory
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
		withdrawID, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			return err
		}

		amount, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			return err
		}

		fee, err := strconv.ParseFloat(line[4], 64)
		if err != nil {
			return err
		}

		updatedAt, err := strconv.ParseUint(line[5], 10, 64)
		if err != nil {
			return err
		}
		logger.Infow("updated at", "time", updatedAt)

		withdrawHistories = append(withdrawHistories, huobi.WithdrawHistory{
			ID:        withdrawID,
			Currency:  line[1],
			Amount:    amount,
			Fee:       fee,
			Type:      "withdraw",
			TxHash:    line[7],
			Address:   line[6],
			UpdatedAt: updatedAt,
			State:     "confirmed",
		})
	}
	return hdb.UpdateWithdrawHistory(withdrawHistories, "huobi_v1_main") // fixed as import data only for v1
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
	} else {
		sugar.Info("No trade history provided. Skip")
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
