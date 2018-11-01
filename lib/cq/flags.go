package cq

import "github.com/urfave/cli"

const (
	cqsDeployFlag  = "--cqs-deploy"
	cqsExecuteFlag = "--cqs-execute"
)

// NewCQFlags creates new cli flags for CQs manager.
// TODO: integrates this to trade-logs-crawler and reserve-rates-crawler
func NewCQFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolTFlag{
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
func ManageCQs(c *cli.Context, cqs []ContinuousQuery) error {
	// TODO: check if cqsDeploy == true --> deploy CQs

	// TODO: check if cqsExecute == true --> run all queries and exit

	return nil
}
