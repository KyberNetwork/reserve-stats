package postgrestorage

import (
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

func (tldb *TradeLogDB) GetTokenHeatmap(asset ethereum.Address, from, to time.Time, timezone int8) (map[string]common.Heatmap, error) {
	return nil, nil
}
