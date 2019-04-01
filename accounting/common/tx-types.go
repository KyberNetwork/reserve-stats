package common

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
	ethereum "github.com/ethereum/go-ethereum/common"
	etherscan "github.com/nanmu42/etherscan-api"
)

// NormalTx holds info from normal tx query.
type NormalTx struct {
	BlockNumber int       `json:"blockNumber,string"`
	Timestamp   time.Time `json:"-"`
	Hash        string    `json:"hash"`
	BlockHash   string    `json:"blockHash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       *big.Int  `json:"value"`
	Gas         int       `json:"gas,string"`
	GasUsed     int       `json:"gasUsed,string"`
	GasPrice    *big.Int  `json:"gasPrice"`
	IsError     int       `json:"isError,string"`
}

// UnmarshalJSON is the custom unmarshaller that read timestamp in unix milliseconds.
func (tx *NormalTx) UnmarshalJSON(data []byte) error {
	type AliasNormalTx NormalTx
	decoded := new(struct {
		AliasNormalTx
		Timestamp uint64 `json:"timestamp"`
	})

	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	tx.BlockNumber = decoded.BlockNumber
	tx.Timestamp = timeutil.TimestampMsToTime(decoded.Timestamp).UTC()
	tx.Hash = decoded.Hash
	tx.BlockHash = decoded.BlockHash
	tx.From = decoded.From
	tx.To = decoded.To
	tx.Value = decoded.Value
	tx.Gas = decoded.Gas
	tx.GasUsed = decoded.GasUsed
	tx.GasPrice = decoded.GasPrice
	tx.IsError = decoded.IsError
	return nil
}

// MarshalJSON is the custom JSON marshaller that output timestamp in unix milliseconds.
func (tx NormalTx) MarshalJSON() ([]byte, error) {
	type AliasTNormalTx NormalTx
	return json.Marshal(struct {
		AliasTNormalTx
		Timestamp uint64 `json:"timestamp"`
	}{
		AliasTNormalTx: (AliasTNormalTx)(tx),
		Timestamp:      timeutil.TimeToTimestampMs(tx.Timestamp),
	})
}

// InternalTx holds info from internal tx query.
type InternalTx struct {
	BlockNumber int       `json:"blockNumber,string"`
	Timestamp   time.Time `json:"timestamp"`
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       *big.Int  `json:"value"`
	Gas         int       `json:"gas,string"`
	GasUsed     int       `json:"gasUsed,string"`
	IsError     int       `json:"isError,string"`
}

// UnmarshalJSON is the custom unmarshaller that read timestamp in unix milliseconds.
func (tx *InternalTx) UnmarshalJSON(data []byte) error {
	type AliasInternalTx InternalTx
	decoded := new(struct {
		AliasInternalTx
		Timestamp uint64 `json:"timestamp"`
	})
	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	tx.BlockNumber = decoded.BlockNumber
	tx.Timestamp = timeutil.TimestampMsToTime(decoded.Timestamp).UTC()
	tx.Hash = decoded.Hash
	tx.From = decoded.From
	tx.To = decoded.To
	tx.Value = decoded.Value
	tx.Gas = decoded.Gas
	tx.GasUsed = decoded.GasUsed
	tx.IsError = decoded.IsError
	return nil
}

// MarshalJSON is the custom JSON marshaller that output timestamp in unix milliseconds.
func (tx InternalTx) MarshalJSON() ([]byte, error) {
	type AliasInternalTx InternalTx
	return json.Marshal(struct {
		AliasInternalTx
		Timestamp uint64 `json:"timestamp"`
	}{
		AliasInternalTx: (AliasInternalTx)(tx),
		Timestamp:       timeutil.TimeToTimestampMs(tx.Timestamp),
	})
}

// ERC20Transfer holds info from ERC20 token transfer event query.
type ERC20Transfer struct {
	BlockNumber     int              `json:"blockNumber,string"`
	Timestamp       time.Time        `json:"-"`
	Hash            ethereum.Hash    `json:"-"`
	From            ethereum.Address `json:"-"`
	ContractAddress ethereum.Address `json:"-"`
	To              ethereum.Address `json:"-"`
	Value           *big.Int         `json:"value"`
	Gas             int              `json:"gas,string"`
	GasUsed         int              `json:"gasUsed,string"`
	GasPrice        *big.Int         `json:"gasPrice"`
}

//MarshalJSON return marshal form of erc20transfer
func (et ERC20Transfer) MarshalJSON() ([]byte, error) {
	type AliasErc20 ERC20Transfer
	var ts *uint64
	if !et.Timestamp.IsZero() {
		millis := timeutil.TimeToTimestampMs(et.Timestamp)
		ts = &millis
	}

	return json.Marshal(struct {
		Timestamp       *uint64 `json:"timestamp"`
		Hash            string  `json:"hash"`
		From            string  `json:"from"`
		ContractAddress string  `json:"contractAddress"`
		To              string  `json:"to"`
		AliasErc20
	}{
		Timestamp:       ts,
		Hash:            et.Hash.Hex(),
		From:            et.From.Hex(),
		ContractAddress: et.ContractAddress.Hex(),
		To:              et.To.Hex(),
		AliasErc20:      (AliasErc20)(et),
	})
}

//UnmarshalJSON return an ERC20Transfer object from JSON form
func (et *ERC20Transfer) UnmarshalJSON(data []byte) error {
	type AliasErc20 ERC20Transfer
	decoded := new(struct {
		Timestamp       *uint64 `json:"timestamp,omitempty"`
		Hash            string  `json:"hash"`
		From            string  `json:"from"`
		ContractAddress string  `json:"contractAddress"`
		To              string  `json:"to"`
		AliasErc20
	})
	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	if decoded.Timestamp != nil {
		et.Timestamp = timeutil.TimestampMsToTime(*decoded.Timestamp).UTC()
	}

	et.BlockNumber = decoded.BlockNumber

	et.Hash = ethereum.HexToHash(decoded.Hash)
	et.From = ethereum.HexToAddress(decoded.From)
	et.To = ethereum.HexToAddress(decoded.To)
	et.ContractAddress = ethereum.HexToAddress(decoded.ContractAddress)
	et.Value = decoded.Value
	et.Gas = decoded.Gas
	et.GasUsed = decoded.GasUsed
	et.GasPrice = decoded.GasPrice

	return nil
}

//EtherscanInternalTxToCommon transforms etherScan.InternalTx to accounting's InternalTx
func EtherscanInternalTxToCommon(tx etherscan.InternalTx) InternalTx {
	return InternalTx{
		BlockNumber: tx.BlockNumber,
		Timestamp:   tx.TimeStamp.Time(),
		Hash:        tx.Hash,
		From:        tx.From,
		To:          tx.To,
		Value:       tx.Value.Int(),
		Gas:         tx.Gas,
		GasUsed:     tx.GasUsed,
		IsError:     tx.IsError,
	}
}

//EtherscanERC20TransferToCommon transforms etherScan.ERC20Trasnfer to accounting's ERC20Transfer
func EtherscanERC20TransferToCommon(tx etherscan.ERC20Transfer) ERC20Transfer {
	return ERC20Transfer{
		BlockNumber:     tx.BlockNumber,
		Timestamp:       tx.TimeStamp.Time(),
		Hash:            ethereum.HexToHash(tx.Hash),
		From:            ethereum.HexToAddress(tx.From),
		ContractAddress: ethereum.HexToAddress(tx.ContractAddress),
		To:              ethereum.HexToAddress(tx.To),
		Value:           tx.Value.Int(),
		Gas:             tx.Gas,
		GasUsed:         tx.GasUsed,
		GasPrice:        tx.GasPrice.Int(),
	}
}

//EtherscanNormalTxToCommon transform etherScan.NormalTx to accounting's normalTx
func EtherscanNormalTxToCommon(tx etherscan.NormalTx) NormalTx {
	return NormalTx{
		BlockNumber: tx.BlockNumber,
		Timestamp:   tx.TimeStamp.Time(),
		Hash:        tx.Hash,
		BlockHash:   tx.BlockHash,
		From:        tx.From,
		To:          tx.To,
		Value:       tx.Value.Int(),
		Gas:         tx.Gas,
		GasUsed:     tx.GasUsed,
		GasPrice:    tx.GasPrice.Int(),
		IsError:     tx.IsError,
	}
}
