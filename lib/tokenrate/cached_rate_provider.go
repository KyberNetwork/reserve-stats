package tokenrate

import (
	"sync"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"go.uber.org/zap"
)

// CachedRateProviderOption is option for CachedRateProvider constructor
type CachedRateProviderOption func(*CachedRateProvider)

// WithTimeout is option to create CachedRateProvider with timeout
func WithTimeout(timeout time.Duration) CachedRateProviderOption {
	return func(crp *CachedRateProvider) {
		crp.timeout = timeout
	}
}

// WithWarningOnly is option to create CachedRateProvider with withwarn is true
// USDRate can return the last non-zero value with warning
func WithWarningOnly() CachedRateProviderOption {
	return func(crp *CachedRateProvider) {
		crp.warningOnly = true
	}
}

//NewCachedRateProvider return cached provider for eth usd rate
func NewCachedRateProvider(sugar *zap.SugaredLogger, provider tokenrate.ETHUSDRateProvider, options ...CachedRateProviderOption) *CachedRateProvider {
	crp := &CachedRateProvider{
		sugar:    sugar,
		provider: provider,
	}
	for _, option := range options {
		option(crp)
	}
	return crp
}

//CachedRateProvider is cached provider for eth usd rate
type CachedRateProvider struct {
	sugar       *zap.SugaredLogger
	timeout     time.Duration
	provider    tokenrate.ETHUSDRateProvider
	warningOnly bool

	mu         sync.RWMutex
	cachedRate float64
	cachedTime time.Time
}

func (crp *CachedRateProvider) isCacheExpired() bool {
	now := time.Now()
	expired := now.Sub(crp.cachedTime) > crp.timeout
	if expired {
		crp.sugar.Infow("cached rate is expired",
			"now", now,
			"cached_time", crp.cachedTime,
			"timeout", crp.timeout,
		)
	}
	return expired
}

//USDRate return usd rate
func (crp *CachedRateProvider) USDRate(timestamp time.Time) (float64, error) {
	logger := crp.sugar.With(
		"func", "users/http/cachedRateProvider.USDRate",
		"timestamp", timestamp,
	)

	crp.mu.RLock()
	cachedRate := crp.cachedRate
	cachedTime := crp.cachedTime
	crp.mu.RUnlock()
	if cachedRate == 0 || cachedTime.IsZero() || crp.isCacheExpired() {
		logger.Debug("cache miss, calling rate provider")
		rate, err := crp.provider.USDRate(timestamp)
		if err != nil && crp.warningOnly && cachedRate != 0 {
			logger.Warnw("failed to fetch USD rate, return last non-zero rate", "rate", cachedRate)
			return cachedRate, nil
		}
		if err != nil {
			return 0, err
		}

		logger.Debugw("got ETH/USD rate",
			"rate", rate)
		crp.mu.Lock()
		crp.cachedRate = rate
		crp.cachedTime = time.Now()
		crp.mu.Unlock()
		return rate, nil
	}

	logger.Debugw("cache hit",
		"rate", cachedRate)
	return cachedRate, nil
}

//Name return the original provider name
func (crp *CachedRateProvider) Name() string {
	return crp.provider.Name()
}
