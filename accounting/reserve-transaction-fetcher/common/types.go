package common

import (
	"math/big"
	"time"
)

// NormalTx holds info from normal tx query.
type NormalTx struct {
	BlockNumber int       `json:"blockNumber,string"`
	TimeStamp   time.Time `json:"timeStamp"`
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

// InternalTx holds info from internal tx query.
type InternalTx struct {
	BlockNumber int       `json:"blockNumber,string"`
	TimeStamp   time.Time `json:"timeStamp"`
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       *big.Int  `json:"value"`
	Gas         int       `json:"gas,string"`
	GasUsed     int       `json:"gasUsed,string"`
	IsError     int       `json:"isError,string"`
}

// ERC20Transfer holds info from ERC20 token transfer event query.
type ERC20Transfer struct {
	BlockNumber     int       `json:"blockNumber,string"`
	TimeStamp       time.Time `json:"timeStamp"`
	Hash            string    `json:"hash"`
	From            string    `json:"from"`
	ContractAddress string    `json:"contractAddress"`
	To              string    `json:"to"`
	Value           *big.Int  `json:"value"`
	Gas             int       `json:"gas,string"`
	GasUsed         int       `json:"gasUsed,string"`
	GasPrice        *big.Int  `json:"gasPrice"`
}
