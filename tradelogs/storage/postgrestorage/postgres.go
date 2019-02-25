package postgrestorage

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// TradeLogDB is storage of tradelog data
type TradeLogDB struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

func (tldb *TradeLogDB) LastBlock() (int64, error) {
	return 0, nil
}

func (tldb *TradeLogDB) SaveTradeLogs(logs []common.TradeLog) error {
	return nil
}
