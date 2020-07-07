package deployment

import (
	"fmt"

	"github.com/urfave/cli"
)

// VersionedStartingBlocks is the list of versioned block for each new contract deployment
type VersionedStartingBlocks struct {
	v4 uint64
	v3 uint64
	v2 uint64
}

// V4 return starting block of KyberNetwork v4
func (v *VersionedStartingBlocks) V4() uint64 {
	return v.v4
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
		v4: 10378366,
		v3: 6997111,
		v2: 5864036,
	},
	Production: {
		v4: 10410610,
		v3: 7019038,
		v2: 5925999,
	},
	// Ropsten starting blocks are for testing purpose
	// TODO suppose to change for more precise log (if needed)
	Ropsten: {
		v4: 8111008, // this block number is not correct, just pick random for test only
		v3: 6899992,
		v2: 6899991,
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
