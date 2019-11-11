package postprocessor

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	schema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/tradelogs-post-processor"
)

const (
	//TradeLogsDatabase is tradelogs database name
	TradeLogsDatabase = "trade_logs"
	//ReportMeasurement is report measurement name
	ReportMeasurement = "monthly_report"
	// timePrecision is precision to save to influx db
	timePrecision = "s"
)

//ReserveVolume is represent for monthly reserve volume
type ReserveVolume struct {
	ETHVolume float64 `json:"eth_volume"`
	USDVolume float64 `json:"usd_Volume"`
}
type FeeVolume struct {
	BurnFee   float64 `json:"burn_fee"`
	WalletFee float64 `json:"wallet_fee"`
}

type PostProcessor struct {
	influxClient client.Client
	logger       *zap.SugaredLogger
	dbName       string
}

func New(influxClient client.Client, logger *zap.SugaredLogger, dbName string) *PostProcessor {
	return &PostProcessor{
		influxClient: influxClient,
		logger:       logger,
		dbName:       dbName,
	}
}
func (p *PostProcessor) Run(beginOfLastMonth, beginOfThisMonth time.Time) error {
	volume, err := p.getVolumeData(beginOfLastMonth, beginOfThisMonth)
	if err != nil {
		return err
	}
	if err = p.writeReserveVolumeMonthly(beginOfLastMonth, volume); err != nil {
		return err
	}

	fee, err := p.getFeeData(beginOfLastMonth, beginOfThisMonth)
	if err != nil {
		return err
	}
	if err = p.writeReserveFeeMonthly(beginOfLastMonth, fee); err != nil {
		return err
	}

	query := fmt.Sprintf(
		`SELECT %[1]s, %[2]s, %[3]s, %[4]s, %[5]s INTO %[6]s FROM %[6]s GROUP BY * FILL(0)`,
		schema.Time.String(),
		schema.BurnFee.String(),
		schema.ETHVolume.String(),
		schema.USDVolume.String(),
		schema.WalletFee.String(),
		ReportMeasurement,
	)
	p.logger.Debug("query ", query)
	_, err = influxdb.QueryDB(p.influxClient, query, p.dbName)
	return errors.Wrap(err, "failed to execute fill zero query")
}

func (p *PostProcessor) fromInfluxResultToMap(res []client.Result, tagKeys string) (map[string]ReserveVolume, error) {
	var (
		result = make(map[string]ReserveVolume)
		logger = p.logger.With("func", caller.GetCurrentFunctionName())
	)
	if len(res) < 1 || len(res[0].Series) < 1 {
		logger.Info("There is no reserve_volume for current query")
	}
	for _, row := range res[0].Series {
		if len(row.Values) < 1 || len(row.Values[0]) != 3 {
			logger.Info("record is not correct format")
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
		result[reserveAddr] = ReserveVolume{
			ETHVolume: ethVolume,
			USDVolume: usdVolume,
		}
	}
	return result, nil
}

func (p *PostProcessor) getVolumeData(beginOfLastMonth, beginOfThisMonth time.Time) (map[string]ReserveVolume, error) {
	var logger = p.logger.With("func", caller.GetCurrentFunctionName())
	// reserve volume monthly for src reserve
	query := fmt.Sprintf(`SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume FROM 
		(SELECT eth_amount, eth_amount*eth_usd_rate as usd_amount FROM trades WHERE time >= '%s' AND time < '%s' AND src_rsv_addr != '' 
		AND ((src_addr != '%s' OR dst_addr != '%s') 
		AND (src_addr != '%s' OR dst_addr != '%s'))
		GROUP BY src_rsv_addr) GROUP BY src_rsv_addr FILL(0)`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339),
		blockchain.ETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.ETHAddr.Hex())

	logger.Debugw("get src reserve volume ", "query", query)

	res, err := influxdb.QueryDB(p.influxClient, query, p.dbName)
	if err != nil {
		return nil, err
	}

	srcReserveVolume, err := p.fromInfluxResultToMap(res, "src_rsv_addr")
	if err != nil {
		return nil, err
	}

	// reserve volume monthly for dst reserve
	query = fmt.Sprintf(`SELECT SUM(eth_amount) AS eth_volume, SUM(usd_amount) AS usd_volume FROM
		(SELECT eth_amount, eth_amount*eth_usd_rate as usd_amount FROM trades WHERE time >= '%s' AND time < '%s' AND dst_rsv_addr != '' 
		AND ((src_addr != '%s' OR dst_addr != '%s') 
		AND (src_addr != '%s' OR dst_addr != '%s'))
		GROUP BY dst_rsv_addr) GROUP BY dst_rsv_addr FILL(0)`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339),
		blockchain.ETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.WETHAddr.Hex(), blockchain.ETHAddr.Hex())

	logger.Debugw("get dst reserve volume ", "query", query)

	res, err = influxdb.QueryDB(p.influxClient, query, p.dbName)
	if err != nil {
		return nil, err
	}

	dstReserveVolume, err := p.fromInfluxResultToMap(res, "dst_rsv_addr")
	if err != nil {
		return nil, err
	}

	for reserveAddr, vol := range dstReserveVolume {
		if _, exist := srcReserveVolume[reserveAddr]; exist {
			srcReserveVolume[reserveAddr] = ReserveVolume{
				ETHVolume: srcReserveVolume[reserveAddr].ETHVolume + vol.ETHVolume,
				USDVolume: srcReserveVolume[reserveAddr].USDVolume + vol.USDVolume,
			}
		} else {
			srcReserveVolume[reserveAddr] = vol
		}
	}
	return srcReserveVolume, nil
}

