package crawler

import "github.com/KyberNetwork/reserve-stats/common"

// Setting defines a set of interface to query for the current reserve setting.
type Setting interface {
	// GetInternalTokens returns a set of internal token from reserve configuration or error if occurs
	GetInternalTokens() ([]common.Token, error)
	// GetActiveToken returns a set of active token from reserve configuration or error if occurs
	GetActiveTokens() ([]common.Token, error)
}
