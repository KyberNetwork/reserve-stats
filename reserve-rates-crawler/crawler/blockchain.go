package crawler

import "time"

type blockTimeResolver interface {
	Resolve(blockNumber uint64) (time.Time, error)
}
