package lastblockdaily

import (
	"time"
)

//BlockInfo is the struct contain block info
type BlockInfo struct {
	Block     uint64    `db:"block"`
	Timestamp time.Time `db:"time"`
}
