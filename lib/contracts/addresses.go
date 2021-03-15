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
	networkContractAddress = deployment.NewAddress(
		// update address for istanbul fork
		[]common.Address{common.HexToAddress("0x7C66550C9c730B6fdd4C03bc2e73c5462c5F7ACC")},
		[]common.Address{common.HexToAddress("0x9CB7bB6D4795A281860b9Bfb7B1441361Cc9A794")},
		[]common.Address{common.HexToAddress("0x920B322D4B8BAB34fb6233646F5c87F87e79952b")},
	)
	internalReserveAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")},
		[]common.Address{common.HexToAddress("0x2C5a182d280EeB5824377B98CD74871f78d6b8BC")},
		[]common.Address{common.HexToAddress("0xEB52Ce516a8d054A574905BDc3D4a176D3a2d51a")},
	)
	pricingContractAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B")},
		[]common.Address{common.HexToAddress("0xe3E415a7a6c287a95DC68a01ff036828073fD2e6")},
		[]common.Address{common.HexToAddress("0xE16E257a25e287AF50C5651A4c2728b32D7e5ef7")},
	)
	burnerContractAddress = deployment.NewAddress(
		// updated address for istanbul fork
		[]common.Address{common.HexToAddress("0x8007aa43792A392b221DC091bdb2191E5fF626d1")},
		[]common.Address{common.HexToAddress("0x39682A7b8E4A03b2c8dC6DA6E0146Aee4E29A306")},
		[]common.Address{common.HexToAddress("0x06b0fbaba8fba5161f725f2159de1e1d6409c35f")},
	)
	oldBurnerContractAddress = deployment.NewAddress(
		[]common.Address{
			common.HexToAddress("0x4E89bc8484B2c454f2F7B25b612b648c45e14A8e"),
			common.HexToAddress("0x07f6e905f2a1559cd9fd43cb92f8a1062a3ca706"),
			// old burner contract v2
			common.HexToAddress("0xed4f53268bfdFF39B36E8786247bA3A02Cf34B04"),
			// old burner contract v3
			common.HexToAddress("0x52166528FCC12681aF996e409Ee3a421a4e128A3"),
		}, // production
		[]common.Address{
			common.HexToAddress("0xB2cB365D803Ad914e63EA49c95eC663715c2F673"),
			// staging old burner contract v2
			common.HexToAddress("0xd6703974Dc30155d768c058189A2936Cf7C62Da6"),
		}, // staging
		[]common.Address{}, // ropsten
	)
)
