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
		[]common.Address{common.HexToAddress("0x214d6546f862c2d36ab6cc7f3ddb3423579dfad3")},
		[]common.Address{common.HexToAddress("0xb3F787E2A95326C9dBb0e43151e9B94B7f850e7b")},
		[]common.Address{common.HexToAddress("0xc1F6B9c8c79fD6BE87aB017566A066c10e92e7e0")},
	)
	pricingContractAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0x58F9546fdE73Ca379fB388BA83292dD9D366156F")},
		[]common.Address{common.HexToAddress("0xD14b35Af7316260454B83739aA377667fDA22627")},
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
