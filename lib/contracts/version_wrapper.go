package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type VersionedWrapper struct {
	WrapperContractV1 *Wrapper
	WrapperContractV2 *Wrapper
}

var (
	//wrapperAddrV1 is the Kyber's wrapper Address before block 5726056
	wrapperAddrV1 = ethereum.HexToAddress("0x533e6d1ffa2b96cf9c157475c76c38d1b13bc584")
	//wrapperAddrV2 is the Kyber's wrapper Address after block 5726056
	wrapperAddrV2 = ethereum.HexToAddress("0x6172AFC8c00c46E0D07ce3AF203828198194620a")
)

const (
	//startingBlockV2 is the block where wrapper contract v2 is deployed and used.
	startingBlockV2 = 5926056
)

func NewVersionedWrapper(client bind.ContractBackend) (*VersionedWrapper, error) {
	wrapperContractV1, err := NewWrapper(wrapperAddrV1, client)
	if err != nil {
		return nil, err
	}
	wrapperContractV2, err := NewWrapper(wrapperAddrV2, client)
	if err != nil {
		return nil, err
	}
	return &VersionedWrapper{
		WrapperContractV1: wrapperContractV1,
		WrapperContractV2: wrapperContractV2,
	}, nil
}

func (vw *VersionedWrapper) GetReserveRate(block uint64, rsvAddr ethereum.Address, srcs, dest []ethereum.Address) ([]*big.Int, []*big.Int, error) {
	if block == 0 {
		//Latest rate from latest block at V2 contract
		return vw.WrapperContractV2.GetReserveRate(&bind.CallOpts{}, rsvAddr, srcs, dest)
	} else if block > startingBlockV2 {
		//V2 contract, call at specific block
		return vw.WrapperContractV2.GetReserveRate(&bind.CallOpts{BlockNumber: big.NewInt(int64(block))}, rsvAddr, srcs, dest)
	}
	//default case: V1 contract.
	return vw.WrapperContractV1.GetReserveRate(&bind.CallOpts{BlockNumber: big.NewInt(int64(block))}, rsvAddr, srcs, dest)
}
