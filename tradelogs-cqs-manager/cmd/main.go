package main

import (
	"fmt"
	"log"
	"os"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs-cqs-manager/cq"
)

const (
	listAllCqsFlag = "list-all-cqs"

	dropCqFlag = "drop-cq"

	executeCqFlag = "execute-cq"
)

//CQs map its name to its cq
var CQs map[string]cq.ContinuousQuery

func main() {
	app := libapp.NewApp()
	app.Name = "Trade Logs cqs manager"
	app.Usage = "Manage trade logs cqs"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = append(app.Flags,
		cli.BoolFlag{
			Name:   listAllCqsFlag,
			Usage:  "List all tradelogs cqs",
			EnvVar: "LIST_ALL_CQS",
		},
		cli.StringFlag{
			Name:   dropCqFlag,
			Usage:  "Drop a cq",
			EnvVar: "DROP_CQ",
		},
		cli.StringFlag{
			Name:   executeCqFlag,
			Usage:  "Execute a cq",
			EnvVar: "EXECUTE_CQ",
		},
	)
	app.Flags = append(app.Flags, cq.NewCQFlags()...)
	app.Flags = append(app.Flags, influxdb.NewCliFlags()...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func listAllCqs() error {
	for cq := range CQs {
		fmt.Printf(cq)
	}
	return nil
}

func dropCq(c client.Client, sugar *zap.SugaredLogger, cqName string) error {
	var (
		logger = sugar.With(
			"func", "tradelogs-cq-manager/cq/DropACQ",
		)
	)
	cq, exist := CQs[cqName]
	if !exist {
		logger.Debugw("cq name does not exist", "name", cqName)
		return nil
	}
	return cq.Drop(c, sugar)
}

func executeCq(c client.Client, sugar *zap.SugaredLogger, cqName string) error {
	var (
		logger = sugar.With(
			"func", "tradelogs-cq-manager/cq/DropACQ",
		)
	)
	cq, exist := CQs[cqName]
	if !exist {
		logger.Debugw("cq name does not exist", "name", cqName)
		return nil
	}
	return cq.Execute(c, sugar)
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

	if c.Bool(listAllCqsFlag) {
		if err := listAllCqs(); err != nil {
			return err
		}
	}

	if c.String(dropCqFlag) != "" {
		cqToDrop := c.String(dropCqFlag)
		if err := dropCq(influxClient, sugar, cqToDrop); err != nil {
			return err
		}
	}

	if c.String(executeCqFlag) != "" {
		cqToExecute := c.String(executeCqFlag)
		if err := executeCq(influxClient, sugar, cqToExecute); err != nil {
			return err
		}
	}

	return nil
}
