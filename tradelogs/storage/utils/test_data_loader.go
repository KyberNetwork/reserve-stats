package utils

import (
	"encoding/json"
	"os"

	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// GetSampleTradeLogs get sample a tradelogs
func GetSampleTradeLogs(dataPath string) ([]common.TradelogV4, error) {
	var tradeLogs []common.TradelogV4
	byteValue, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}
	if err = json.NewDecoder(byteValue).Decode(&tradeLogs); err != nil {
		return nil, err
	}
	return tradeLogs, nil
}
