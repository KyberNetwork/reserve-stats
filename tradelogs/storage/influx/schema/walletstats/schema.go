package walletstats

// FieldName define a list of field names for a TradeLog record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for time
	Time FieldName = iota //time
	//UniqueAddresses is enumerated fieldname for number of Unique Addresses in trade summary
	UniqueAddresses // unique_addresses
	//ETHVolume is the enumerated fieldname for sum total volume in ETH
	ETHVolume // eth_volume
	//USDVolume is the enumerated fieldname for total total USD amount
	USDVolume // usd_volume
	//TotalTrade is the enumerated fieldname for total trades
	TotalTrade // total_trade
	//USDPerTrade is the enumerated fieldname for average usd per trade
	USDPerTrade // usd_per_trade
	//ETHPerTrade is the enumerated fieldname for average eth per trade
	ETHPerTrade // eth_per_trade
	//TotalBurnFee is the enumerated fieldname for total burnfee
	TotalBurnFee // total_burn_fee
	//NewUniqueAddresses is the enumerated fieldname for new Unique Address in trade summary
	NewUniqueAddresses // new_unique_addresses
	//KYCedAddresses is the enumerated fieldname for number kyced address
	KYCedAddresses // kyced
	//WalletAddress is the enumerated field name for wallet address
	WalletAddress // wallet_addr
)

//walletStatsFields translates the stringer of wallet stats fields into its enumerated form
var walletStatsFields = map[string]FieldName{
	"time":                 Time,
	"unique_addresses":     UniqueAddresses,
	"eth_volume":           ETHVolume,
	"usd_volume":           USDVolume,
	"total_trade":          TotalTrade,
	"usd_per_trade":        USDPerTrade,
	"eth_per_trade":        ETHPerTrade,
	"total_burn_fee":       TotalBurnFee,
	"new_unique_addresses": NewUniqueAddresses,
	"kyced":                KYCedAddresses,
	"wallet_addr":          WalletAddress,
}
