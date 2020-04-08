package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/now"
	"github.com/urfave/cli"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	storage "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/postprocessor"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs Post Processor"
	app.Usage = "Fetch trade logs on KyberNetwork"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		startTime time.Time
	)
	sugar, flush, err := libapp.NewSugaredLogger(c)
	if err != nil {
		return err
	}
	defer flush()

	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Debugw("influx client initiate successfully", "influx client", influxClient)

	// get first timestamp from db
	q := fmt.Sprintf(`SELECT eth_volume from %s ORDER BY DESC LIMIT 1`, storage.ReportMeasurement)
	sugar.Info(q)
	res, err := influxdb.QueryDB(influxClient, q, storage.TradeLogsDatabase)
	if err != nil {
		return err
	}
	sugar.Info(fmt.Sprintf("%+v", res))
	if len(res) != 1 || len(res[0].Series) != 1 || len(res[0].Series[0].Values) < 1 || len(res[0].Series[0].Values[0]) != 2 {
		q = `SELECT eth_amount FROM trades ORDER BY ASC LIMIT 1`
		res, err := influxdb.QueryDB(influxClient, q, storage.TradeLogsDatabase)
		if err != nil {
			return err
		}
		sugar.Info(fmt.Sprintf("%+v", res))
		if len(res) != 1 || len(res[0].Series) != 1 || len(res[0].Series[0].Values) < 1 || len(res[0].Series[0].Values[0]) != 2 {
			sugar.Info("There is no trade in tradelogs")
			return nil
		}

		startTimeString := res[0].Series[0].Values[0][0].(string)
		startTime, err = time.Parse(time.RFC3339, startTimeString)
		if err != nil {
			return err
		}
	} else {
		startTimeString := res[0].Series[0].Values[0][0].(string)
		startTime, err = time.Parse(time.RFC3339, startTimeString)
		if err != nil {
			return err
		}
	}
	postprocessor := storage.New(influxClient, sugar, storage.TradeLogsDatabase)
	// run each day
	for {
		beginOfThisMonth := now.New(startTime).BeginningOfMonth()
		previousMonth := beginOfThisMonth.Add(-24 * time.Hour)
		beginOfLastMonth := now.New(previousMonth).BeginningOfMonth()
		if err := postprocessor.Run(beginOfLastMonth, beginOfThisMonth); err != nil {
			return err
		}

		if beginOfThisMonth.Equal(now.New(time.Now().In(time.UTC)).BeginningOfMonth()) {
			sugar.Info("Finish aggregating...")
			break
		} else {
			nextMonth := now.New(startTime).EndOfMonth().Add(24 * time.Hour)
			startTime = nextMonth
		}
	}
	return nil
}
