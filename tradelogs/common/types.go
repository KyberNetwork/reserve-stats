package common

import (
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// BurnFee represent burnFee event on KyberNetwork
type BurnFee struct {
	ReserveAddress ethereum.Address `json:"reserve_addr"`
	Amount         *big.Int         `json:"amount"`
}

// WalletFee represent feeToWallet event on KyberNetwork
type WalletFee struct {
	ReserveAddress ethereum.Address `json:"reserve_addr"`
	WalletAddress  ethereum.Address `json:"wallet_addr"`
	Amount         *big.Int         `json:"amount"`
}

// TradeLog represent trade event on KyberNetwork
type TradeLog struct {
	Timestamp       time.Time     `json:"timestamp"`
	BlockNumber     uint64        `json:"block_number"`
	TransactionHash ethereum.Hash `json:"tx_hash"`

	EtherReceivalSender ethereum.Address `json:"eth_receival_sender"`
	EtherReceivalAmount *big.Int         `json:"eth_receival_amount"`

	UserAddress ethereum.Address `json:"user_addr"`
	SrcAddress  ethereum.Address `json:"src_addr"`
	DestAddress ethereum.Address `json:"dst_addr"`
	SrcAmount   *big.Int         `json:"src_amount"`
	DestAmount  *big.Int         `json:"dst_amount"`
	FiatAmount  float64          `json:"fiat_amount"`

	BurnFees   []BurnFee   `json:"burn_fees"`
	WalletFees []WalletFee `json:"wallet_fees"`

	IP      string `json:"ip"`
	Country string `json:"country"`
}
