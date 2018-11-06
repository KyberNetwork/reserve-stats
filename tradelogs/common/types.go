package common

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
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

	ETHUSDRate     float64 `json:"-"`
	ETHUSDProvider string  `json:"-"`
}

// MarshalJSON implements custom JSON marshaler for TradeLog to format timestamp in unix millis instead of RFC3339.
func (tl *TradeLog) MarshalJSON() ([]byte, error) {
	type AliasTradeLog TradeLog
	return json.Marshal(struct {
		Timestamp uint64 `json:"timestamp"`
		*AliasTradeLog
	}{
		AliasTradeLog: (*AliasTradeLog)(tl),
		Timestamp:     timeutil.TimeToTimestampMs(tl.Timestamp),
	})
}

// VolumeStats struct holds all the volume fields of volume in a specfic time
type VolumeStats struct {
	ETHAmount float64 `json:"eth_amount"`
	USDAmount float64 `json:"usd_amount"`
	Volume    float64 `json:"volume"`
}

// TradeSummary struct holds all the fields required for trade summary
type TradeSummary struct {
	ETHVolume       float64 `json:"total_eth_volume"`
	USDAmount       float64 `json:"total_usd_amount"`
	TotalBurnFee    float64 `json:"burn_fee"`
	TotalTrade      uint64  `json:"total_trade"`
	UniqueAddresses uint64  `json:"unique_addresses"`
	KYCEDAddresses  uint64  `json:"kyced_addresses"`
	NewUquAddresses uint64  `json:"new_unique_addresses"`
	USDPerTrade     float64 `json:"usd_per_trade"`
	ETHPerTrade     float64 `json:"eth_per_trade"`
}
