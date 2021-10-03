package withdrawalstorage

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/lib/pq"
)

//BinanceStorage is storage for binance fetcher including trade history and withdraw history
type BinanceStorage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//NewDB return a new instance of binance storage
func NewDB(sugar *zap.SugaredLogger, db *sqlx.DB) (*BinanceStorage, error) {
	var (
		logger = sugar.With("func", caller.GetCurrentFunctionName())
	)

	const schemaFmt = `CREATE TABLE IF NOT EXISTS "binance_withdrawals"
	(
	  id   text NOT NULL,
	  data JSONB,
	  CONSTRAINT binance_withdrawals_pk PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS binance_withdrawals_time_idx ON binance_withdrawals ((data ->> 'applyTime'));
	CREATE INDEX IF NOT EXISTS binance_withdrawals_txid_idx ON binance_withdrawals ((data ->> 'txId'));

	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS account TEXT;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS timestamp TIMESTAMP;

	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS amount FLOAT;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS address TEXT;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS asset TEXT;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS tx TEXT;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS apply_time BIGINT;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS status INTEGER;
	ALTER TABLE binance_withdrawals ADD COLUMN IF NOT EXISTS tx_fee FLOAT;
	`

	s := &BinanceStorage{
		sugar: sugar,
		db:    db,
	}

	logger.Debugw("create table query", "query", schemaFmt)

	if _, err := db.Exec(schemaFmt); err != nil {
		return nil, err
	}

	logger.Info("binance table init successfully")

	return s, nil
}

//Close database connection
func (bd *BinanceStorage) Close() error {
	if bd.db != nil {
		return bd.db.Close()
	}
	return nil
}

//UpdateWithdrawHistory save withdraw history to db
func (bd *BinanceStorage) UpdateWithdrawHistory(withdrawHistories []binance.WithdrawHistory, account string) (err error) {
	var (
		logger                      = bd.sugar.With("func", caller.GetCurrentFunctionName())
		withdrawJSON                []byte
		ids, assets, addresses, txs []string
		amounts, txFees             []float64
		statuses, applyTimes        []int64
		dataJSON                    [][]byte
		timestamps                  []time.Time
	)
	const updateQuery = `INSERT INTO binance_withdrawals (id, data, timestamp, account, amount, address, asset, tx, apply_time, status, tx_fee)
	VALUES(
		unnest($1::TEXT[]),
		unnest($2::JSONB[]),
		unnest($3::TIMESTAMP[]),
		$4,
		unnest($5::FLOAT[]),
		unnest($6::TEXT[]),
		unnest($7::TEXT[]),
		unnest($8::TEXT[]),
		unnest($9::BIGINT[]),
		unnest($10::INTEGER[]),
		unnest($11::FLOAT[])
	) ON CONFLICT ON CONSTRAINT binance_withdrawals_pk DO UPDATE SET data = EXCLUDED.data;
	`

	tx, err := bd.db.Beginx()
	if err != nil {
		return
	}

	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	logger.Debugw("query update withdraw history", "query", updateQuery)

	// prepare data to insert into db
	for _, withdraw := range withdrawHistories {
		withdrawJSON, err = json.Marshal(withdraw)
		if err != nil {
			return
		}
		ids = append(ids, withdraw.ID)
		timestamp, sErr := time.Parse("2006-01-02 15:04:05", withdraw.ApplyTime)
		if sErr != nil {
			return
		}
		timestamps = append(timestamps, timestamp)
		dataJSON = append(dataJSON, withdrawJSON)
		assets = append(assets, withdraw.Asset)
		addresses = append(addresses, withdraw.Address)
		txs = append(txs, withdraw.TxID)
		amount, err := strconv.ParseFloat(withdraw.Amount, 64)
		if err != nil {
			return err
		}
		amounts = append(amounts, amount)
		txFee, err := strconv.ParseFloat(withdraw.TxFee, 64)
		if err != nil {
			return err
		}
		txFees = append(txFees, txFee)
		statuses = append(statuses, withdraw.Status)
		applyTimes = append(applyTimes, timestamp.UnixMilli())
	}

	if _, err = tx.Exec(updateQuery, pq.Array(ids), pq.Array(dataJSON), pq.Array(timestamps), account, pq.Array(amounts), pq.Array(addresses), pq.Array(assets), pq.Array(txs), pq.Array(applyTimes), pq.Array(statuses), pq.Array(txFees)); err != nil {
		return err
	}

	return err
}

