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

// BurnFee represent burnFee event on KyberNetwork
type BurnFee struct {
	ReserveAddress ethereum.Address `json:"reserve_addr"`
	Amount         *big.Int         `json:"amount"`
	Index          uint             `json:"index"` // the index of event log in transaction receipt
}

// WalletFee represent feeToWallet event on KyberNetwork
type WalletFee struct {
	ReserveAddress ethereum.Address `json:"reserve_addr"`
	WalletAddress  ethereum.Address `json:"wallet_addr"`
	WalletName     string           `json:"wallet_name"`
	Amount         *big.Int         `json:"amount"`
	Index          uint             `json:"index"` // the index of event log in transaction receipt
}

// TradeLog represent trade event on KyberNetwork
type TradeLog struct {
	Timestamp       time.Time     `json:"timestamp"`
	BlockNumber     uint64        `json:"block_number"`
	TransactionHash ethereum.Hash `json:"tx_hash"`
	// EthAmount = OriginalEthAmount * len(BurnFees)
	EthAmount         *big.Int `json:"eth_amount"`
	OriginalEthAmount *big.Int `json:"original_eth_amount"`

	SrcAddress  ethereum.Address `json:"src_addr"`
	SrcSymbol   string           `json:"src_symbol,omitempty"`
	DestAddress ethereum.Address `json:"dst_addr"`
	DestSymbol  string           `json:"dst_symbol,omitempty"`

	UserAddress       ethereum.Address `json:"user_addr"`
	ReceiverAddress   ethereum.Address `json:"receiver_address"`
	SrcReserveAddress ethereum.Address `json:"src_reserve_addr"`
	DstReserveAddress ethereum.Address `json:"dst_reserve_addr"`
	SrcAmount         *big.Int         `json:"src_amount"`
	DestAmount        *big.Int         `json:"dst_amount"`
	FiatAmount        float64          `json:"fiat_amount"`
	WalletAddress     ethereum.Address `json:"wallet_addr"`
	WalletName        string           `json:"wallet_name"`

	SrcBurnAmount      float64 `json:"src_burn_amount"`
	DstBurnAmount      float64 `json:"dst_burn_amount"`
	SrcWalletFeeAmount float64 `json:"src_wallet_fee_amount"`
	DstWalletFeeAmount float64 `json:"dst_wallet_fee_amount"`

	BurnFees       []BurnFee   `json:"-"`
	WalletFees     []WalletFee `json:"-"`
	IntegrationApp string      `json:"integration_app"`

	IP       string           `json:"ip"`
	Country  string           `json:"country"`
	UID      string           `json:"uid"`
	TxSender ethereum.Address `json:"tx_sender"`

	ETHUSDRate     float64 `json:"eth_usd_rate"`
	ETHUSDProvider string  `json:"-"`

	UserName  string `json:"user_name"`
	ProfileID int64  `json:"profile_id"`
	Index     uint   `json:"index"` // the index of event log in transaction receipt

	GasUsed        uint64   `json:"gas_used"`
	GasPrice       *big.Int `json:"gas_price"`
	TransactionFee *big.Int `json:"transaction_fee"`
}

// CrawlResult is result of the crawl
type CrawlResult struct {
	Reserves      []Reserve    `json:"reserves"` // reserve update on this
	UpdateWallets []Reserve    `json:"update_wallets"`
	Trades        []TradelogV4 `json:"trades"`
}

// TradelogV4 is object for tradelog after katalyst upgrade
type TradelogV4 struct {
	Timestamp       time.Time     `json:"timestamp"`
	BlockNumber     uint64        `json:"block_number"`
	TransactionHash ethereum.Hash `json:"tx_hash"`

	TokenInfo TradeTokenInfo `json:"token_info"`
	// support version before katalyst
	SrcReserveAddress ethereum.Address `json:"-"`
	DstReserveAddress ethereum.Address `json:"-"`

	// After katalyst info
	T2EReserves  [][32]byte    `json:"-"` // reserve_id of reserve for trade from token to ether
	E2TReserves  [][32]byte    `json:"-"` // reserve_id of reserve for trade from ether to token
	T2ESrcAmount []*big.Int    `json:"-"`
	E2TSrcAmount []*big.Int    `json:"-"`
	T2ERates     []*big.Int    `json:"-"`
	E2TRates     []*big.Int    `json:"-"`
	Fees         []TradelogFee `json:"fees"`

	// EthAmount = OriginalEthAmount * len(BurnFees)
	EthAmount         *big.Int `json:"eth_amount"`
	OriginalEthAmount *big.Int `json:"original_eth_amount"`
	SrcAmount         *big.Int `json:"src_amount"`
	DestAmount        *big.Int `json:"dst_amount"`
	FiatAmount        float64  `json:"fiat_amount"`
	ETHUSDRate        float64  `json:"eth_usd_rate"`
	ETHUSDProvider    string   `json:"-"`

	WalletAddress ethereum.Address `json:"wallet_addr"`
	WalletName    string           `json:"wallet_name"`

	IntegrationApp string `json:"integration_app"`

	User            KyberUserInfo    `json:"user"`
	ReceiverAddress ethereum.Address `json:"receiver_address"`
	TxDetail        TxDetail         `json:"tx_detail"`
	Split           []TradeSplit     `json:"split"`

	Index   uint `json:"index"`
	Version uint `json:"version"`
}

