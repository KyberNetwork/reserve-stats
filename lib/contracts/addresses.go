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

// ProxyContractAddress returns the address of network contract of all deployments.
func ProxyContractAddress() deployment.Address {
	return proxyContractAddress
}

// BurnerContractAddress returns the address of burner contract of all deployments.
func BurnerContractAddress() deployment.Address {
	return burnerContractAddress
}

// OldNetworkContractAddress returns old network address of all deployments.
func OldNetworkContractAddress() deployment.Address {
	return oldNetworkContractAddress
}

// OldBurnerContractAddress returns old burner address of all deployments.
func OldBurnerContractAddress() deployment.Address {
	return oldBurnerContractAddress
}

// OldProxyContractAddress returns old proxy address of all deployment
func OldProxyContractAddress() deployment.Address {
	return oldProxyContractAddress
}

// VolumeExcludedReserves return volume excluded reserve of all deployments
func VolumeExcludedReserves() deployment.Address {
	return volumeExcludedReserves
}

// KyberStorageContractAddress return contract address of kyber storage for all deployments
func KyberStorageContractAddress() deployment.Address {
	return kyberStorageContractAddress
}

// KyberFeeHandlerContractAddress return contract address of kyber fee handler for all deployments
func KyberFeeHandlerContractAddress() deployment.Address {
	return feeHandlerContractAddress
}

func OldFeeHandlerContractAddress() deployment.Address {
	return oldFeeHandlerContractAddress
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
	proxyContractAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0x9AAb3f75489902f3a48495025729a0AF77d4b11e")},
		[]common.Address{common.HexToAddress("0xc153eeAD19e0DBbDb3462Dcc2B703cC6D738A37c")},
		[]common.Address{common.HexToAddress("0xd719c34261e099Fdb33030ac8909d5788D3039C4")},
	)
	burnerContractAddress = deployment.NewAddress(
		// updated address for istanbul fork
		[]common.Address{common.HexToAddress("0x8007aa43792A392b221DC091bdb2191E5fF626d1")},
		[]common.Address{common.HexToAddress("0x39682A7b8E4A03b2c8dC6DA6E0146Aee4E29A306")},
		[]common.Address{common.HexToAddress("0x06b0fbaba8fba5161f725f2159de1e1d6409c35f")},
	)

	feeHandlerContractAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0x9Fb131eFbac23b735d7764AB12F9e52cC68401CA")}, // production
		[]common.Address{common.HexToAddress("0xEc30037C9A8A6A3f42734c30Dfa0a208aF71b40C")}, // staging
		[]common.Address{common.HexToAddress("0xfF456D9A8cbB5352eF77dEc2337bAC8dEC63bEAC")}, // ropsten
	)

	kyberStorageContractAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0xC8fb12402cB16970F3C5F4b48Ff68Eb9D1289301")},
		[]common.Address{common.HexToAddress("0xB18D90bE9ADD2a6c9F2c3943B264c3dC86E30cF5")},
		[]common.Address{common.HexToAddress("0x688bf5EeC43E0799c5B9c1612F625F7b93FE5434")},
	)

	oldFeeHandlerContractAddress = deployment.NewAddress(
		[]common.Address{
			common.HexToAddress("0xd3d2b5643e506c6d9B7099E9116D7aAa941114fe"),
		},
		[]common.Address{},
		[]common.Address{},
	)

	oldProxyContractAddress = deployment.NewAddress(
		[]common.Address{
			common.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"),
		},
		[]common.Address{
			common.HexToAddress("0x65897aDCBa42dcCA5DD162c647b1cC3E31238490"),
			// old proxy contract v3
			common.HexToAddress("0x6326dd73E368c036D4C4997053a021CBc52c7367"),
		},
		[]common.Address{
			common.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"),
		},
	)

	oldNetworkContractAddress = deployment.NewAddress(
		[]common.Address{
			common.HexToAddress("0x964F35fAe36d75B1e72770e244F6595B68508CF5"),
			// production old network v2
			common.HexToAddress("0x91a502C678605fbCe581eae053319747482276b9"),
			// production old network v3
			common.HexToAddress("0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950"),
			//
			common.HexToAddress("0x65bF64Ff5f51272f729BDcD7AcFB00677ced86Cd"),
		},
		[]common.Address{
			common.HexToAddress("0xD2D21FdeF0D054D2864ce328cc56D1238d6b239e"),
			// staging old network v2
			common.HexToAddress("0x706aBcE058DB29eB36578c463cf295F180a1Fe9C"),
			// staging old network contract
			common.HexToAddress("0xC14f34233071543E979F6A79AA272b0AB1B4947D"),
			// staging old network contract v3
			common.HexToAddress("0xafBf0D08269a7eEe8d587121f3B0616c8CeF5077"),
		},
		[]common.Address{
			common.HexToAddress("0x753fe1914db38ee744e071baadd123f50f9c8e46"),
		},
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

	volumeExcludedReserves = deployment.NewAddress(
		[]common.Address{
			common.HexToAddress("0x2295fc6BC32cD12fdBb852cFf4014cEAc6d79C10"), // PT Reserve
			common.HexToAddress("0x57f8160e1c59D16C01BbE181fD94db4E56b60495"), // WETH Reserve
			common.HexToAddress("0x0000000000000000000000000000000000000000"), // Self Reserve
		},
		[]common.Address{
			common.HexToAddress("0x0000000000000000000000000000000000000000"), // Self Reserve
			common.HexToAddress("0x29382a4c3b22a39B83c76F261439bBCC78c72dd0"), // PT Reserve
		}, // staging
		[]common.Address{
			common.HexToAddress("0x0000000000000000000000000000000000000000"), // Self Reserve
		}, // ropsten
	)
)
