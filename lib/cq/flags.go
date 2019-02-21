package cq

import (
	"fmt"
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
	//listAllCqsFlag is set to true to list all cqs will be or ever created in tradelogs
	listAllCqsFlag = "list-all-cqs"
	//dropCqFlag will drop a cq
	dropCqFlag = "drop-cq"
	//executeCqFlag will execute a single cq
	executeCqFlag = "execute-cq"
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
	}
}

// ManageCQs manages the given Continous Queries.
func ManageCQs(c *cli.Context, cqs []*ContinuousQuery, influxClient client.Client, sugar *zap.SugaredLogger) error {
	var (
		deploy  = c.Bool(cqsDeployFlag)
		execute = c.Bool(cqsExecuteFlag)
		listAll = c.Bool(listAllCqsFlag)
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

	if listAll {
		for _, cQuery := range cqs {
			fmt.Println(cQuery.Name)
		}
	}

	if c.String(dropCqFlag) != "" {
		cqToDrop := c.String(dropCqFlag)
		ok := false
		for _, cQuery := range cqs {
			if cQuery.Name == cqToDrop {
				if err := cQuery.Drop(influxClient, sugar); err != nil {
					return err
				}
				ok = true
				sugar.Infow("Drop cq successfully", "cq", cqToDrop)
				break
			}
		}
		if !ok {
			sugar.Infow("cq does not exist", "cq", cqToDrop)
		}
	}

	if c.String(executeCqFlag) != "" {
		cqToExecute := c.String(executeCqFlag)
		ok := false
		for _, cQuery := range cqs {
			if cQuery.Name == cqToExecute {
				if err := cQuery.Execute(influxClient, sugar); err != nil {
					return err
				}
				ok = true
				sugar.Infow("Execute cq successfully", "cq", cqToExecute)
				break
			}
		}
		if !ok {
			sugar.Infow("cq does not exist", "cq", cqToExecute)
		}
	}

	if exit {
		sugar.Info("CQ management process completed, exiting")
		os.Exit(0)
	}

	return nil
}
