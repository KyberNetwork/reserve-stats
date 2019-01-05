package contracts

import (
	"context"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	timeout = 30 * time.Second
)

//InternalNetworkContractAddress returns the address of internal network of all deployments
func InternalNetworkContractAddress(proxyAddress common.Address, client bind.ContractBackend, blockNumber *big.Int) (common.Address, error) {
	proxyContract, err := NewProxy(proxyAddress, client)
	if err != nil {
		return common.Address{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// if not input blocknumber then get from latest block
	callOpt := &bind.CallOpts{Context: ctx}
	if blockNumber.Cmp(big.NewInt(1)) != 0 {
		callOpt.BlockNumber = blockNumber
	}
	return proxyContract.KyberNetworkContract(callOpt)
}

// PricingContractAddress returns the address of pricing contract of all deployments.
func PricingContractAddress(internalNetworkAddress common.Address, client bind.ContractBackend, blockNumber *big.Int) (common.Address, error) {
	internalNetworkContract, err := NewInternalNetwork(internalNetworkAddress, client)
	if err != nil {
		return common.Address{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// if not input blocknumber then get from latest block
	callOpt := &bind.CallOpts{Context: ctx}
	if blockNumber.Cmp(big.NewInt(1)) != 0 {
		callOpt.BlockNumber = blockNumber
	}
	return internalNetworkContract.ExpectedRateContract(callOpt)
}

// InternalReserveAddress returns the address of reserve contract of all deployments.
func InternalReserveAddress() deployment.Address {
	return internalReserveAddress
}

//ProxyContractAddress returns the address of proxy contract of all deployments
func ProxyContractAddress() deployment.Address {
	return proxyContractAddress
}

// BurnerContractAddress returns the address of burner contract of all deployments.
func BurnerContractAddress(internalNetworkAddress common.Address, client bind.ContractBackend, blockNumber *big.Int) (common.Address, error) {
	internalNetworkContract, err := NewInternalNetwork(internalNetworkAddress, client)
	if err != nil {
		return common.Address{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// if not input blocknumber then get from latest block
	callOpt := &bind.CallOpts{Context: ctx}
	if blockNumber.Cmp(big.NewInt(1)) != 0 {
		callOpt.BlockNumber = blockNumber
	}
	return internalNetworkContract.FeeBurnerContract(callOpt)
}

// OldNetworkContractAddress returns old network address of all deployments.
// func OldNetworkContractAddress() app.Address {
// 	return oldNetworkContractAddress
// }

// OldBurnerContractAddress returns old burner address of all deployments.
// func OldBurnerContractAddress() app.Address {
// 	return oldBurnerContractAddress
// }

var (
	// internalNetworkContractAddress = app.NewAddress(
	// 	[]common.Address{common.HexToAddress("0x9ae49C0d7F8F9EF4B864e004FE86Ac8294E20950")},
	// 	[]common.Address{common.HexToAddress("0x65897aDCBa42dcCA5DD162c647b1cC3E31238490")},
	// )
	internalReserveAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")},
		[]common.Address{common.HexToAddress("0x2C5a182d280EeB5824377B98CD74871f78d6b8BC")},
	)
	// pricingContractAddress = app.NewAddress(
	// 	[]common.Address{common.HexToAddress("0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B")},
	// 	[]common.Address{common.HexToAddress("0xe3E415a7a6c287a95DC68a01ff036828073fD2e6")},
	// )
	proxyContractAddress = deployment.NewAddress(
		[]common.Address{common.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755")},
		[]common.Address{common.HexToAddress("0xC14f34233071543E979F6A79AA272b0AB1B4947D")},
	)
	// burnerContractAddress = app.NewAddress(
	// 	[]common.Address{common.HexToAddress("0x52166528FCC12681aF996e409Ee3a421a4e128A3")},
	// 	[]common.Address{common.HexToAddress("0x39682A7b8E4A03b2c8dC6DA6E0146Aee4E29A306")},
	// )

	// oldNetworkContractAddress = app.NewAddress(
	// 	[]common.Address{
	// 		common.HexToAddress("0x964F35fAe36d75B1e72770e244F6595B68508CF5"),
	// 		// production old internal network v2
	// 		common.HexToAddress("0x91a502C678605fbCe581eae053319747482276b9")},
	// 	[]common.Address{
	// 		common.HexToAddress("0xD2D21FdeF0D054D2864ce328cc56D1238d6b239e"),
	// 		// staging old internal network v2
	// 		common.HexToAddress("0x706aBcE058DB29eB36578c463cf295F180a1Fe9C")},
	// )

	// oldBurnerContractAddress = app.NewAddress(
	// 	[]common.Address{
	// 		common.HexToAddress("0x4E89bc8484B2c454f2F7B25b612b648c45e14A8e"),
	// 		common.HexToAddress("0x07f6e905f2a1559cd9fd43cb92f8a1062a3ca706"),
	// 		// old burner contract v2
	// 		common.HexToAddress("0xed4f53268bfdFF39B36E8786247bA3A02Cf34B04")},
	// 	[]common.Address{
	// 		common.HexToAddress("0xB2cB365D803Ad914e63EA49c95eC663715c2F673"),
	// 		// staging old burner contract v2
	// 		common.HexToAddress("0xd6703974Dc30155d768c058189A2936Cf7C62Da6")},
	// )
)
