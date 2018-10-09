package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

var testStorage *InfluxStorage

type mockAmountFormatter struct {
}

func (fmt *mockAmountFormatter) FormatAmount(address ethereum.Address, amount *big.Int) (float64, error) {
	return 100, nil
}

func newTestInfluxStorage() (*InfluxStorage, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		return nil, err
	}

	storage, err := NewInfluxStorage(
		sugar,
		"test_db",
		influxClient,
		&mockAmountFormatter{},
	)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

// tearDown remove the database that storing trade logs measurements.
func (is *InfluxStorage) tearDown() error {
	_, err := is.queryDB(is.influxClient, fmt.Sprintf("DROP DATABASE %s", is.dbName))
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
	if err != nil {
		t.Fatal(err)
	}
	if err = testStorage.SaveTradeLogs(tradeLogs); err != nil {
		t.Error("get unexpected error when save trade logs", "err", err.Error())
	}

	// TODO: validate number of records inserted
}

func TestMain(m *testing.M) {
	var err error
	if testStorage, err = newTestInfluxStorage(); err != nil {
		log.Fatal("get unexpected error when create storage", "err", err.Error())
	}
	defer testStorage.tearDown()

	ret := m.Run()

	if err = testStorage.tearDown(); err != nil {
		log.Fatal(err)
	}

	os.Exit(ret)
}
