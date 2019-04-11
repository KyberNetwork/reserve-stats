package storage

import "github.com/KyberNetwork/reserve-stats/app-names/common"

// Interface is the common interface of app name storage implementations.
type Interface interface {
	CreateOrUpdate(app common.Application) (id int64, update bool, err error)
	Get(appID int64) (common.Application, error)
	GetAll(name *string, active *bool) ([]common.Application, error)
	Update(app common.Application) (err error)
	Delete(appID int64) (err error)
}
