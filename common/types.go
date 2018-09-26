package common

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// KNLog is the common interface of some important logging events.
type KNLog interface {
	TxHash() ethereum.Hash
	BlockNo() uint64
	Type() string
}

// TradeLog represent trade event on KyberNetwork
type TradeLog struct {
	Timestamp       uint64
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

func (self TradeLog) BlockNo() uint64       { return self.BlockNumber }
func (self TradeLog) Type() string          { return "TradeLog" }
func (self TradeLog) TxHash() ethereum.Hash { return self.TransactionHash }
