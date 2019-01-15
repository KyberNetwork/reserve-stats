package app

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"

	"github.com/urfave/cli"
)

// VersionedStartingBlocks is the list of versioned block for each new contract deployment
type VersionedStartingBlocks struct {
	v3 uint64
}

// V3 returns starting block of KyberNetwork v3.
func (v *VersionedStartingBlocks) V3() uint64 {
	return v.v3
}

//DeploymentToStartingBlocks map deployment to its according starting blocks
var DeploymentToStartingBlocks = map[deployment.Deployment]VersionedStartingBlocks{
	deployment.Staging: {
		v3: 6997111,
	},
	deployment.Production: {
		v3: 7024980,
	},
}

//MustGetStartingBlocksFromContext return starting blocks from context
func MustGetStartingBlocksFromContext(c *cli.Context) VersionedStartingBlocks {
	dpl := c.GlobalString(Flag)
	deploymentMode, err := stringToDeploymentMode(dpl)
	if err != nil {
		panic(err)
	}
	result, ok := DeploymentToStartingBlocks[deploymentMode]
	if !ok {
		panic(fmt.Errorf("starting blocks for deployment Mode %s is not supported",
			deploymentMode.String()))
	}
	return result
}
