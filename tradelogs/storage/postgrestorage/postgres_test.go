package postgrestorage

import (
	"fmt"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
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
	driverName = "postgres"
)

func newTestTradeLogPostgresql(dbName string) (*TradeLogDB, error) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()

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

func TestNewTradeLogDB(t *testing.T) {
	blockNumber, err := testStorage.LastBlock()
	if err != nil {
		t.Fatal("call max block failed: ", err)
	}
	var assertion = assert.New(t)
	assertion.Equal(blockNumber, int64(0))
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
	var dbName = "test_log_db"
	if err = dropDatabaseIfExists(dbName); err != nil {
		log.Fatal("create database failed ", "error ", err.Error())
	}
	if err = createDatabase(dbName); err != nil {
		log.Fatal("create database failed ", "error ", err.Error())
	}
	if testStorage, err = newTestTradeLogPostgresql(dbName); err != nil {
		log.Fatal("get unexpected error when create storage ", "error ", err.Error())
	}
	ret := m.Run()
	if err = testStorage.tearDown(dbName); err != nil {
		log.Fatal(err)
	}
	os.Exit(ret)
}
