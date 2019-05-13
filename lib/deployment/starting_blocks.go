package deployment

import (
	"fmt"

	"github.com/urfave/cli"
)

// VersionedStartingBlocks is the list of versioned block for each new contract deployment
type VersionedStartingBlocks struct {
	v3 uint64
	v2 uint64
}

// V3 returns starting block of KyberNetwork v3.
func (v *VersionedStartingBlocks) V3() uint64 {
	return v.v3
}

// V2 returns starting block of KyberNetwork v2.
func (v *VersionedStartingBlocks) V2() uint64 {
	return v.v2
}

//StartingBlocks map deployment to its according starting blocks
var StartingBlocks = map[Deployment]VersionedStartingBlocks{
	Staging: {
		v3: 6997111,
		v2: 5864036,
	},
	Production: {
		v3: 7019038,
		v2: 5925999,
	},
}

//MustGetStartingBlocksFromContext return starting blocks from context
func MustGetStartingBlocksFromContext(c *cli.Context) VersionedStartingBlocks {
	deploymentMode := MustGetDeploymentFromContext(c)
	result, ok := StartingBlocks[deploymentMode]
	if !ok {
		panic(fmt.Errorf("starting blocks for deployment Mode %s is not supported",
			deploymentMode.String()))
	}
	return result
}
