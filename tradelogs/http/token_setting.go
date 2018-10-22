package http

import "github.com/KyberNetwork/reserve-stats/lib/core"

// coreSetting defines a set of interface to query for the current reserve setting.
type coreSetting interface {
	// GetActiveToken returns a set of active token from reserve configuration or error if occurs
	GetActiveTokens() ([]core.Token, error)
}
