package binancestorage

import (
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	binanceTradeTable = "binance_trades"
)

//BinanceStorage is storage for binance fetcher including trade history and withdraw history
type BinanceStorage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) *BinanceStorage {
	var (
		logger = sugar.With("func", "accounting/binance-storage/binancestorage.NewDB")
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "%s"
	(
	  id   SERIAL PRIMARY KEY,
	  data JSONB,
	);
	`

	logger.Debugw("schema init db")
	return &BinanceStorage{
		sugar: sugar,
		db:    db,
	}
}

//UpdateTradeHistory save trade history into a postgres db
func (bd *BinanceStorage) UpdateTradeHistory(map[string]binance.TradeHistory) error {
	return nil
}
