package common

import "github.com/ethereum/go-ethereum/common"

// AppObject an app
type AppObject struct {
	ID        int64            `json:"id,omitempty"`
	AppName   string           `json:"app_name" binding:"required"`
	Addresses []common.Address `json:"addresses" binding:"required"`
}
