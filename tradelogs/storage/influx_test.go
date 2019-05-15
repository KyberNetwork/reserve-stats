package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

var testStorage *InfluxStorage

func newTestInfluxStorage(db string) (*InfluxStorage, error) {
	sugar := testutil.MustNewDevelopmentSugaredLogger()

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://127.0.0.1:8086",
	})
	if err != nil {
		return nil, err
	}

	storage, err := NewInfluxStorage(
		sugar,
		db,
		influxClient,
		blockchain.NewMockTokenAmountFormatter(),
	)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

// tearDown remove the database that storing trade logs measurements.
func (is *InfluxStorage) tearDown() error {
	_, err := influxdb.QueryDB(is.influxClient, fmt.Sprintf("DROP DATABASE %s", is.dbName), is.dbName)
	return err
}

func getSampleTradeLogs(dataPath string) ([]common.TradeLog, error) {
	var tradeLogs []common.TradeLog
	byteValue, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(byteValue).Decode(&tradeLogs); err != nil {
		return nil, err
	}
	return tradeLogs, nil
}

func TestSaveTradeLogs(t *testing.T) {
	tradeLogs, err := getSampleTradeLogs("testdata/trade_logs.json")
	require.NoError(t, err)
	if err = testStorage.SaveTradeLogs(tradeLogs); err != nil {
		t.Error("get unexpected error when save trade logs", "err", err.Error())
	}
}

func TestSaveFirstTradeLogs(t *testing.T) {
	tradeLogs, err := getSampleTradeLogs("testdata/trade_logs.json")
	require.NoError(t, err)
	if err = testStorage.SaveTradeLogs(tradeLogs); err != nil {
		t.Error("get unexpected error when save trade logs", "err", err.Error())
	}
	testUser := ethereum.HexToAddress("0x85c5c26dc2af5546341fc1988b9d178148b4838b")
	traded, err := testStorage.userTraded(testUser)
	if err != nil {
		t.Error(err)
	}
	if !traded {
		t.Error("Expect user 0x85c5c26dc2af5546341fc1988b9d178148b4838b traded, but result suggests otherwise")
	}
}

func TestMain(m *testing.M) {
	var err error
	if testStorage, err = newTestInfluxStorage("test_log_db"); err != nil {
		log.Fatal("get unexpected error when create storage", "err", err.Error())
	}

	ret := m.Run()
	if err = testStorage.tearDown(); err != nil {
		log.Fatal(err)
	}
	os.Exit(ret)
}