func (p *PostProcessor) fromInfluxResultToMapWalletFee(res []client.Result, tagKeys string) (map[string]FeeVolume, error) {
	var (
		logger = p.logger.With("func", caller.GetCurrentFunctionName())
		result = make(map[string]FeeVolume)
	)
	if len(res) < 1 || len(res[0].Series) < 1 {
		logger.Info("There is no reserve_volume for current query")
	}
	for _, row := range res[0].Series {
		if len(row.Values) < 1 || len(row.Values[0]) != 3 {
			logger.Info("record is not correct format")
			break
		}
		reserveAddr := row.Tags[tagKeys]
		burnFee, err := influxdb.GetFloat64FromInterface(row.Values[0][1])
		if err != nil {
			return nil, err
		}
		walletFee, err := influxdb.GetFloat64FromInterface(row.Values[0][2])
		if err != nil {
			return nil, err
		}
		result[reserveAddr] = FeeVolume{
			BurnFee:   burnFee,
			WalletFee: walletFee,
		}
	}
	return result, nil
}

func (p *PostProcessor) getFeeData(beginOfLastMonth time.Time, beginOfThisMonth time.Time) (map[string]FeeVolume, error) {
	var logger = p.logger.With("func", caller.GetCurrentFunctionName())
	// reserve fee monthly for src reserve
	query := fmt.Sprintf(`SELECT SUM(src_burn_amount) AS burn_amount, SUM(src_wallet_fee_amount) AS wallet_fee_amount 
		FROM trades 
		WHERE time >= '%s' AND time < '%s' AND src_rsv_addr != ''  
		GROUP BY src_rsv_addr FILL(0)`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))
	logger.Debugw("get src reserve fee ", "query", query)

	res, err := influxdb.QueryDB(p.influxClient, query, p.dbName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query db")
	}

	srcReserveFee, err := p.fromInfluxResultToMapWalletFee(res, "src_rsv_addr")
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert result")
	}
	// reserve fee monthly for dst reserve
	query = fmt.Sprintf(`SELECT SUM(dst_burn_amount) AS burn_amount, SUM(dst_wallet_fee_amount) AS wallet_fee_amount 
		FROM trades 
		WHERE time >= '%s' AND time < '%s' AND dst_rsv_addr != ''  
		GROUP BY dst_rsv_addr FILL(0)`, beginOfLastMonth.Format(time.RFC3339), beginOfThisMonth.Format(time.RFC3339))
	logger.Debugw("get dst reserve fee ", "query", query)

	res, err = influxdb.QueryDB(p.influxClient, query, p.dbName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query db")
	}
	dstReserveFee, err := p.fromInfluxResultToMapWalletFee(res, "dst_rsv_addr")
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert result")
	}

	for reserveAddr, vol := range dstReserveFee {
		if _, exist := srcReserveFee[reserveAddr]; exist {
			srcReserveFee[reserveAddr] = FeeVolume{
				BurnFee:   srcReserveFee[reserveAddr].BurnFee + vol.BurnFee,
				WalletFee: srcReserveFee[reserveAddr].WalletFee + vol.WalletFee,
			}
		} else {
			srcReserveFee[reserveAddr] = vol
		}
	}
	return srcReserveFee, nil
}

// writeReserveVolumeMonthly writes aggregated monthly volume to influx
func (p *PostProcessor) writeReserveVolumeMonthly(timestamp time.Time, reserveVolume map[string]ReserveVolume) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  p.dbName,
		Precision: timePrecision,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create new batch point")
	}
	for reserveAddr, vol := range reserveVolume {
		tags := map[string]string{
			schema.ReserveAddr.String(): reserveAddr,
		}
		fields := map[string]interface{}{
			schema.ETHVolume.String(): vol.ETHVolume,
			schema.USDVolume.String(): vol.USDVolume,
		}
		point, err := client.NewPoint(ReportMeasurement, tags, fields, timestamp)
		if err != nil {
			return errors.Wrap(err, "failed to create new point")
		}
		bp.AddPoint(point)
	}
	return errors.Wrap(p.influxClient.Write(bp), "failed to write batch point")
}

// writeReserveFeeMonthly writes aggregated monthly data for reserve fee to influx
func (p *PostProcessor) writeReserveFeeMonthly(timestamp time.Time, reserveFee map[string]FeeVolume) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  p.dbName,
		Precision: timePrecision,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create new batch point")
	}

	for rsvAddr, fee := range reserveFee {
		tags := map[string]string{
			schema.ReserveAddr.String(): rsvAddr,
		}
		fields := map[string]interface{}{
			schema.BurnFee.String():   fee.BurnFee,
			schema.WalletFee.String(): fee.WalletFee,
		}
		point, err := client.NewPoint(ReportMeasurement, tags, fields, timestamp)
		if err != nil {
			return errors.Wrap(err, "failed to create new point")
		}
		bp.AddPoint(point)
	}
	return errors.Wrap(p.influxClient.Write(bp), "failed to write batch point")
}
