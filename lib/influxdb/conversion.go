package influxdb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// GetAddressFromInterface converts given value to ethereum address.
func GetAddressFromInterface(value interface{}) (common.Address, error) {
	var address common.Address
	if value == nil {
		return common.Address{}, nil
	}
	s, ok := value.(string)
	if !ok {
		return address, fmt.Errorf("invalid address value %v", value)
	}
	return common.HexToAddress(s), nil
}

// GetTxHashFromInterface converts given value to ethereum tx hash.
func GetTxHashFromInterface(value interface{}) (common.Hash, error) {
	var txHash common.Hash
	s, ok := value.(string)
	if !ok {
		return txHash, fmt.Errorf("invalid hash value %v", value)
	}
	return common.HexToHash(s), nil
}

// GetFloat64FromInterface converts given value to float64.
func GetFloat64FromInterface(value interface{}) (float64, error) {
	if value == nil {
		return 0, nil
	}
	number, ok := value.(json.Number)
	if !ok {
		return 0, fmt.Errorf("invalid number value %v", value)
	}
	return number.Float64()
}

// GetInt64FromInterface converts given value to int64.
// Since influx doesn't support fill without group by queries, any nil interface will be considered as 0
func GetInt64FromInterface(value interface{}) (int64, error) {
	if value == nil {
		return 0, nil
	}
	number, ok := value.(json.Number)
	if !ok {
		return 0, fmt.Errorf("invalid int64 value %v", value)
	}
	return number.Int64()
}

// GetInt64FromTagValue converts given value to int64
// The original tag value should be string
func GetInt64FromTagValue(value interface{}) (int64, error) {
	if value == nil {
		return 0, nil
	}
	number, ok := value.(string)
	if !ok {
		return 0, fmt.Errorf("invalid int64 tag value %v", value)
	}
	return strconv.ParseInt(number, 10, 64)
}

// GetUint64FromTagValue converts given value to uint64
// The original tag value should be string
func GetUint64FromTagValue(value interface{}) (uint64, error) {
	number, ok := value.(string)
	if !ok {
		return 0, fmt.Errorf("invalid uint64 tag value %v", value)
	}
	return strconv.ParseUint(number, 10, 64)
}

// GetTimeFromInterface converts given value to time.
func GetTimeFromInterface(value interface{}) (time.Time, error) {
	var result time.Time

	s, ok := value.(string)
	if !ok {
		return result, fmt.Errorf("invalid time value %v", value)
	}

	result, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return result, err
	}

	return result, nil
}

// GetUint64FromInterface converts given value to int64.
func GetUint64FromInterface(value interface{}) (uint64, error) {
	number, err := GetInt64FromInterface(value)
	if err != nil {
		return 0, err
	}
	return uint64(number), nil
}
