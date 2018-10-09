package common

import (
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// TradeLog represent trade event on KyberNetwork
type TradeLog struct {
	Timestamp       time.Time
	BlockNumber     uint64
	TransactionHash ethereum.Hash
	Index           uint

	EtherReceivalSender ethereum.Address
	EtherReceivalAmount *big.Int

	UserAddress ethereum.Address
	SrcAddress  ethereum.Address
	DestAddress ethereum.Address
	SrcAmount   *big.Int
	DestAmount  *big.Int
	FiatAmount  float64

	ReserveAddress ethereum.Address
	WalletAddress  ethereum.Address
	WalletFee      *big.Int
	BurnFee        *big.Int
	IP             string
	Country        string
}

// AmountFormatter is the formatter used to format the amount from big number
// to float using preconfigured decimals. The intended implementation is from Core API client.
type AmountFormatter interface {
	FormatAmount(address ethereum.Address, amount *big.Int) (float64, error)
}
