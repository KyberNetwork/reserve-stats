package deployment

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

// TODO: consider merging this to lib/app

// Flag is the cli flag of deployment.
const Flag = "deployment"

// Address is a wrapper of ethereum common Address that supports multiple deployments.
type Address map[string]common.Address

//  Address returns an Address instance. Address of all deployments should be present.
func NewAddress(prodAddr, stagingAddr common.Address) Address {
	return map[string]common.Address{
		Production: prodAddr,
		Staging:    stagingAddr,
	}
}

// NewCrossDeploymentAddress returns an Address with given same address for all deployments.
func NewCrossDeploymentAddress(addr common.Address) Address {
	return map[string]common.Address{
		Production: addr,
		Staging:    addr,
	}
}

// MustGetFromContext returns the common address for given deployment from context.
func (a Address) MustGetFromContext(c *cli.Context) common.Address {
	dpl := c.GlobalString(Flag)
	addr, ok := a[dpl]
	if !ok {
		panic(fmt.Errorf("address is not available for deployment: %s", dpl))
	}
	return addr
}
