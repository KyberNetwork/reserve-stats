package tokenrate

import "time"

// ETHUSDRate represent rate for usd
type ETHUSDRate struct {
	Timestamp   time.Time
	Rate        float64
	Provider    string
	BlockNumber uint64
}

// TokenRate represent rate for usd
type TokenRate struct {
	Timestamp time.Time
	Rate      float64
	Provider  string
	TokenID   string
	Currency  string
}
