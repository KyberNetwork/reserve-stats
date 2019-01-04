package app

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"

	"github.com/urfave/cli"
)

//VersionedStartingBlocks is the list of versioned block for each new contract deployment
type VersionedStartingBlocks struct {
	V1 uint64
	V2 uint64
	V3 uint64
}

//DeploymentToStartingBlocks map deployment to its according starting blocks
var DeploymentToStartingBlocks = map[deployment.Deployment]VersionedStartingBlocks{
	deployment.Staging: VersionedStartingBlocks{
		V1: 0,
		V2: 0,
		V3: 6997111,
	},
	deployment.Production: VersionedStartingBlocks{
		V1: 0,
		V2: 5926056,
		V3: 0,
	},
}

//GetStartingBlocksFromContext return starting blocks from context
func GetStartingBlocksFromContext(c *cli.Context) (VersionedStartingBlocks, error) {
	dpl := c.GlobalString(Flag)
	deploymentMode, err := stringToDeploymentMode(dpl)
	if err != nil {
		return VersionedStartingBlocks{}, err
	}
	result, ok := DeploymentToStartingBlocks[deploymentMode]
	if !ok {
		return result, fmt.Errorf("starting blocks for deployment Mode %s is not supported", deploymentMode.String())
	}
	return result, nil
}
