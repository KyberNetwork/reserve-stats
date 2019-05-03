package kycupdater

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

// KYCUpdateJob describe a
type KYCUpdateJob struct {
	c    *cli.Context
	from time.Time
	to   time.Time
}

// NewKYCUpdateJob returns a KYCUpdateJob
func NewKYCUpdateJob(c *cli.Context, from, to time.Time) *KYCUpdateJob {
	return &KYCUpdateJob{
		c:    c,
		from: from,
		to:   to,
	}
}

// Execute get the tradelogs and re-kyc it
func (ku *KYCUpdateJob) Execute(sugar *zap.SugaredLogger) error {
	return ku.update(sugar)
}

func (ku *KYCUpdateJob) update(sugar *zap.SugaredLogger) error {
	var (
		step   = time.Hour * 24 * 30 // 30 days
		logger = sugar.With(
			"from", ku.from.String(),
			"to", ku.to.String(),
		)
	)
	logger.Debugw("getting trade logs from DB")
	c := ku.c
	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}

	tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
	if err != nil {
		return err
	}
	influxStorage, err := storage.NewInfluxStorage(
		sugar,
		common.DatabaseName,
		influxClient,
		tokenAmountFormatter,
	)

	if err != nil {
		return err
	}

	for ts := ku.from; ku.to.After(ts); ts = ts.Add(step) {
		var (
			point     *client.Point
			kycPoints []*client.Point
			bp        client.BatchPoints
		)
		logger.Infow("updating KYC information",
			"step_from", ts,
			"step_to", ts.Add(step),
		)

		tradeLogs, err := influxStorage.LoadTradeLogs(ts, ts.Add(step))
		if err != nil {
			return err
		}

		for _, logItem := range tradeLogs {
			if point, err = influxStorage.AssembleKYCPoint(logItem); err != nil {
				return err
			}
			kycPoints = append(kycPoints, point)
		}

		bp, err = client.NewBatchPoints(client.BatchPointsConfig{
			Database:  common.DatabaseName,
			Precision: "s",
		})
		if err != nil {
			return err
		}

		for _, p := range kycPoints {
			bp.AddPoint(p)
		}

		if err = influxClient.Write(bp); err != nil {
			return err
		}
	}
	return nil
}
