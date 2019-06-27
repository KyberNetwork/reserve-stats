package postgrestorage

import (
	"fmt"
	"testing"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

const (
	driverName   = "postgres"
	testDataFile = "../testdata/postgresql/export.sql"
)

func newTestTradeLogPostgresql(dbName string) (*TradeLogDB, error) {
	var err error
	sugar := testutil.MustNewDevelopmentSugaredLogger()
	if err = dropDatabaseIfExists(dbName); err != nil {
		return nil, err
	}
	if err = createDatabase(dbName); err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, "reserve_stats", "reserve_stats", dbName,
	)
	db, err := sqlx.Connect(driverName, connStr)
	if err != nil {
		return nil, err
	}
	tokenAmountFormatter := blockchain.NewMockTokenAmountFormatter()
	storage, err := NewTradeLogDB(sugar, db, tokenAmountFormatter)
	if err != nil {
		return nil, err
	}
	sugar.Debugw("create test database object successful", "db_name", dbName)
	return storage, nil
}

func createDatabase(dbName string) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		"127.0.0.1", 5432, "reserve_stats", "reserve_stats",
	)
	db, err := sqlx.Connect(driverName, connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec(fmt.Sprintf(`CREATE DATABASE %v`, dbName)); err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}
	return nil
}

func (tldb *TradeLogDB) tearDown(dbName string) error {
	var err error
	if err = tldb.db.Close(); err != nil {
		return err
	}
	err = dropDatabaseIfExists(dbName)
	return err
}

func dropDatabaseIfExists(dbName string) error {
	const driverName = "postgres"
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		"127.0.0.1", 5432, "reserve_stats", "reserve_stats",
	)
	db, err := sqlx.Connect(driverName, connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err = db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %v`, dbName)); err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}
	return nil
}

func loadTestData(db *sqlx.DB, path string) error {
	_, err := sqlx.LoadFile(db, path)
	if err != nil {
		return err
	}
	_, err = db.Exec("SET search_path = public")
	return err
}

func TestTradeLogDB_LastBlock(t *testing.T) {
	const (
		dbName = "test_last_block"
	)

	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()

	blockNumber, err := testStorage.LastBlock()
	require.NoError(t, err, "call max block failed")
	t.Log(blockNumber)
	require.Equal(t, blockNumber, int64(0))
	require.NoError(t, loadTestData(testStorage.db, testDataFile))
	// call lastBlock after load test storage
	blockNumber, err = testStorage.LastBlock()
	require.NoError(t, err, "call max block failed after load data")
	require.NotEqual(t, blockNumber, int64(0))
	t.Logf("blocknumber: %v\n", blockNumber)
}

func TestSaveTradeLogs(t *testing.T) {
	const (
		dbName = "test_save_trade_log"
	)
	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()

	tradeLogs, err := utils.GetSampleTradeLogs("../testdata/trade_logs.json")
	require.NoError(t, err)
	require.NoError(t, testStorage.SaveTradeLogs(tradeLogs))
}

func TestTradeLogDB_LoadTradeLogs(t *testing.T) {
	const (
		fromTime = 1539000000000
		toTime   = 1539250666000
		dbName   = "test_load_trade_log"
	)

	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()

	require.NoError(t, loadTestData(testStorage.db, testDataFile))

	tradeLogs, err := testStorage.LoadTradeLogs(timeutil.TimestampMsToTime(fromTime), timeutil.TimestampMsToTime(toTime))
	require.NoError(t, err)
	t.Log(len(tradeLogs))
}
