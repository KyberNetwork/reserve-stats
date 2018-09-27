package crawler

import "github.com/KyberNetwork/reserve-stats/common"

type Setting interface {
	GetInternalTokens() ([]common.Token, error)
	GetActiveTokens() ([]common.Token, error)
}
