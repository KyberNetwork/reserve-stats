package main

import (
	"fmt"
	"log"
	"os"
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	schema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelogs-post-processor"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jinzhu/now"
	"github.com/urfave/cli"
	"go.uber.org/zap"
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

func fromInfluxResultToMap(res []client.Result, sugar *zap.SugaredLogger, tagKeys string) (map[string]storage.ReserveVolume, error) {
	var (
		result = make(map[string]storage.ReserveVolume)
	)
	if len(res) < 1 || len(res[0].Series) < 1 {
		sugar.Info("There is no reserve_volume for current query")
	}
	for _, row := range res[0].Series {
		if len(row.Values) < 1 || len(row.Values[0]) != 3 {
			sugar.Info("record is not correct format")
			break
		}
		reserveAddr := row.Tags[tagKeys]
		ethVolume, err := influxdb.GetFloat64FromInterface(row.Values[0][1])
		if err != nil {
			return nil, err
		}
		usdVolume, err := influxdb.GetFloat64FromInterface(row.Values[0][2])
		if err != nil {
			return nil, err
		}
		result[reserveAddr] = storage.ReserveVolume{
			ETHVolume: ethVolume,
			USDVolume: usdVolume,
		}
	}
	return result, nil
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
		q = fmt.Sprintf(`SELECT eth_amount FROM trades ORDER BY ASC LIMIT 1`)
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
	// run each day
	for {
		beginOfThisMonth := now.New(startTime).BeginningOfMonth()
		previousMonth := beginOfThisMonth.Add(-24 * time.Hour)
		beginOfLastMonth := now.New(previousMonth).BeginningOfMonth()
		// reserve volume monthly for dst reserve
		query := fmt.Sprintf(`SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume FROM 
		(SELECT eth_amount, eth_amount*eth_usd_rate as usd_amount FROM trades WHERE time >= '%s' AND time < '%s' AND src_rsv_addr != '' 
		AND ((src_addr != '%s' OR dst_addr != '%s') 
		AND (src_addr != '%s' OR dst_addr != '%s'))
		GROUP BY src_rsv_addr) GROUP BY src_rsv_addr FILL(0)`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339),
			blockchain.ETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.ETHAddr.Hex())

		sugar.Debug("src query ", query)

		res, err := influxdb.QueryDB(influxClient, query, storage.TradeLogsDatabase)
		if err != nil {
			return err
		}

		srcReserveVolume, err := fromInfluxResultToMap(res, sugar, "src_rsv_addr")
		if err != nil {
			return err
		}

		// reserve volume monthly for dst reserve
		query = fmt.Sprintf(`SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume FROM
		(SELECT eth_amount, eth_amount*eth_usd_rate as usd_amount FROM trades WHERE time >= '%s' AND time < '%s' AND dst_rsv_addr != '' 
		AND ((src_addr != '%s' OR dst_addr != '%s') 
		AND (src_addr != '%s' OR dst_addr != '%s'))
		GROUP BY dst_rsv_addr) GROUP BY dst_rsv_addr FILL(0)`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339),
			blockchain.ETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.ETHAddr.Hex())

		sugar.Debug("dst query ", query)

		res, err = influxdb.QueryDB(influxClient, query, storage.TradeLogsDatabase)
		if err != nil {
			return err
		}

		dstReserveVolume, err := fromInfluxResultToMap(res, sugar, "dst_rsv_addr")
		if err != nil {
			return err
		}

		for reserveAddr, vol := range dstReserveVolume {
			if _, exist := srcReserveVolume[reserveAddr]; exist {
				srcReserveVolume[reserveAddr] = storage.ReserveVolume{
					ETHVolume: srcReserveVolume[reserveAddr].ETHVolume + vol.ETHVolume,
					USDVolume: srcReserveVolume[reserveAddr].USDVolume + vol.USDVolume,
				}
			} else {
				srcReserveVolume[reserveAddr] = vol
			}
		}

		if err = storage.WriteReserveVolumeMonthly(influxClient, beginOfLastMonth, srcReserveVolume); err != nil {
			return err
		}

		// burn fee monthly

		query = fmt.Sprintf(`SELECT SUM(amount) as burn_fee INTO %s 
		FROM burn_fees WHERE time >= '%s' AND time < '%s'
		GROUP BY reserve_addr FILL(0)`, storage.ReportMeasurement, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("query ", query)

		_, err = influxdb.QueryDB(influxClient, query, storage.TradeLogsDatabase)
		if err != nil {
			return err
		}

		// wallet fee monthly

		query = fmt.Sprintf(`SELECT SUM(amount) as wallet_fee INTO %s  
		FROM wallet_fees WHERE time >= '%s' AND time < '%s'
		GROUP BY reserve_addr FILL(0)`, storage.ReportMeasurement, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))

		sugar.Debug("query ", query)

		_, err = influxdb.QueryDB(influxClient, query, storage.TradeLogsDatabase)
		if err != nil {
			return err
		}

		query = fmt.Sprintf(
			`SELECT %[1]s, %[2]s, %[3]s, %[4]s, %[5]s INTO %[6]s FROM %[6]s GROUP BY * FILL(0)`,
			schema.Time.String(),
			schema.BurnFee.String(),
			schema.ETHVolume.String(),
			schema.USDVolume.String(),
			schema.WalletFee.String(),
			storage.ReportMeasurement,
		)
		sugar.Debug("query ", query)

		_, err = influxdb.QueryDB(influxClient, query, storage.TradeLogsDatabase)
		if err != nil {
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
