package influxstorage

import (
	"time"

	schema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influxstorage/schema/tradelogs-post-processor"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	//TradeLogsDatabase is tradelogs database name
	TradeLogsDatabase = "trade_logs"
	//ReportMeasurement is report measurement name
	ReportMeasurement = "monthly_report"
)

//ReserveVolume is represent for monthly reserve volume
type ReserveVolume struct {
	ETHVolume float64 `json:"eth_volume"`
	USDVolume float64 `json:"usd_Volume"`
}

//WriteReserveVolumeMonthly write aggregated monthly data to databases
func WriteReserveVolumeMonthly(influxclient client.Client, timestamp time.Time, reserveVolume map[string]ReserveVolume) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  TradeLogsDatabase,
		Precision: timePrecision,
	})
	if err != nil {
		return err
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
			return err
		}
		bp.AddPoint(point)
	}
	return influxclient.Write(bp)
}
