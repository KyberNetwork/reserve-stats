package contracts

import (
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/ethereum/go-ethereum/common"
)

// NetworkContractAddress returns the address of network contract of all deployments.
func NetworkContractAddress() deployment.Address {
	return networkContractAddress
}

// InternalReserveAddress returns the address of reserve contract of all deployments.
func InternalReserveAddress() deployment.Address {
	return internalReserveAddress
}

// PricingContractAddress returns the address of pricing contract of all deployments.
func PricingContractAddress() deployment.Address {
	return pricingContractAddress
}

// OldBurnerContractAddress returns old burner address of all deployments.
func OldBurnerContractAddress() deployment.Address {
	return oldBurnerContractAddress
}

// BurnerContractAddress returns the address of burner contract of all deployments.
func BurnerContractAddress() deployment.Address {
	return burnerContractAddress
}

var (
	networkContractAddress = deployment.NewAddress( // we don't have network contract on bsc, this for compatible purpose
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	internalReserveAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("")},
		[]common.Address{common.HexToAddress("")},
		[]common.Address{common.HexToAddress("0xc1F6B9c8c79fD6BE87aB017566A066c10e92e7e0")},
	)
	pricingContractAddress = deployment.NewAddress(
		[]common.Address{},
		[]common.Address{},
		[]common.Address{common.HexToAddress("0xfe96C7566c5f5D5555EbDC88B687456081114EBe")},
	)
	burnerContractAddress = deployment.NewAddress( // we don't have burner for bsc yet, keep this for compatible purpose
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	oldBurnerContractAddress = deployment.NewAddress(
		[]common.Address{}, // production
		[]common.Address{}, // staging
		[]common.Address{}, // ropsten
	)
)
