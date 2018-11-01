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
type Address map[deployment.Mode][]common.Address

// NewAddress returns an Address instance. Address of all deployments should be present.
func NewAddress(prodAddr, stagingAddr []common.Address) Address {
	return map[deployment.Mode][]common.Address{
		deployment.Production: prodAddr,
		deployment.Staging:    stagingAddr,
	}
}

// NewCrossDeploymentAddress returns an Address with given same address for all deployments.
func NewCrossDeploymentAddress(addr []common.Address) Address {
	return map[deployment.Mode][]common.Address{
		deployment.Production: addr,
		deployment.Staging:    addr,
	}
}

// MustGetFromContext returns the common address for given deployment from context.
func (a Address) MustGetFromContext(c *cli.Context) []common.Address {
	dpl := c.GlobalInt(Flag)
	addr, ok := a[deployment.Mode(dpl)]
	if !ok {
		panic(fmt.Errorf("address is not available for deployment: %d", dpl))
	}
	return addr
}
