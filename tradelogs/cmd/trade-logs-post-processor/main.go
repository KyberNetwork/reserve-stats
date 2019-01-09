package main

import (
	"fmt"
	"log"
	"os"
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/now"
	"github.com/urfave/cli"
)

const (
	tradeLogsDatabase = "trade_logs"
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

// queryDB convenience function to query the database
func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: tradeLogsDatabase,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func run(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}

	defer logger.Sync()

	sugar := logger.Sugar()

	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	sugar.Debugw("influx client initiate successfully", "influx client", influxClient)

	// get first timestamp from db
	q := fmt.Sprintf(`SELECT eth_amount FROM trades ORDER BY ASC LIMIT 1`)
	res, err := queryDB(influxClient, q)
	if err != nil {
		return err
	}
	sugar.Info(fmt.Sprintf("%+v", res))
	if len(res) != 1 || len(res[0].Series) != 1 || len(res[0].Series[0].Values[0]) != 2 {
		sugar.Info("There is no trade in tradelogs")
		return nil
	}

	startTimeString := res[0].Series[0].Values[0][0].(string)
	startTime, err := time.Parse(time.RFC3339, startTimeString)
	if err != nil {
		return err
	}
	// run each day
	for {
		beginOfThisMonth := now.New(startTime).BeginningOfMonth()
		previousMonth := beginOfThisMonth.Add(-24 * time.Hour)
		beginOfLastMonth := now.New(previousMonth).BeginningOfMonth()

		// reserve volume monthly
		query := fmt.Sprintf(`SELECT SUM(eth_amount) AS eth_volume, SUM(usd_volume) AS usd_volume
		INTO monthly_reserve_volume FROM trades WHERE time >= '%s' AND time < '%s'
		GROUP BY src_rsv_addr, dst_rsv_addr`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("query ", query)

		_, err := queryDB(influxClient, query)
		if err != nil {
			return err
		}

		// burn fee monthly

		query = fmt.Sprintf(`SELECT SUM(amount) as burn_fee INTO monthly_fee 
		FROM burn_fees WHERE time >= '%s' AND time <= '%s'
		GROUP BY reserve_addr`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("query ", query)

		_, err = queryDB(influxClient, query)
		if err != nil {
			return err
		}

		// wallet fee monthly

		query = fmt.Sprintf(`SELECT SUM(amount) as wallet_fee INTO monthly_fee 
		FROM wallet_fees WHERE time >= '%s' AND time <= '%s'
		GROUP BY reserve_addr`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("query ", query)

		_, err = queryDB(influxClient, query)
		if err != nil {
			return err
		}

		if beginOfThisMonth.Equal(now.New(time.Now().In(time.UTC)).BeginningOfMonth()) {
			time.Sleep(24 * time.Hour)
			startTime = time.Now()
		} else {
			nextMonth := now.New(startTime).EndOfMonth().Add(24 * time.Hour)
			startTime = nextMonth
		}
	}
}
