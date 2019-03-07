package listedtokenstorage

import "github.com/KyberNetwork/reserve-stats/accounting/common"

//Interface represent interface for accounting lsited token service
type Interface interface {
	CreateOrUpdate(tokens []common.ListedToken) error
}
