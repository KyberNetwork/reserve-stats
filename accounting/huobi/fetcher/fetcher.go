package fetcher

import (
	"time"

	storage "github.com/KyberNetwork/reserve-stats/accounting/huobi/storage/withdrawal-history"
	"github.com/KyberNetwork/reserve-stats/lib/huobi"

	"go.uber.org/zap"
)

//Fetcher is the struct to fetch/store Data from huobi
type Fetcher struct {
	sugar      *zap.SugaredLogger
	client     huobi.Interface
	retryDelay time.Duration
	attempt    int
	db         storage.Interface
}

//NewFetcher return a fetcher object
func NewFetcher(sugar *zap.SugaredLogger, client huobi.Interface, retryDelay time.Duration, attempt int, strg ...storage.Interface) *Fetcher {
	fetcher := &Fetcher{
		sugar:      sugar,
		client:     client,
		retryDelay: retryDelay,
		attempt:    attempt,
	}
	if len(strg) != 0 {
		sugar.Info("fetcher is init with DB. Assigning...")
		fetcher.db = strg[0]
	}
	return fetcher
}
