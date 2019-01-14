package deployment

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

// Flag is the cli Flag of deployment.
const Flag = "deployment"

// Address is a wrapper of ethereum common Address that supports multiple deployments.
type Address map[Deployment][]common.Address

// NewAddress returns an Address instance. Address of all deployments should be present.
func NewAddress(prodAddr, stagingAddr []common.Address) Address {
	return map[Deployment][]common.Address{
		Production: prodAddr,
		Staging:    stagingAddr,
	}
}

// DeploymentFromContext returns deployment from cli context.
func MustGetDeploymentFromContext(c *cli.Context) Deployment {
	dpl := c.GlobalString(Flag)
	switch dpl {
	case Staging.String():
		return Staging
	case Production.String():
		return Production
	default:
		panic(fmt.Errorf("invalid deployment %s", dpl))
	}
}

// NewCrossDeploymentAddress returns an Address with given same address for all deployments.
func NewCrossDeploymentAddress(addr []common.Address) Address {
	return map[Deployment][]common.Address{
		Production: addr,
		Staging:    addr,
	}
}

// MustGetFromContext returns the common address for given deployment from context.
func (a Address) MustGetFromContext(c *cli.Context) []common.Address {
	dpl := MustGetDeploymentFromContext(c)
	addr, ok := a[dpl]
	if !ok {
		panic(fmt.Errorf("address is not available for deployment: %s", dpl))
	}
	return addr
}

// MustGetOneFromContext returns one common address for given deployment from context
func (a Address) MustGetOneFromContext(c *cli.Context) common.Address {
	dpl := MustGetDeploymentFromContext(c)
	addr, ok := a[dpl]
	if !ok {
		panic(fmt.Errorf("address is not available for deployment: %s", dpl))
	}
	if len(addr) != 1 {
		panic(fmt.Errorf("address should return only one address for this mode: %s", dpl))
	}
	return addr[0]
}
