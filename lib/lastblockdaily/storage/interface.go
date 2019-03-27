package storage

import (
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
)

//Interface define functionalities of lastBlockDaily Storage
type Interface interface {
	GetBlockInfo(time time.Time) (common.BlockInfo, error)
	UpdateBlockInfo(blockInfo common.BlockInfo) error
}
