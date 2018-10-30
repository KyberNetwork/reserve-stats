package contracts

import (
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/ethereum/go-ethereum/common"
)

// InternalNetworkContractAddress returns the address of internal network contract of all deployments.
func InternalNetworkContractAddress() deployment.Address {
	return internalNetworkContractAddress
}

// InternalReserveAddress returns the address of reserve contract of all deployments.
func InternalReserveAddress() deployment.Address {
	return internalReserveAddress
}

// PricingContractAddress returns the address of pricing contract of all deployments.
func PricingContractAddress() deployment.Address {
	return pricingContractAddress
}

// NetworkContractAddress returns the address of network contract of all deployments.
func NetworkContractAddress() deployment.Address {
	return networkContractAddress
}

// BurnerContractAddress returns the address of burner contract of all deployments.
func BurnerContractAddress() deployment.Address {
	return burnerContractAddress
}

var (
	internalNetworkContractAddress = deployment.NewAddress(
		common.HexToAddress("0x91a502C678605fbCe581eae053319747482276b9"),
		common.HexToAddress("0x706aBcE058DB29eB36578c463cf295F180a1Fe9C"),
	)
	internalReserveAddress = deployment.NewAddress(
		common.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"),
		common.HexToAddress("0x2C5a182d280EeB5824377B98CD74871f78d6b8BC"),
	)
	pricingContractAddress = deployment.NewAddress(
		common.HexToAddress("0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B"),
		common.HexToAddress("0xe3E415a7a6c287a95DC68a01ff036828073fD2e6"),
	)
	networkContractAddress = deployment.NewAddress(
		common.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"),
		common.HexToAddress("0xC14f34233071543E979F6A79AA272b0AB1B4947D"),
	)
	burnerContractAddress = deployment.NewAddress(
		common.HexToAddress("0xed4f53268bfdFF39B36E8786247bA3A02Cf34B04"),
		common.HexToAddress("0xd6703974Dc30155d768c058189A2936Cf7C62Da6"),
	)
)
