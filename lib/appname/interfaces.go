package appname

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

// AddrToAppName define required function of an appname instance
type AddrToAppName interface {
	GetAddrToAppName() (map[ethereum.Address]string, error)
}
