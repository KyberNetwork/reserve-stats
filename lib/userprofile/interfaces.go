package userprofile

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

//Interface define functionality of this package
type Interface interface {
	LookUpUserProfile(addr ethereum.Address) (UserProfile, error)
}
