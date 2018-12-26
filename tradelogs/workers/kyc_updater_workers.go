package workers

import (
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
	"github.com/urfave/cli"
	"go.uber.org/zap"
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
	logger := sugar.With(
		"from", ku.from.String(),
		"to", ku.to.String(),
	)
	logger.Debugw("getting trade logs from DB")
	c := ku.c
	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}
	db, err := libapp.NewDBFromContext(c)
	if err != nil {
		return err
	}
	kycChecker := storage.NewUserKYCChecker(sugar, db)

	tokenAmountFormatter, err := blockchain.NewToKenAmountFormatterFromContext(c)
	if err != nil {
		return err
	}
	influxStorage, err := storage.NewInfluxStorage(
		sugar,
		common.DatabaseName,
		influxClient,
		tokenAmountFormatter,
		kycChecker,
	)
	tradeLogs, err := influxStorage.LoadTradeLogs(ku.from, ku.to)
	if err != nil {
		return err
	}
	logger.Debugw("saving trade logs from DB", "n. trade logs", len(tradeLogs))
	return influxStorage.SaveTradeLogs(tradeLogs)
}
