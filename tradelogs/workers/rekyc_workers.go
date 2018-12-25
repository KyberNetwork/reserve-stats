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

// ReKYCJob describe a
type ReKYCJob struct {
	c     *cli.Context
	order int
	from  time.Time
	to    time.Time
}

// NewReKYCJob returns a ReKYCJob
func NewReKYCJob(c *cli.Context, order int, from, to time.Time) *ReKYCJob {
	return &ReKYCJob{
		c:     c,
		order: order,
		from:  from,
		to:    to,
	}
}

// Execute get the tradelogs and re-kyc it
func (rk *ReKYCJob) Execute(sugar *zap.SugaredLogger) error {
	return rk.getAndReKYC(sugar)
}

func (rk *ReKYCJob) getAndReKYC(sugar *zap.SugaredLogger) error {
	logger := sugar.With(
		"from", rk.from.String(),
		"to", rk.to.String(),
	)
	logger.Debugw("getting trade logs from DB")
	c := rk.c
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
	tradeLogs, err := influxStorage.LoadTradeLogs(rk.from, rk.to)
	if err != nil {
		return err
	}
	logger.Debugw("saving trade logs from DB", "n. trade logs", len(tradeLogs))
	return influxStorage.SaveTradeLogs(tradeLogs)
}