// TradeSplit split the trade
type TradeSplit struct {
	ReserveAddress ethereum.Address `json:"reserve_address"`
	SrcToken       ethereum.Address `json:"src_token"`
	DstToken       ethereum.Address `json:"dst_token"`
	SrcAmount      *big.Int         `json:"src_amount"`
	DstAmount      *big.Int         `json:"dst_amount"`
	Rate           *big.Int         `json:"rate"`
	Index          uint             `json:"index"`
}

// Reserve represent a reserve in KN
type Reserve struct {
	Address      ethereum.Address `json:"reserve"`
	ReserveID    [32]byte         `json:"reserve_id"`
	ReserveType  uint64           `json:"reserve_type"`
	RebateWallet ethereum.Address `json:"rebate_wallet"`
	BlockNumber  uint64           `json:"block_number"` // block number where reserve value (address, rebate_wallet) is applied
}

// TradelogFee is fee for a trade
type TradelogFee struct {
	ReserveAddr ethereum.Address `json:"reserve_addr"`    // backward compatible for tradelog before katalyst
	WalletName  string           `json:"wallet_name"`     // backward compatible for tradelog before katalyst
	WalletFee   *big.Int         `json:"platform_rebate"` // backward compatible for tradelog before katalyst

	PlatformFee               *big.Int           `json:"platform_fee"`
	PlatformWallet            ethereum.Address   `json:"platform_wallet"`
	Reward                    *big.Int           `json:"reward"`
	Rebate                    *big.Int           `json:"reserve_rebate"`
	RebateWallets             []ethereum.Address `json:"rebate_wallet"`
	RebatePercentBpsPerWallet []*big.Int         `json:"rebate_percent_per_wallet"`
	Burn                      *big.Int           `json:"burn"`
	Index                     uint               `json:"index"`
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
	UserName    string           `json:"user_name"`
	ProfileID   int64            `json:"profile_id"`
	Index       uint             `json:"index"` // the index of event log in transaction receipt
	IP          string           `json:"ip"`
	Country     string           `json:"country"`
	UID         string           `json:"uid"`
}

// BigTradeLog represent trade event on KyberNetwork
type BigTradeLog struct {
	TradelogID        uint64        `json:"tradelog_id"`
	Timestamp         time.Time     `json:"timestamp"`
	TransactionHash   ethereum.Hash `json:"tx_hash"`
	EthAmount         *big.Int      `json:"eth_amount"`
	OriginalETHAmount *big.Int      `json:"original_eth_amount"`
	SrcSymbol         string        `json:"src_symbol,omitempty"`
	DestSymbol        string        `json:"dst_symbol,omitempty"`
	FiatAmount        float64       `json:"fiat_amount"`
}

// MarshalJSON implements custom JSON marshaller for TradeLog to format timestamp in unix millis instead of RFC3339.
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

