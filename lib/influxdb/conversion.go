package influxdb

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

// GetAddressFromInterface convert given value to ethereum address
func GetAddressFromInterface(value interface{}) (common.Address, error) {
	var address common.Address
	s, ok := value.(string)
	if !ok {
		return address, errors.New("invalid address value")
	}
	return common.HexToAddress(s), nil
}

// GetTxHashFromInterface convert given value to ethereum tx hash
func GetTxHashFromInterface(value interface{}) (common.Hash, error) {
	var txHash common.Hash
	s, ok := value.(string)
	if !ok {
		return txHash, errors.New("invalid hash value")
	}
	return common.HexToHash(s), nil
}

// GetFloat64FromInterface convert given value to float64
func GetFloat64FromInterface(value interface{}) (float64, error) {
	number, ok := value.(json.Number)
	if !ok {
		return 0, errors.New("invalid number value")
	}
	return number.Float64()
}

// GetTimeFromInterface convert given value to time
func GetTimeFromInterface(value interface{}) (time.Time, error) {
	var result time.Time

	s, ok := value.(string)
	if !ok {
		return result, errors.New("invalid time value")
	}

	result, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return result, err
	}

	return result, nil
}
