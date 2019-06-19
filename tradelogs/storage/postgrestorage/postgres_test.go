package postgrestorage

import (
	"database/sql"
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgrestorage/schema"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

var testStorage *TradeLogDB

const (
	driverName   = "postgres"
	testDbName   = "test_log_db"
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

func TestNewTradeLogDB(t *testing.T) {
	blockNumber, err := testStorage.LastBlock()
	require.NoError(t, err, "call max block failed")
	t.Log(blockNumber)
	assert.Equal(t, blockNumber, int64(0))
	err = loadTestData(testStorage.db, testDataFile)
	require.NoError(t, err)
	// call lastBlock after load test storage
	blockNumber, err = testStorage.LastBlock()
	require.NoError(t, err, "call max block failed after load data")
	t.Log(blockNumber)
}

func TestTradeLogDB_LastBlock(t *testing.T) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"127.0.0.1", 5432, "reserve_stats", "reserve_stats", testDbName,
	)
	db, err := sqlx.Connect(driverName, connStr)
	require.NoError(t, err)
	stmt := fmt.Sprintf(`SELECT MAX("block_number") FROM "%v"`, schema.TradeLogsTableName)
	t.Log(stmt)
	var result sql.NullInt64
	row := db.QueryRow(stmt)
	err = row.Scan(&result)
	require.NoError(t, err)
	t.Log(result.Int64)
}

func TestSaveTradeLogs(t *testing.T) {
	tradeLogs, err := utils.GetSampleTradeLogs("../testdata/trade_logs.json")
	require.NoError(t, err)
	if err = testStorage.SaveTradeLogs(tradeLogs); err != nil {
		t.Error("get unexpected error when save trade logs", "err", err.Error())
	}
}

func TestMain(m *testing.M) {
	var err error
	if testStorage, err = newTestTradeLogPostgresql(testDbName); err != nil {
		log.Fatal("get unexpected error when create storage ", "error ", err.Error())
	}
	ret := m.Run()
	if err = testStorage.tearDown(testDbName); err != nil {
		log.Fatal(err)
	}
	os.Exit(ret)
}
