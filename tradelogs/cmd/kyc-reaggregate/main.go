package main

import (
	"log"
	"os"
	"time"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/workers"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	defaultDB    = "reserve_stats"
	fromTimeFlag = "fromTime"
	toTimeFlag   = "toTime"
	//This is 01 Feb 2018 UTC
	defaultFromTime = 1517443200000
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs KYC re-aggregate"
	app.Version = "0.0.1"
	app.Action = run
	app.Flags = append(app.Flags,
		cli.Uint64Flag{
			Name:   fromTimeFlag,
			Usage:  "re-KYC from time. Default is 00:00:00",
			EnvVar: "FROM_TIME",
			Value:  defaultFromTime,
		},
		cli.Uint64Flag{
			Name:   toTimeFlag,
			Usage:  "re-KYC to time",
			EnvVar: "TO_TIME",
			Value:  0,
		},
	)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, libapp.NewEthereumNodeFlags())
	app.Flags = append(app.Flags, libapp.NewPostgreSQLFlags(defaultDB)...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		err error
	)
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	fromTime := timeutil.TimestampMsToTime(c.Uint64(fromTimeFlag))
	toTime := timeutil.TimestampMsToTime(c.Uint64(toTimeFlag))
	if c.Uint64(toTimeFlag) == 0 {
		toTime = time.Now()
	}
	worker := workers.NewReKYCJob(c, 0, fromTime, toTime)
	if err := worker.Execute(sugar); err != nil {
		return err
	}

	if err = reRunCqs(c, sugar); err != nil {
		return err
	}

	return nil
}

func reRunCqs(c *cli.Context, sugar *zap.SugaredLogger) error {
	influxClient, err := influxdb.NewClientFromContext(c)
	if err != nil {
		return err
	}
	var cqs []*cq.ContinuousQuery
	summaryCqs, err := tradelogcq.CreateSummaryCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, summaryCqs...)

	countryStatsCqs, err := tradelogcq.CreateCountryCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, countryStatsCqs...)

	walletStatsCqs, err := tradelogcq.CreateWalletStatsCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, walletStatsCqs...)

	for _, cQuery := range cqs {
		if err := cQuery.Execute(influxClient, sugar); err != nil {
			sugar.Fatalw("failed to execute CQs", err)
		}
	}
	return nil
}
