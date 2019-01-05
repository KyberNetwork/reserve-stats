package contracts

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	timeout               = 30 * time.Second
	blockNumberV1         = 5926000
	startingBlockV3 int64 = 6997111
)

// CachedContractAddressClient is a client to get contract address
type CachedContractAddressClient struct {
	mu                           *sync.RWMutex
	client                       bind.ContractBackend
	cachedInternalNetworkAddress map[int64]common.Address
	cachedPricingAddress         map[int64]common.Address
	cachedBurnerAddress          map[int64]common.Address
}

//NewCachedContractAddressClient return new cached contract address client
func NewCachedContractAddressClient(client bind.ContractBackend) *CachedContractAddressClient {
	return &CachedContractAddressClient{
		mu:                           &sync.RWMutex{},
		client:                       client,
		cachedInternalNetworkAddress: make(map[int64]common.Address),
		cachedPricingAddress:         make(map[int64]common.Address),
		cachedBurnerAddress:          make(map[int64]common.Address),
	}
}

//InternalNetworkContractAddress returns the address of internal network of all deployments
func (cc CachedContractAddressClient) InternalNetworkContractAddress(proxyAddress common.Address, blockNumber *big.Int) (common.Address, error) {
	if blockNumber.Cmp(big.NewInt(1)) == 0 {
		return cc.getInternalNetworkAddress(proxyAddress, 1)
	}
	if blockNumber.Cmp(big.NewInt(startingBlockV2)) == -1 {
		return cc.getInternalNetworkAddress(proxyAddress, blockNumberV1)
	}
	if blockNumber.Cmp(big.NewInt(startingBlockV2)) >= 0 {
		return cc.getInternalNetworkAddress(proxyAddress, startingBlockV2)
	}
	return cc.getInternalNetworkAddress(proxyAddress, startingBlockV3)
}

func (cc CachedContractAddressClient) getInternalNetworkAddress(proxyAddress common.Address, blockNumber int64) (common.Address, error) {
	if blockNumber != 1 {
		cc.mu.RLock()
		if _, exist := cc.cachedInternalNetworkAddress[blockNumber]; exist {
			cc.mu.RUnlock()
			return cc.cachedInternalNetworkAddress[blockNumber], nil
		}
		cc.mu.RUnlock()
	}

	proxyContract, err := NewProxy(proxyAddress, cc.client)
	if err != nil {
		return common.Address{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	callOpt := &bind.CallOpts{Context: ctx}
	// if not input blocknumber then get from latest block
	if blockNumber != 1 {
		callOpt.BlockNumber = big.NewInt(blockNumber)
	}

	address, err := proxyContract.KyberNetworkContract(callOpt)
	if err != nil {
		return common.Address{}, err
	}
	cc.mu.Lock()
	cc.cachedInternalNetworkAddress[blockNumber] = address
	cc.mu.Unlock()
	return address, nil
}

// PricingContractAddress returns the address of pricing contract of all deployments.
func (cc CachedContractAddressClient) PricingContractAddress(internalNetworkAddress common.Address, blockNumber *big.Int) (common.Address, error) {
	if blockNumber.Cmp(big.NewInt(1)) == 0 {
		return cc.getPricingContractAddress(internalNetworkAddress, 1)
	}
	if blockNumber.Cmp(big.NewInt(startingBlockV2)) == -1 {
		return cc.getPricingContractAddress(internalNetworkAddress, blockNumberV1)
	}
	if blockNumber.Cmp(big.NewInt(startingBlockV2)) >= 0 {
		return cc.getPricingContractAddress(internalNetworkAddress, startingBlockV2)
	}
	return cc.getPricingContractAddress(internalNetworkAddress, startingBlockV3)
}

func (cc *CachedContractAddressClient) getPricingContractAddress(internalNetworkAddress common.Address, blockNumber int64) (common.Address, error) {
	if blockNumber != 1 {
		cc.mu.RLock()
		if _, exist := cc.cachedPricingAddress[blockNumber]; exist {
			cc.mu.RUnlock()
			return cc.cachedPricingAddress[blockNumber], nil
		}
		cc.mu.RUnlock()
	}
	internalNetworkContract, err := NewInternalNetwork(internalNetworkAddress, cc.client)
	if err != nil {
		return common.Address{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// if not input blocknumber then get from latest block
	callOpt := &bind.CallOpts{Context: ctx}
	if blockNumber != 1 {
		callOpt.BlockNumber = big.NewInt(blockNumber)
	}
	address, err := internalNetworkContract.ExpectedRateContract(callOpt)
	if err != nil {
		return common.Address{}, err
	}
	cc.mu.Lock()
	cc.cachedPricingAddress[blockNumber] = address
	cc.mu.Unlock()
	return address, nil
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
func (cc CachedContractAddressClient) BurnerContractAddress(internalNetworkAddress common.Address, blockNumber *big.Int) (common.Address, error) {
	if blockNumber.Cmp(big.NewInt(1)) == 0 {
		return cc.getBurnerContractAddress(internalNetworkAddress, 1)
	}
	if blockNumber.Cmp(big.NewInt(startingBlockV2)) == -1 {
		return cc.getBurnerContractAddress(internalNetworkAddress, blockNumberV1)
	}
	if blockNumber.Cmp(big.NewInt(startingBlockV2)) >= 0 {
		return cc.getBurnerContractAddress(internalNetworkAddress, startingBlockV2)
	}
	return cc.getBurnerContractAddress(internalNetworkAddress, startingBlockV3)
}

func (cc CachedContractAddressClient) getBurnerContractAddress(internalNetworkAddress common.Address, blockNumber int64) (common.Address, error) {
	if blockNumber != 1 {
		cc.mu.RLock()
		if _, exist := cc.cachedBurnerAddress[blockNumber]; exist {
			cc.mu.RUnlock()
			return cc.cachedBurnerAddress[blockNumber], nil
		}
		cc.mu.RUnlock()
	}
	internalNetworkContract, err := NewInternalNetwork(internalNetworkAddress, cc.client)
	if err != nil {
		return common.Address{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// if not input blocknumber then get from latest block
	callOpt := &bind.CallOpts{Context: ctx}
	if blockNumber != 1 {
		callOpt.BlockNumber = big.NewInt(blockNumber)
	}
	address, err := internalNetworkContract.FeeBurnerContract(callOpt)
	cc.mu.Lock()
	cc.cachedBurnerAddress[blockNumber] = address
	cc.mu.Unlock()
	return address, nil
}

var (
	internalReserveAddress = app.NewAddress(
		[]common.Address{common.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")},
		[]common.Address{common.HexToAddress("0x2C5a182d280EeB5824377B98CD74871f78d6b8BC")},
	)
	proxyContractAddress = app.NewAddress(
		[]common.Address{common.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755")},
		[]common.Address{common.HexToAddress("0xC14f34233071543E979F6A79AA272b0AB1B4947D")},
	)
)
