package cq

import "github.com/urfave/cli"

const (
	// CqsDeployFlag is set to true to deploy cqs at start to aggregate in comming data.
	CqsDeployFlag = "cqs-deploy"
	// CqsExecuteFlag is set to true to execute cqs at start to aggregate historical data.
	CqsExecuteFlag = "cqs-execute"
)

// NewCQFlags creates new cli flags for CQs manager.
// Default these flags will be false
func NewCQFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:   CqsDeployFlag,
			Usage:  "deploy Continuous Queries on startup",
			EnvVar: "CQS_DEPLOY",
		},
		cli.BoolFlag{
			Name:   CqsExecuteFlag,
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
