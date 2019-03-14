package binance

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

const (
	defaultHardLimit = 1200 / 60
	defaultWafLimit  = 4000 / 60 / 5
	// defaultMaxWeight should be greater than max weight required by a request in Binance client.
	defaultMaxWeight = 5
)

// RateLimiter implements Limiter interface.
// Binance API has three kinds of limit (https://support.binance.com/hc/en-us/articles/360004492232-API-Frequently-Asked-Questions-FAQ-#%E2%80%9Climits_hardlimits%E2%80%9D):
// 1. Hard-Limits: 1,200 request weight per minute (https://api.binance.com/api/v1/exchangeInfo)
// 2. Machine Learning Limits
// 3. Web Application Firewall Limits: no detail is provided in documentation. @cTnko (Stefan) in Binance API
//    suggests to limit number of requests below 4000 / 5 minutes.
type RateLimiter struct {
	m sync.Mutex

	// hardLimiter limits the resource by weight of each request.
	// https://api.binance.com/api/v1/exchangeInfo
	hardLimiter *rate.Limiter
	// wafLimiter limits the resources by keeping requests count low enough to avoid CloudFront 403 error.
	// https://support.binance.com/hc/en-us/articles/360004492232-API-Frequently-Asked-Questions-FAQ-#%E2%80%9Cwafban%E2%80%9D
	wafLimiter *rate.Limiter
}

// NewRateLimiter returns a new instance of RateLimiter to use with Binance client.
func NewRateLimiter(hardLimit float64) *RateLimiter {
	hardLimiter := rate.NewLimiter(rate.Limit(hardLimit), defaultMaxWeight)

	// wafLimit will be guessed from configured hard limit
	wafLimit := hardLimit / defaultHardLimit * defaultWafLimit
	if wafLimit > defaultWafLimit {
		wafLimit = defaultWafLimit
	}

	wafLimiter := rate.NewLimiter(rate.Limit(wafLimit), 1)
	return &RateLimiter{
		m:           sync.Mutex{},
		hardLimiter: hardLimiter,
		wafLimiter:  wafLimiter,
	}
}

// WaitN waits until enough resources are available for a request with given weight.
func (r *RateLimiter) WaitN(ctx context.Context, n int) error {
	r.m.Lock()
	defer r.m.Unlock()
	errGr := errgroup.Group{}
	errGr.Go(func() error {
		if err := r.wafLimiter.Wait(ctx); err != nil {
			return err
		}
		return nil
	},
	)

	errGr.Go(func() error {
		if err := r.hardLimiter.WaitN(ctx, n); err != nil {
			return err
		}
		return nil
	},
	)

	return errGr.Wait()
}
