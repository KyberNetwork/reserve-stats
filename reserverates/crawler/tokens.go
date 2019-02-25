package crawler

import (
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/ethereum/go-ethereum/common"
)

type reserveTokenFetcherInterface interface {
	Tokens(common.Address, uint64) ([]blockchain.TokenInfo, error)
}
