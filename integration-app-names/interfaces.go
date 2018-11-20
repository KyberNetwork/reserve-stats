package appname

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

//AddrToAppName define a set of interface required to translate address to app name
type AddrToAppName interface {
	LoadFromFile(path string) error
	UpdateMapAddrAppName(addr ethereum.Address, name string) error
	GetAddrToAppName() map[ethereum.Address]string
}
