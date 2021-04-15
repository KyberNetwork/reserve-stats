package common

import (
	"encoding/json"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

// IsETHAddress return true if address is eth
func IsETHAddress(addr ethereum.Address) bool {
	return addr == ethereum.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee")
}

// CrawlResult is result of the crawl
type CrawlResult struct {
	Reserves []Reserve  `json:"reserves"` // reserve update on this
	Trades   []Tradelog `json:"trades"`
}

// TradelogV4 is object for tradelog after katalyst upgrade
type Tradelog struct {
	Timestamp       time.Time     `json:"timestamp"`
	BlockNumber     uint64        `json:"block_number"`
	TransactionHash ethereum.Hash `json:"tx_hash"`

	TokenInfo TradeTokenInfo `json:"token_info"`
	// support version before katalyst
	SrcReserveAddress ethereum.Address `json:"-"`
	DstReserveAddress ethereum.Address `json:"-"`

	USDTAmount         *big.Int         `json:"usdt_amount"`
	OriginalUSDTAmount *big.Int         `json:"original_usdt_amount"`
	SrcAmount          *big.Int         `json:"src_amount"`
	DestAmount         *big.Int         `json:"dst_amount"`
	FiatAmount         float64          `json:"fiat_amount"`
	ReserveAddress     ethereum.Address `json:"reserve_address"`

	User            KyberUserInfo    `json:"user"`
	ReceiverAddress ethereum.Address `json:"receiver_address"`
	TxDetail        TxDetail         `json:"tx_detail"`

	Index   uint `json:"index"`
	Version uint `json:"version"`
}

// TradeSplit split the trade
type TradeSplit struct {
	ReserveAddress ethereum.Address `json:"reserve_addr"`
	SrcToken       ethereum.Address `json:"src_token"`
	DstToken       ethereum.Address `json:"dst_token"`
	SrcAmount      *big.Int         `json:"src_amount"`
	DstAmount      *big.Int         `json:"dst_amount"`
	Rate           *big.Int         `json:"rate"`
	Index          uint             `json:"index"`
}

// Reserve represent a reserve in KN
type Reserve struct {
	Address     ethereum.Address `json:"reserve"`
	ReserveID   [32]byte         `json:"reserve_id"`
	ReserveType uint64           `json:"reserve_type"`
	BlockNumber uint64           `json:"block_number"` // block number where reserve value (address, rebate_wallet) is applied
}

// TradeTokenInfo is token info
type TradeTokenInfo struct {
	SrcAddress  ethereum.Address `json:"src_addr"`
	SrcSymbol   string           `json:"src_symbol,omitempty"`
	DestAddress ethereum.Address `json:"dst_addr"`
	DestSymbol  string           `json:"dst_symbol,omitempty"`
}

// TxDetail about the tx fee
type TxDetail struct {
	GasUsed        uint64           `json:"gas_used"`
	GasPrice       *big.Int         `json:"gas_price"`
	TransactionFee *big.Int         `json:"transaction_fee"`
	TxSender       ethereum.Address `json:"tx_sender"`
}

// KyberUserInfo if available from KS
type KyberUserInfo struct {
	UserAddress ethereum.Address `json:"user_addr"`
	// UserName    string           `json:"user_name"`
	// ProfileID   int64            `json:"profile_id"`
	// Index       uint             `json:"index"` // the index of event log in transaction receipt
	// IP          string           `json:"ip"`
	// Country     string           `json:"country"`
	// UID         string           `json:"uid"`
}

// MarshalJSON implements custom JSON marshaller for TradeLog to format timestamp in unix millis instead of RFC3339.
func (tl *Tradelog) MarshalJSON() ([]byte, error) {
	type AliasTradeLog Tradelog
	return json.Marshal(struct {
		Timestamp uint64 `json:"timestamp"`
		*AliasTradeLog
	}{
		AliasTradeLog: (*AliasTradeLog)(tl),
		Timestamp:     timeutil.TimeToTimestampMs(tl.Timestamp),
	})
}

// UnmarshalJSON implements custom JSON unmarshal for TradeLog
func (tl *Tradelog) UnmarshalJSON(b []byte) error {
	type AliasTradeLog Tradelog
	type mask struct {
		Timestamp uint64 `json:"timestamp"`
		*AliasTradeLog
	}
	m := mask{
		Timestamp:     0,
		AliasTradeLog: (*AliasTradeLog)(tl),
	}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	tl.Timestamp = timeutil.TimestampMsToTime(m.Timestamp)
	return nil
}

// VolumeStats struct holds all the volume fields of volume in a specfic time
type VolumeStats struct {
	USDAmount float64 `json:"usd_amount"`
	Volume    float64 `json:"volume"`
}

//UserVolume represent volume of an user from time to time
type UserVolume struct {
	USDAmount float64 `json:"usd_amount"`
}

// UserInfo represent trade stats of an address
type UserInfo struct {
	Addr      string  `json:"user_address" db:"user_address"`
	ETHVolume float64 `json:"total_eth_volume" db:"total_eth_volume"`
	USDVolume float64 `json:"total_usd_volume" db:"total_usd_volume"`
}

//UserList - list of user
type UserList []UserInfo

//Len length of user list for sorting function
func (u UserList) Len() int {
	return len(u)
}

//Swap swap 2 item of user list
func (u UserList) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

//Less for sorting function
func (u UserList) Less(i, j int) bool {
	return u[i].ETHVolume < u[j].ETHVolume
}

// IntegrationVolume represent kyberSwap and non kyberswap volume
type IntegrationVolume struct {
	KyberSwapVolume    float64 `json:"kyber_swap_volume"`
	NonKyberSwapVolume float64 `json:"non_kyber_swap_volume"`
}

// StatsResponse reponse for stats
type StatsResponse struct {
	ETHVolume        float64 `json:"eth_volume"`
	USDVolume        float64 `json:"usd_volume"`
	UniqueAddresses  uint64  `json:"unique_addresses"`
	NewAdresses      uint64  `json:"new_addresses"`
	TotalTrades      uint64  `json:"total_trades"`
	FeeCollected     float64 `json:"fee_collected"`
	AverageTradeSize float64 `json:"average_size"`
}

// TopTokens by volume
// map token symbol with its volume
type TopTokens map[string]float64

// TopIntegrations by volume
// map integration name and its volume
type TopIntegrations map[string]float64

// TopReserves by volume
// map reserve name and its volume
type TopReserves map[string]float64
