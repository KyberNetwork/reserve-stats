package cq

import (
	"os"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

const (
	// cqsDeployFlag is set to true to deploy cqs at start to aggregate in comming data.
	cqsDeployFlag = "cqs-deploy"
	// cqsExecuteFlag is set to true to execute cqs at start to aggregate historical data.
	cqsExecuteFlag = "cqs-execute"
)

// NewCQFlags creates new cli flags for CQs manager.
// Default these flags will be false
func NewCQFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   cqsDeployFlag,
			Usage:  "deploy Continuous Queries on startup",
			EnvVar: "CQS_DEPLOY",
		},
		cli.BoolFlag{
			Name:   cqsExecuteFlag,
			Usage:  "execute all Continuous Queries and exit",
			EnvVar: "CQS_EXECUTE",
		},
	}
}

// ManageCQs manages the given Continous Queries.
func ManageCQs(c *cli.Context, cqs []*ContinuousQuery, influxClient client.Client, sugar *zap.SugaredLogger) error {
	var (
		deploy  = c.Bool(cqsDeployFlag)
		execute = c.Bool(cqsExecuteFlag)
		exit    = deploy || execute
	)

	if deploy {
		for _, cQuery := range cqs {
			if err := cQuery.Deploy(influxClient, sugar); err != nil {
				sugar.Fatalw("failed to deploy CQs", err)
			}
		}
	}

	if execute {
		for _, cQuery := range cqs {
			if err := cQuery.Execute(influxClient, sugar); err != nil {
				sugar.Fatalw("failed to deploy CQs", err)
			}
		}
	}

	if exit {
		sugar.Info("CQ management process completed, exiting")
		os.Exit(0)
	}

	return nil
}
