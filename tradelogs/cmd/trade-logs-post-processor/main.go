package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/now"
	"github.com/urfave/cli"
)

const (
	tradeLogsDatabase = "trade_logs"
	timePrecision     = "s"
	reportMeasurement = "monthly_report"
)

type reserveVolume struct {
	ethVolume float64
	usdVolume float64
}

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

func writeReserveVolumeMonthly(influxclient client.Client, timestamp time.Time, ethVolume, usdVolume float64, reserveAddr string) error {
	tags := map[string]string{
		"reserve_addr": reserveAddr,
	}
	fields := map[string]interface{}{
		"eth_volume": ethVolume,
		"usd_volume": usdVolume,
	}
	point, err := client.NewPoint(reportMeasurement, tags, fields, timestamp)
	if err != nil {
		return err
	}
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  tradeLogsDatabase,
		Precision: timePrecision,
	})
	bp.AddPoint(point)
	return influxclient.Write(bp)
}

func fromInfluxResultToMap(res []client.Result, sugar *zap.SugaredLogger, tagKeys string) (map[string]reserveVolume, error) {
	var (
		result = make(map[string]reserveVolume)
	)
	if len(res) < 1 || len(res[0].Series) < 1 {
		sugar.Info("There is no trades")
	}
	for _, row := range res[0].Series {
		if len(row.Values[0]) != 3 {
			sugar.Info("record is not correct format")
			break
		}
		reserveAddr := row.Tags[tagKeys]
		if reserveAddr == "" {
			continue
		}
		ethVolume, err := influxdb.GetFloat64FromInterface(row.Values[0][1])
		if err != nil {
			return nil, err
		}
		usdVolume, err := influxdb.GetFloat64FromInterface(row.Values[0][2])
		if err != nil {
			return nil, err
		}
		result[reserveAddr] = reserveVolume{
			ethVolume: ethVolume,
			usdVolume: usdVolume,
		}
	}
	return result, nil
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

		// reserve volume monthly for dst reserve
		query := fmt.Sprintf(`SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume FROM 
		(SELECT eth_amount, eth_amount*eth_usd_rate as usd_amount FROM trades WHERE time >= '%s' AND time < '%s'
		GROUP BY src_rsv_addr) GROUP BY src_rsv_addr`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("src query ", query)

		res, err := queryDB(influxClient, query)
		if err != nil {
			return err
		}

		srcReserveVolume, err := fromInfluxResultToMap(res, sugar, "src_rsv_addr")
		if err != nil {
			return err
		}

		// reserve volume monthly for dst reserve
		query = fmt.Sprintf(`SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume FROM
		(SELECT eth_amount, eth_amount*eth_usd_rate as usd_amount FROM trades WHERE time >= '%s' AND time < '%s'
		GROUP BY dst_rsv_addr) GROUP BY dst_rsv_addr`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("dst query ", query)

		res, err = queryDB(influxClient, query)
		if err != nil {
			return err
		}

		dstReserveVolume, err := fromInfluxResultToMap(res, sugar, "dst_rsv_addr")
		if err != nil {
			return err
		}

		for reserveAddr, vol := range dstReserveVolume {
			if _, exist := srcReserveVolume[reserveAddr]; exist {
				srcReserveVolume[reserveAddr] = reserveVolume{
					ethVolume: srcReserveVolume[reserveAddr].ethVolume + vol.ethVolume,
					usdVolume: srcReserveVolume[reserveAddr].usdVolume + vol.usdVolume,
				}
			} else {
				srcReserveVolume[reserveAddr] = vol
			}
		}

		for reserveAddr, vol := range srcReserveVolume {
			// insert into report monthly
			if err = writeReserveVolumeMonthly(influxClient, beginOfLastMonth, vol.ethVolume, vol.usdVolume, reserveAddr); err != nil {
				return err
			}
		}

		// burn fee monthly

		query = fmt.Sprintf(`SELECT SUM(amount) as burn_fee INTO %s 
		FROM burn_fees WHERE time >= '%s' AND time < '%s'
		GROUP BY reserve_addr`, reportMeasurement, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("query ", query)

		_, err = queryDB(influxClient, query)
		if err != nil {
			return err
		}

		// wallet fee monthly

		query = fmt.Sprintf(`SELECT SUM(amount) as wallet_fee INTO %s  
		FROM wallet_fees WHERE time >= '%s' AND time < '%s'
		GROUP BY reserve_addr`, reportMeasurement, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("query ", query)

		_, err = queryDB(influxClient, query)
		if err != nil {
			return err
		}

		if beginOfThisMonth.Equal(now.New(time.Now().In(time.UTC)).BeginningOfMonth()) {
			sugar.Info("Sleep for 10 day...")
			time.Sleep(24 * time.Hour)
			startTime = time.Now()
		} else {
			nextMonth := now.New(startTime).EndOfMonth().Add(24 * time.Hour)
			startTime = nextMonth
		}
	}
}
