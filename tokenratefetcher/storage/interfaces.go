package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/tokenratefetcher/common"
)

//Interface abstracts the implementation of storage functionality.
type Interface interface {
	LastTimePoint(providerName, tokenID, currencyID string) (time.Time, error)
	SaveRates(rates []common.TokenRate) error
}
