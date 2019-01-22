package main

import (
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs-cqs-manager/cq"
	"github.com/urfave/cli"
)

const (
	listAllCqsFlag = "list-all-cqs"

	dropCqFlag = "drop-cq"

	executeCqFlag = "execute-cq"
)

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
	return nil
}

func dropCq(cq string) error {
	return nil
}

func executeCq(cq string) error {
	return nil
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

	sugar.Info("influxClient: ", influxClient)

	if c.Bool(listAllCqsFlag) {
		if err := listAllCqs(); err != nil {
			return err
		}
	}

	if c.String(dropCqFlag) != "" {
		cqToDrop := c.String(dropCqFlag)
		if err := dropCq(cqToDrop); err != nil {
			return err
		}
	}

	if c.String(executeCqFlag) != "" {
		cqToExecute := c.String(executeCqFlag)
		if err := executeCq(cqToExecute); err != nil {
			return err	
		}
	}

	return nil
}
