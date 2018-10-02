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

var (
	//ReserveAddr is the Kyber's own reserve address
	ReserveAddr = ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F")
)

const (
	//InfuraEndpoint is url for infura node
	InfuraEndpoint = "https://mainnet.infura.io"
)
