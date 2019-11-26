package tokenrate

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/tokenrate"
)

const defaultExpire = time.Hour
const todayDefaultExpire = time.Minute * 5

// CachedRateProviderOption is option for CachedRateProvider constructor
type CachedRateProviderOption func(*CachedRateProvider)

// WithExpires is option to create CachedRateProvider with expiresDuration
func WithExpires(timeout time.Duration) CachedRateProviderOption {
	return func(crp *CachedRateProvider) {
		crp.expiresDuration = timeout
	}
}

// WithTodayExpires is option to create CachedRateProvider with todayExpiresDuration
func WithTodayExpires(exp time.Duration) CachedRateProviderOption {
	return func(crp *CachedRateProvider) {
		crp.todayExpiresDuration = exp
	}
}

//NewCachedRateProvider return cached provider for eth usd rate
func NewCachedRateProvider(sugar *zap.SugaredLogger, provider tokenrate.ETHUSDRateProvider, options ...CachedRateProviderOption) *CachedRateProvider {
	crp := &CachedRateProvider{
		sugar:                sugar,
		expiresDuration:      defaultExpire,
		todayExpiresDuration: todayDefaultExpire,
		provider:             provider,
	}
	for _, option := range options {
		option(crp)
	}
	crp.cache = gocache.New(crp.expiresDuration, crp.expiresDuration)
	return crp
}

//CachedRateProvider is cached provider for eth usd rate
type CachedRateProvider struct {
	sugar                *zap.SugaredLogger
	expiresDuration      time.Duration
	todayExpiresDuration time.Duration
	provider             tokenrate.ETHUSDRateProvider
	cache                *gocache.Cache
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
	cacheKey := timestamp.UTC().Format("2006-01-02")
	if item, found := crp.cache.Get(cacheKey); !found {
		logger.Debug("cache miss, calling rate provider")
		if rate, err = crp.provider.USDRate(timestamp); err != nil {
			return 0, err
		}
		logger.Debugw("got ETH/USD rate", "rate", rate)
		today := time.Now().UTC().Format("2006-01-02")
		if cacheKey == today {
			crp.cache.Set(cacheKey, rate, crp.todayExpiresDuration)
		} else {
			crp.cache.Set(cacheKey, rate, crp.expiresDuration)
		}
	} else {
		rate, _ = item.(float64)
	}
	return rate, nil
}

//Name return the original provider name
func (crp *CachedRateProvider) Name() string {
	return crp.provider.Name()
}
