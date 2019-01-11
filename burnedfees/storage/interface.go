package storage

import "github.com/KyberNetwork/reserve-stats/burnedfees/common"

// Interface is the database interaction of burned-fees-crawler service.
type Interface interface {
	Store([]common.BurnAssignedFeesEvent) error
	LastBlock() (int64, error)
}
