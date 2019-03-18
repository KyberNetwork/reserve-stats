package fetcher

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/huobi"

	"go.uber.org/zap"
)

//Fetcher is the struct to fetch/store Data from huobi
type Fetcher struct {
	sugar      *zap.SugaredLogger
	client     huobi.Interface
	retryDelay time.Duration
	attempt    int
}

//NewFetcher return a fetcher object
func NewFetcher(sugar *zap.SugaredLogger, client huobi.Interface, retryDelay time.Duration, attempt int) *Fetcher {
	return &Fetcher{
		sugar:      sugar,
		client:     client,
		retryDelay: retryDelay,
		attempt:    attempt,
	}
}
