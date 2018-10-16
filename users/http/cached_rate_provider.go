package http

import (
	"sync"
	"time"

	"github.com/KyberNetwork/tokenrate"
	"go.uber.org/zap"
)

func newCachedRateProvider(sugar *zap.SugaredLogger, provider tokenrate.ETHUSDRateProvider, timeout time.Duration) *cachedRateProvider {
	return &cachedRateProvider{
		sugar:    sugar,
		timeout:  timeout,
		provider: provider,
	}
}

type cachedRateProvider struct {
	sugar    *zap.SugaredLogger
	timeout  time.Duration
	provider tokenrate.ETHUSDRateProvider

	mu         sync.RWMutex
	cachedRate float64
	cachedTime time.Time
}

func (crp *cachedRateProvider) isCacheExpired() bool {
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

func (crp *cachedRateProvider) USDRate(timestamp time.Time) (float64, error) {
	logger := crp.sugar.With(
		"func", "users/http/cachedRateProvider.USDRate",
		"timestamp", timestamp,
	)

	crp.mu.RLock()
	if crp.cachedRate == 0 || crp.cachedTime.IsZero() || crp.isCacheExpired() {
		crp.mu.RUnlock()

		logger.Debug("cache miss, calling rate provider")
		rate, err := crp.provider.USDRate(timestamp)
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
	crp.mu.RUnlock()

	logger.Debugw("cache hit",
		"rate", crp.cachedRate)
	return crp.cachedRate, nil
}

func (crp *cachedRateProvider) Name() string {
	return crp.provider.Name()
}
