package app

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

// Flag is the cli flag of deployment.
const Flag = "deployment"

// Address is a wrapper of ethereum common Address that supports multiple deployments.
type Address map[deployment.Deployment][]common.Address

// NewAddress returns an Address instance. Address of all deployments should be present.
func NewAddress(prodAddr, stagingAddr, ropstenAddr []common.Address) Address {
	return map[deployment.Deployment][]common.Address{
		deployment.Production: prodAddr,
		deployment.Staging:    stagingAddr,
		deployment.Ropsten:    ropstenAddr,
	}
}

// NewCrossDeploymentAddress returns an Address with given same address for all deployments.
func NewCrossDeploymentAddress(addr []common.Address) Address {
	return map[deployment.Deployment][]common.Address{
		deployment.Production: addr,
		deployment.Staging:    addr,
		deployment.Ropsten:    addr,
	}
}

// MustGetFromContext returns the common address for given deployment from context.
func (a Address) MustGetFromContext(c *cli.Context) []common.Address {
	dpl := c.GlobalString(Flag)
	deploymentMode, err := stringToDeploymentMode(dpl)
	if err != nil {
		panic(err)
	}
	addr, ok := a[deploymentMode]
	if !ok {
		panic(fmt.Errorf("address is not available for deployment: %s", dpl))
	}
	return addr
}

// MustGetOneFromContext returns one common address for given deployment from context
func (a Address) MustGetOneFromContext(c *cli.Context) common.Address {
	dpl := c.GlobalString(Flag)
	deploymentMode, err := stringToDeploymentMode(dpl)
	if err != nil {
		panic(err)
	}
	addr, ok := a[deploymentMode]
	if !ok {
		panic(fmt.Errorf("address is not available for deployment: %s", dpl))
	}
	if len(addr) != 1 {
		panic(fmt.Errorf("address should return only one address for this mode: %s", dpl))
	}
	return addr[0]
}
