package tokenrate

import "time"

// Provider is the common interface to query historical rates of any
// token to real worldp currencies.
// **Experimental**: the token, each provider intepretered currency,
// parameters differently.
type Provider interface {
	Rate(token, currency string, timestamp time.Time) (float64, error)
}

// ETHUSDRateProvider is the common interface to query historical
// rates of ETH to USD.
type ETHUSDRateProvider interface {
	USDRate(time.Time) (float64, error)
}
