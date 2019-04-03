package common

import ac "github.com/KyberNetwork/reserve-stats/accounting/common"

type AllAddressesResponse struct {
	Version int64               `json:"version"`
	Data    []ac.ReserveAddress `json:"data"`
}
