package main

import (
	"log"
	"os"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	tradelogcq "github.com/KyberNetwork/reserve-stats/tradelogs/storage/cq"
)

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs cqs manager"
	app.Usage = "Manage trade logs cqs"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags)
	app.Flags = append(app.Flags, cq.NewCQFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func manageCQFromContext(c *cli.Context, influxClient client.Client, sugar *zap.SugaredLogger) error {
	//Deploy CQ	before get/store trade logs
	cqs, err := tradelogcq.CreateAssetVolumeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	reserveVolumeCqs, err := tradelogcq.CreateReserveVolumeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, reserveVolumeCqs...)
	userVolumeCqs, err := tradelogcq.CreateUserVolumeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, userVolumeCqs...)
	burnFeeCqs, err := tradelogcq.CreateBurnFeeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, burnFeeCqs...)
	walletFeeCqs, err := tradelogcq.CreateWalletFeeCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, walletFeeCqs...)
	summaryCqs, err := tradelogcq.CreateSummaryCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, summaryCqs...)
	walletStatsCqs, err := tradelogcq.CreateWalletStatsCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, walletStatsCqs...)
	countryStatsCqs, err := tradelogcq.CreateCountryCqs(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, countryStatsCqs...)
	integrationVolumeCqs, err := tradelogcq.CreateIntegrationVolumeCq(common.DatabaseName)
	if err != nil {
		return err
	}
	cqs = append(cqs, integrationVolumeCqs...)

	return cq.ManageCQs(c, cqs, influxClient, sugar)
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

	sugar.Info("initialized influxClient successfully: ", influxClient)

	if err = manageCQFromContext(c, influxClient, sugar); err != nil {
		return err
	}

	return nil
}
