package postgres

import (
	"fmt"
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/utils"
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
	storage, err := NewTradeLogDB(sugar, db, tokenAmountFormatter, blockchain.KNCAddr)
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
	t.Skip()
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
	t.Skip()
	const (
		dbName = "test_save_trade_log"
	)
	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()
	var result common.CrawlResult
	result.Trades, err = utils.GetSampleTradeLogs("../testdata/trade_logs.json")
	if err != nil {
		fmt.Println("error: ", err)
	}
	require.NoError(t, err)
	require.NoError(t, testStorage.SaveTradeLogs(&result))

	tls, err := testStorage.LoadTradeLogs(timeutil.TimestampMsToTime(1554353231000), timeutil.TimestampMsToTime(1554353231000))
	require.NoError(t, err)
	require.Equal(t, 1, len(tls))
}

func TestSaveTradeLogs_Overwrite(t *testing.T) {
	t.Skip()
	const (
		dbName = "test_save_trade_log"
	)
	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()
	var timestampMs uint64 = 1554336000000
	timestamp := timeutil.TimestampMsToTime(timestampMs)
	tradelog := common.Tradelog{
		Timestamp:       timestamp,
		BlockNumber:     uint64(6100010),
		TransactionHash: ethereum.HexToHash("0x33dcdbed63556a1d90b7e0f626bfaf20f6f532d2ae8bf24c22abb15c4e1fff01"),
		User: common.KyberUserInfo{
			UserAddress: ethereum.HexToAddress("0x85c5c26dc2af5546341fc1988b9d178148b4838b"),
		},
		TokenInfo: common.TradeTokenInfo{
			SrcAddress:  ethereum.HexToAddress("0x0f5d2fb29fb7d3cfee444a200298f468908cc942"),
			DestAddress: ethereum.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"),
		},
		SrcAmount:          big.NewInt(99995137653743),
		DestAmount:         big.NewInt(99995137653743773),
		USDTAmount:         big.NewInt(100000000000000000),
		OriginalUSDTAmount: big.NewInt(100000000000000000),
	}
	result := &common.CrawlResult{
		Trades: []common.Tradelog{tradelog},
	}
	require.NoError(t, testStorage.SaveTradeLogs(result))
	tradelog2 := tradelog
	tradelog2.USDTAmount = big.NewInt(0).Mul(big.NewInt(2), tradelog.USDTAmount)
	require.NoError(t, testStorage.SaveTradeLogs(result))

	tls, err := testStorage.LoadTradeLogs(timestamp, timestamp)
	require.NoError(t, err)
	require.Equal(t, len(tls), 1)
	assert.Equal(t, tradelog2.USDTAmount, tls[0].USDTAmount)
}

func TestTradeLogDB_LoadTradeLogs(t *testing.T) {
	t.Skip()
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
	for _, log := range tradeLogs {
		t.Logf("%+v %+v", log.OriginalUSDTAmount, log.USDTAmount)
	}
}

func TestTradeLogDB_LoadTradeLogsByTxHash(t *testing.T) {
	t.Skip()
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
	require.NotZero(t, len(tradeLogs))

	tradeLogsBS, err := testStorage.LoadTradeLogsByTxHash(tradeLogs[0].TransactionHash)
	require.NoError(t, err)
	require.NotZero(t, len(tradeLogsBS))
}

func TestTokenSymbol(t *testing.T) {
	t.Skip()
	const (
		dbName = "test_get_token_symbol"
	)
	var (
		ethAddress = ethereum.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee").Hex()
	)

	testStorage, err := newTestTradeLogPostgresql(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, testStorage.tearDown(dbName))
	}()
	require.NoError(t, loadTestData(testStorage.db, testDataFile))

	symbol, err := testStorage.GetTokenSymbol(ethAddress)
	assert.NoError(t, err)
	assert.Equal(t, "", symbol)

	err = testStorage.UpdateTokens([]string{ethAddress}, []string{"ETH"})
	assert.NoError(t, err)

	symbol, err = testStorage.GetTokenSymbol(ethAddress)
	assert.NoError(t, err)
	assert.Equal(t, "ETH", symbol)
}
