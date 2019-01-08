package tokenratefetcher

import (
	"time"
)

// TokenRate represent rate for usd
type TokenRate struct {
	Timestamp time.Time
	Rate      float64
	Provider  string
	TokenID   string
	Currency  string
}
