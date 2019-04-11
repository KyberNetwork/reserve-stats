package storage

import (
	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/app-names/common"
)

// FilterConf is the configuration of GetAll function.
type FilterConf struct {
	Name    *string
	Address *string
	Active  *bool
}

// Filter is a filter of GetAll method.
type Filter func(*FilterConf)

// WithNameFilter filters the applications list by name.
func WithNameFilter(name string) Filter {
	return func(filters *FilterConf) {
		filters.Name = &name
	}
}

// WithAddressFilter filters the applications list by address.
func WithAddressFilter(address ethereum.Address) Filter {
	return func(filters *FilterConf) {
		addressFilter := address.Hex()
		filters.Address = &addressFilter
	}
}

// WithActiveFilter filters the applications to only returns active ones.
func WithActiveFilter() Filter {
	return func(filters *FilterConf) {
		active := true
		filters.Active = &active
	}
}

// WithInactiveFilter filters the applications to only returns inactive ones.
func WithInactiveFilter() Filter {
	return func(filters *FilterConf) {
		active := false
		filters.Active = &active
	}
}

// Interface is the common interface of app name storage implementations.
type Interface interface {
	CreateOrUpdate(app common.Application) (id int64, update bool, err error)
	Get(appID int64) (common.Application, error)
	GetAll(filters ...Filter) ([]common.Application, error)
	Update(app common.Application) (err error)
	Delete(appID int64) (err error)
}
