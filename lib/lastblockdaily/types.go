package lastblockdaily

import (
	"time"
)

//BlockInfo is the struct contain block info
type BlockInfo struct {
	Block     uint64
	Timestamp time.Time
}
