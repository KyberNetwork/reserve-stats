package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	blockInfoTable = "block_info"
)

//Option define init behaviour for db storage.
type Option func(*BlockInfoStorage) error

//WithBlockInfoTableName return Option to set trade table Name
func WithBlockInfoTableName(name string) Option {
	return func(hs *BlockInfoStorage) error {
		if hs.tableNames == nil {
			hs.tableNames = make(map[string]string)
		}
		hs.tableNames[blockInfoTable] = name
		return nil
	}
}

//BlockInfoStorage defines the object to store Huobi data
type BlockInfoStorage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	tableNames map[string]string
}

// NewDB return the BlockInfoStorage instance. User must call Close() before exit.
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB, options ...Option) (*BlockInfoStorage, error) {
	const schemaFMT = `
	CREATE TABLE IF NOT EXISTS %[1]s
(
	block bigint NOT NULL,
	time timestamp NOT NULL,
	CONSTRAINT %[1]s_pk PRIMARY KEY(block)
) ;
CREATE INDEX IF NOT EXISTS %[1]s_time_idx ON %[1]s (time);
`
	var (
		logger     = sugar.With("func", caller.GetCurrentFunctionName())
		tableNames = map[string]string{blockInfoTable: blockInfoTable}
	)
	hs := &BlockInfoStorage{
		sugar:      sugar,
		db:         db,
		tableNames: tableNames,
	}
	for _, opt := range options {
		if err := opt(hs); err != nil {
			return nil, err
		}
	}

	query := fmt.Sprintf(schemaFMT, hs.tableNames[blockInfoTable])
	logger.Debugw("initializing database schema", "query", query)
	if _, err := hs.db.Exec(query); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")
	return hs, nil
}

//Close close DB connection
func (bidb *BlockInfoStorage) Close() error {
	if bidb.db != nil {
		return bidb.db.Close()
	}
	return nil
}

//UpdateBlockInfo store the blockInfo
func (bidb *BlockInfoStorage) UpdateBlockInfo(blockInfo common.BlockInfo) error {
	var (
		logger = bidb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"Block", blockInfo.Block,
		)
	)

	const updateStmt = `INSERT INTO %[1]s(block, time)
	VALUES ( 
		$1,
		$2
	)
	ON CONFLICT ON CONSTRAINT %[1]s_pk DO NOTHING;`
	query := fmt.Sprintf(updateStmt,
		bidb.tableNames[blockInfoTable],
	)
	logger.Debugw("updating blockInfo...", "query", query)
	_, err := bidb.db.Exec(query, blockInfo.Block, blockInfo.Timestamp)
	return err
}

//GetBlockInfo return blockInfo on the day of input timestamp
//return sql.ErrNorows if that day hasn't got any data
func (bidb *BlockInfoStorage) GetBlockInfo(atTime time.Time) (common.BlockInfo, error) {
	var (
		logger = bidb.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"time", atTime.String(),
		)
		result common.BlockInfo
	)
	const selectStmt = `SELECT block, time FROM %[1]s WHERE time>$1 AND time<$2 Limit 1`
	query := fmt.Sprintf(selectStmt, bidb.tableNames[blockInfoTable])
	logger.Debugw("querying blockInfo...", "query", query)
	if err := bidb.db.Get(&result, query, timeutil.Midnight(atTime), timeutil.Midnight(atTime).AddDate(0, 0, 1)); err != nil {
		return common.BlockInfo{}, err
	}
	return result, nil
}
