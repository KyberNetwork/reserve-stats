package common

import (
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// BurnFee represent burn event on KyberNetwork
type BurnFee struct {
	ReserveAddress ethereum.Address
	Amount         *big.Int
}

// FeeToWallet represent feeToWallet on KyberNetwork
type FeeToWallet struct {
	ReserveAddress ethereum.Address
	WalletAddress  ethereum.Address
	Amount         *big.Int
}

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

	BurnFees     []BurnFee
	FeeToWallets []FeeToWallet

	IP      string
	Country string
}
