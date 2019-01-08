package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/token-rate-fetcher/common"
)

//Interface abstracts the implementation of storage functionalities
type Interface interface {
	GetFirstTimePoint(providerName, tokenID, currencyID string) (time.Time, error)
	SaveRates(rates []common.TokenRate) error
}
