package postgres

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/burnedfees/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
)

const schema = `
	CREATE TABLE IF NOT EXISTS "burned_fee" (
		id SERIAL PRIMARY KEY,
		block_number INTEGER,
		timestamp TIMESTAMP,
		tx_hash TEXT,
		amount FLOAT,
		sender TEXT,
		reserve TEXT
	);
`

// Storage postgres storage for burned fee component
type Storage struct {
	l                    *zap.SugaredLogger
	db                   *sqlx.DB
	blockTimeResolver    blockchain.BlockTimeResolverInterface
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface
}

// NewPostgresStorage return new instance of postgres storage
func NewPostgresStorage(db *sqlx.DB, l *zap.SugaredLogger, blockTimeResolver blockchain.BlockTimeResolverInterface,
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface) (*Storage, error) {
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &Storage{
		db:                   db,
		l:                    l,
		blockTimeResolver:    blockTimeResolver,
		tokenAmountFormatter: tokenAmountFormatter,
	}, nil
}

// Store save burn fee event into db
func (ps *Storage) Store(records []common.BurnAssignedFeesEvent) error {
	var (
		logger = ps.l.With("func", caller.GetCurrentFunctionName())
		query  = `INSERT INTO "burned_fee" (block_number, timestamp, tx_hash, amount, sender, reserve)
		VALUES ($1, $2, $3, $4, $5, $6);`
	)
	logger.Infow("saving burn fee event into database", "query", query)
	for _, record := range records {
		qty, err := ps.tokenAmountFormatter.FromWei(blockchain.KNCAddr, record.Quantity)
		if err != nil {
			return err
		}
		ts, err := ps.blockTimeResolver.Resolve(record.BlockNumber)
		if err != nil {
			return err
		}
		if _, err := ps.db.Exec(query, record.BlockNumber, ts, record.TxHash.Hex(), qty, record.Sender.Hex(), record.Reserve.Hex()); err != nil {
			return err
		}
	}
	return nil
}

// LastBlock return last block saved in storage
func (ps *Storage) LastBlock() (int64, error) {
	var (
		logger    = ps.l.With("func", caller.GetCurrentFunctionName())
		query     = `SELECT coalesce(max(block_number), 0) FROM "burned_fee";`
		lastBlock int64
	)
	logger.Infow("get last block from db", "query", query)
	if err := ps.db.Get(&lastBlock, query); err != nil {
		return 0, err
	}
	return lastBlock, nil
}
