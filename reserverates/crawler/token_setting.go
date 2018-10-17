package crawler

import "github.com/KyberNetwork/reserve-stats/lib/core"

// tokenSetting defines a set of interface to query for the current reserve setting.
type tokenSetting interface {
	// GetInternalTokens returns a set of internal token from reserve configuration or error if occurs
	GetInternalTokens() ([]core.Token, error)
	// GetActiveToken returns a set of active token from reserve configuration or error if occurs
	GetActiveTokens() ([]core.Token, error)
}
