package common

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

// ETHToken is the token object represent ETH in the system
var ETHToken = Token{
	ID:       "ETH",
	Name:     "Ethereum",
	Address:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	Decimals: 18,
	Active:   true,
	Internal: true,
}

var ( //WrapperAddrV1 is the Kyber's wrapper Address before block 5726056
	WrapperAddrV1 = ethereum.HexToAddress("0x533e6d1ffa2b96cf9c157475c76c38d1b13bc584")
	//WrapperAddrV2 is the Kyber's wrapper Address after block 5726056
	WrapperAddrV2 = ethereum.HexToAddress("0x6172AFC8c00c46E0D07ce3AF203828198194620a")
	//ReserveAddr is the Kyber's own reserve address
	ReserveAddr = ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")
)

const (
	//InfuraEndpoint: url for infura node
	InfuraEndpoint = "https://mainnet.infura.io"
	//StartingBlockV2 is the block where wrapper contract v2 is deployed and used.
	StartingBlockV2 = 5926056
)

// WrapperontractAddrs returns the proper network, contract addresses to use with given block number.
func WrapperContractAddr(block uint64) (wrapperAddr ethereum.Address) {
	if block < StartingBlockV2 {
		wrapperAddr = WrapperAddrV1
	} else {
		wrapperAddr = WrapperAddrV2
	}
	return
}
