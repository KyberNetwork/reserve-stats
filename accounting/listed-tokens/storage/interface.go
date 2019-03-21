package storage

import "github.com/KyberNetwork/reserve-stats/accounting/common"

//Interface represent interface for accounting lsited token service
type Interface interface {
	CreateOrUpdate(tokens map[string]common.ListedToken) error
	GetTokens() (map[string]common.ListedToken, error)
}
