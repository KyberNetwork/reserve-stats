package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/priceanalytics/common"
)

// Interface is the interface for price analytic storage
type Interface interface {
	UpdatePriceAnalytic(data common.PriceAnalytic) error
	GetPriceAnalytic(fromTime, toTime time.Time) ([]common.PriceAnalytic, error)
}
