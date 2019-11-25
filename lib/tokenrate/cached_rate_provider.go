package tokenrate

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/tokenrate"
)

const defaultTimeout = time.Minute * 5

// CachedRateProviderOption is option for CachedRateProvider constructor
type CachedRateProviderOption func(*CachedRateProvider)

// WithTimeout is option to create CachedRateProvider with timeout
func WithTimeout(timeout time.Duration) CachedRateProviderOption {
	return func(crp *CachedRateProvider) {
		crp.timeout = timeout
	}
}

//NewCachedRateProvider return cached provider for eth usd rate
func NewCachedRateProvider(sugar *zap.SugaredLogger, provider tokenrate.ETHUSDRateProvider, options ...CachedRateProviderOption) *CachedRateProvider {
	crp := &CachedRateProvider{
		sugar:    sugar,
		timeout:  defaultTimeout,
		provider: provider,
	}
	for _, option := range options {
		option(crp)
	}
	crp.cache = gocache.New(crp.timeout, crp.timeout)
	return crp
}

//CachedRateProvider is cached provider for eth usd rate
type CachedRateProvider struct {
	sugar    *zap.SugaredLogger
	timeout  time.Duration
	provider tokenrate.ETHUSDRateProvider
	cache    *gocache.Cache
}

//USDRate return usd rate
func (crp *CachedRateProvider) USDRate(timestamp time.Time) (float64, error) {
	var (
		rate   float64
		err    error
		logger = crp.sugar.With("func", caller.GetCurrentFunctionName(),
			"timestamp", timestamp,
		)
	)
	if item, found := crp.cache.Get(timestamp.Format("2006-01-02")); !found {
		logger.Debug("cache miss, calling rate provider")
		if rate, err = crp.provider.USDRate(timestamp); err != nil {
			return 0, err
		}
		logger.Debugw("got ETH/USD rate", "rate", rate)
		crp.cache.Set(timestamp.Format("2006-01-02"), rate, crp.timeout)
	} else {
		rate, _ = item.(float64)
	}
	return rate, nil
}

//Name return the original provider name
func (crp *CachedRateProvider) Name() string {
	return crp.provider.Name()
}
