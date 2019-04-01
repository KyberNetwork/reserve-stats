package client

import "github.com/KyberNetwork/reserve-stats/accounting/common"

//Interface define the functionality of reserve addresses client
type Interface interface {
	ReserveAddresses(...common.AddressType) ([]common.ReserveAddress, error)
}