//WithdrawRecord represent a record of binace withdraw
type WithdrawRecord struct {
	Account    string          `db:"account"`
	Ids        pq.StringArray  `db:"ids"`
	Amounts    pq.Float64Array `db:"amounts"`
	Assets     pq.StringArray  `db:"assets"`
	Addresses  pq.StringArray  `db:"addresses"`
	ApplyTimes pq.Int64Array   `db:"apply_times"`
	Txs        pq.StringArray  `db:"txs"`
	Statuses   pq.Int64Array   `db:"statuses"`
	TxFees     pq.Float64Array `db:"tx_fees"`
}

//GetWithdrawHistory return list of withdraw fromTime to toTime
func (bd *BinanceStorage) GetWithdrawHistory(fromTime, toTime time.Time) (map[string][]binance.WithdrawHistory, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = make(map[string][]binance.WithdrawHistory)
		dbResult []WithdrawRecord
	)
	const selectStmt = `SELECT account, 
	ARRAY_AGG(id) as ids, 
	ARRAY_AGG(asset) as assets, 
	ARRAY_AGG(address) as addresses, 
	ARRAY_AGG(apply_time) as apply_times, 
	ARRAY_AGG(tx) as txs, 
	ARRAY_AGG(tx_fee) as tx_fees, 
	ARRAY_AGG(status) as statuses, 
	ARRAY_AGG(amount) as amounts FROM binance_withdrawals WHERE timestamp >=$1 AND timestamp <=$2 GROUP BY account;`

	logger.Debugw("querying trade history...", "query", selectStmt)

	if err := bd.db.Select(&dbResult, selectStmt, fromTime, toTime); err != nil {
		logger.Errorw("failed to get withdraw history", "error", err)
		return result, err
	}

	for _, record := range dbResult {
		arrResult := []binance.WithdrawHistory{}
		for index := range record.Ids {
			applyTime := time.UnixMilli(record.ApplyTimes[index])
			apTime := applyTime.Format("2006-01-02 15:04:05")
			arrResult = append(arrResult, binance.WithdrawHistory{
				ID:        record.Ids[index],
				Asset:     record.Assets[index],
				Amount:    strconv.FormatFloat(record.Amounts[index], 'f', -1, 64),
				Address:   record.Addresses[index],
				TxID:      record.Txs[index],
				ApplyTime: apTime,
				Status:    record.Statuses[index],
				TxFee:     strconv.FormatFloat(record.TxFees[index], 'f', -1, 64),
			})
		}
		result[record.Account] = arrResult
	}

	return result, nil
}

//GetLastStoredTimestamp return last timestamp stored in database
func (bd *BinanceStorage) GetLastStoredTimestamp(account string) (time.Time, error) {
	var (
		logger   = bd.sugar.With("func", caller.GetCurrentFunctionName())
		result   = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
		dbResult time.Time
		statuses = []string{strconv.Itoa(int(common.AwaitingApproval)), strconv.Itoa(int(common.Processing))}
	)
	const (
		selectStmt = `SELECT timestamp FROM binance_withdrawals WHERE account = $1 ORDER BY timestamp DESC LIMIT 1;`
		//handle not completed withdraw
		latestNotCompleted = `SELECT timestamp FROM binance_withdrawals WHERE data->>'status' = any($1) AND account = $2 ORDER BY timestamp ASC LIMIT 1;`
	)
	logger.Debugw("querying last stored timestamp", "query", selectStmt)

	err := bd.db.Get(&dbResult, latestNotCompleted, pq.Array(statuses), account)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}
	logger.Debugw("min processing record time", "time", dbResult)

	if dbResult.IsZero() {
		err := bd.db.Get(&dbResult, selectStmt, account)
		if err != nil && err != sql.ErrNoRows {
			return result, err
		}
	}
	if !dbResult.IsZero() {
		result = dbResult
	}

	return result, nil
}

//UpdateWithdrawHistoryWithFee update fee into withdraw history table
func (bd *BinanceStorage) UpdateWithdrawHistoryWithFee(withdrawHistories []binance.WithdrawHistory, account string) (err error) {
	var (
		logger = bd.sugar.With("func", caller.GetCurrentFunctionName())
	)

	const (
		updateStmt = `UPDATE binance_withdrawals SET data = $1 WHERE data->>'txId' = $2`
	)
	logger.Debugw("update withdraw history", "query", updateStmt)
	tx, err := bd.db.Beginx()
	if err != nil {
		return
	}
	defer pgsql.CommitOrRollback(tx, bd.sugar, &err)

	for _, withdraw := range withdrawHistories {
		withdrawJSON, err := json.Marshal(withdraw)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(updateStmt, withdrawJSON, withdraw.TxID); err != nil {
			return err
		}
	}

	return nil
}
