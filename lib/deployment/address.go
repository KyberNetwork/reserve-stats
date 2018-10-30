package deployment

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// Address is a wrapper of ethereum common Address that supports multiple deployments.
type Address map[Deployment]common.Address

//  Address returns an Address instance. Address of all deployments should be present.
func NewAddress(prodAddr, stagingAddr common.Address) Address {
	return map[Deployment]common.Address{
		ProdDeployment:    prodAddr,
		StagingDeployment: stagingAddr,
	}
}

// NewDuplicatedAddress returns an Address with given common address for all deployments.
func NewDuplicatedAddress(addr common.Address) Address {
	return map[Deployment]common.Address{
		ProdDeployment:    addr,
		StagingDeployment: addr,
	}
}

func (a Address) MustGet() common.Address {
	addr, ok := a[dpl]
	if !ok {
		panic(fmt.Errorf("address is not available for deployment: %s", dpl.String()))
	}
	return addr
}
