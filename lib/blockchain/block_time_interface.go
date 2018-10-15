package blockchain

import "time"

// BlockTimeResolverInterface define the functionality
type BlockTimeResolverInterface interface {
	Resolve(blockNumber uint64) (time.Time, error)
}