// UnmarshalJSON implements custom JSON unmarshal for TradeLog
func (tl *TradeLog) UnmarshalJSON(b []byte) error {
	type AliasTradeLog TradeLog
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

// MarshalJSON implements custom JSON marshaller for TradeLog to format timestamp in unix millis instead of RFC3339.
func (tl *TradelogV4) MarshalJSON() ([]byte, error) {
	type AliasTradeLog TradelogV4
	return json.Marshal(struct {
		Timestamp uint64 `json:"timestamp"`
		*AliasTradeLog
	}{
		AliasTradeLog: (*AliasTradeLog)(tl),
		Timestamp:     timeutil.TimeToTimestampMs(tl.Timestamp),
	})
}

// UnmarshalJSON implements custom JSON unmarshal for TradeLog
func (tl *TradelogV4) UnmarshalJSON(b []byte) error {
	type AliasTradeLog TradelogV4
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

// MarshalJSON implements custom JSON marshaller for TradeLog to format timestamp in unix millis instead of RFC3339.
func (tl *BigTradeLog) MarshalJSON() ([]byte, error) {
	type AliasTradeLog BigTradeLog
	return json.Marshal(struct {
		Timestamp uint64 `json:"timestamp"`
		*AliasTradeLog
	}{
		AliasTradeLog: (*AliasTradeLog)(tl),
		Timestamp:     timeutil.TimeToTimestampMs(tl.Timestamp),
	})
}

// UnmarshalJSON implements custom JSON unmarshal for TradeLog
func (tl *BigTradeLog) UnmarshalJSON(b []byte) error {
	type AliasTradeLog BigTradeLog
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
	ETHAmount float64 `json:"eth_amount"`
	USDAmount float64 `json:"usd_amount"`
	Volume    float64 `json:"volume"`
}

// TradeSummary struct holds all the fields required for trade summary
type TradeSummary struct {
	ETHVolume          float64 `json:"eth_volume"`
	USDAmount          float64 `json:"usd_volume"`
	TotalBurnFee       float64 `json:"burn_fee"`
	TotalTrade         uint64  `json:"total_trade"`
	UniqueAddresses    uint64  `json:"unique_addresses"`
	KYCEDAddresses     uint64  `json:"kyced_addresses"`
	NewUniqueAddresses uint64  `json:"new_unique_addresses"`
	USDPerTrade        float64 `json:"usd_per_trade"`
	ETHPerTrade        float64 `json:"eth_per_trade"`
}

//CountryStats stats for a country a day
type CountryStats struct {
	TotalETHVolume     float64 `json:"eth_volume"`
	TotalUSDVolume     float64 `json:"usd_volume"`
	TotalBurnFee       float64 `json:"burn_fee"`
	TotalTrade         uint64  `json:"total_trade"`
	UniqueAddresses    uint64  `json:"unique_addresses"`
	KYCEDAddresses     uint64  `json:"kyced_addresses"`
	NewUniqueAddresses uint64  `json:"new_unique_addresses"`
	USDPerTrade        float64 `json:"usd_per_trade"`
	ETHPerTrade        float64 `json:"eth_per_trade"`
}

//UserVolume represent volume of an user from time to time
type UserVolume struct {
	ETHAmount float64 `json:"eth_amount"`
	USDAmount float64 `json:"usd_amount"`
}

// WalletStats represent stat for a wallet address
type WalletStats struct {
	ETHVolume          float64 `json:"eth_volume"`
	USDVolume          float64 `json:"usd_volume"`
	BurnFee            float64 `json:"burn_fee"`
	TotalTrade         int64   `json:"total_trade"`
	UniqueAddresses    int64   `json:"unique_addresses"`
	KYCEDAddresses     int64   `json:"kyced_addresses"`
	NewUniqueAddresses int64   `json:"new_unique_addresses"`
	USDPerTrade        float64 `json:"usd_per_trade"`
	ETHPerTrade        float64 `json:"eth_per_trade"`
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

//Heatmap represent a country heatmap
type Heatmap struct {
	TotalETHValue        float64 `json:"total_eth_value"`
	TotalTokenValue      float64 `json:"total_token_value"`
	TotalFiatValue       float64 `json:"total_fiat_value"`
	ToTalBurnFee         float64 `json:"total_burn_fee"`
	TotalTrade           int64   `json:"total_trade"`
	TotalUniqueAddresses int64   `json:"total_unique_addr"`
	TotalKYCUser         int64   `json:"total_kyc_user"`
}

var kyberWallets = map[ethereum.Address]struct{}{
	ethereum.HexToAddress("0x440bBd6a888a36DE6e2F6A25f65bc4e16874faa9"): {},
	ethereum.HexToAddress("0xEA1a7dE54a427342c8820185867cF49fc2f95d43"): {},
}

// LengthWalletFees return true if tradelog have no wallet fee
func LengthWalletFees(tradelog TradelogV4) int {
	count := 0
	for _, fee := range tradelog.Fees {
		if fee.WalletFee != nil {
			if fee.WalletFee.Cmp(big.NewInt(0)) != 0 {
				count++
			}
		}
	}
	return count
}

// LengthBurnFees return number of burn fee in a tradelogs
func LengthBurnFees(tradelog TradelogV4) int {
	count := 0
	for _, fee := range tradelog.Fees {
		if fee.Burn != nil {
			if fee.Burn.Cmp(big.NewInt(0)) != 0 {
				count++
			}
		}
	}
	return count
}

func isKyberWallet(addr ethereum.Address) bool {
	if _, exist := kyberWallets[addr]; exist {
		return true
	}
	return false
}

//IsKyberSwap determine if the tradelog is through KyberSwap
func (tl TradelogV4) IsKyberSwap() bool {
	// since block 6715130 KyberSwap add wallet_addr to its tx
	// then we use only this logic to detect if a tx a KyberSwap tx or not
	if tl.BlockNumber >= 6715130 {
		return isKyberWallet(tl.WalletAddress)
	}
	// with older block we use logic below to detect if a tx is a KyberSwap tx
	// if a trade log has no feeToWalletEvent, it is KyberSwap
	if LengthWalletFees(tl) == 0 {
		return true
	}
	// if Wallet Address < maxUint128, it is KyberSwap
	// as a result  of history we used to put block number as wallet address (while other put their real wallet addr)
	// then we use logic below to check if a tx is Kyber Swap tx
	if new(big.Int).SetBytes(tl.WalletAddress.Bytes()).Cmp(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil)) == -1 {
		return true
	}

	return isKyberWallet(tl.WalletAddress)
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
