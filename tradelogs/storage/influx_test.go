package storage

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

type mockAmountFormatter struct {
}

func (fmt *mockAmountFormatter) FormatAmount(address ethereum.Address, amount *big.Int) (float64, error) {
	return 0, nil
}

func createStorage() (*InfluxStorage, error) {
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

func getSampleTradeLogs(dataPath string) ([]common.TradeLog, error) {
	var tradeLogs []common.TradeLog
	byteValue, _ := os.Open(dataPath)
	err := json.NewDecoder(byteValue).Decode(&tradeLogs)
	if err != nil {
		return nil, err
	}
	return tradeLogs, nil
}

func TestSaveTradeLogs(t *testing.T) {
	tradeLogs, err := getSampleTradeLogs("testdata/trade_logs.json")
	storage, err := createStorage()
	if err != nil {
		t.Error("get unexpected error when create storage", "err", err.Error())
	}
	defer storage.RemoveDB()

	err = storage.SaveTradeLogs(tradeLogs)
	if err != nil {
		t.Error("get unexpected error when save trade logs", "err", err.Error())
	}
}
