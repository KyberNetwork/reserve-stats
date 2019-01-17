package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/workers"
)

const (
	defaultDB = "reserve_stats"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs KYC re-aggregate"
	app.Version = "0.0.1"
	app.Action = run
	app.Flags = append(app.Flags, timeutil.NewTimeRangeCliFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
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
	fromTime, err := timeutil.FromTimeFromContext(c)
	if err != nil {
		return err
	}
	toTime, err := timeutil.ToTimeFromContext(c)
	if err != nil {
		return err
	}
	worker := workers.NewKYCUpdateJob(c, fromTime, toTime)
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
