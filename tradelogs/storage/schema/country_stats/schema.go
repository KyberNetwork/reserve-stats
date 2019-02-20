package countrystats

// FieldName define a list of field names for a TradeLog record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for time
	Time FieldName = iota //time
	//UniqueAddresses is enumerated fieldname for number of Unique Addresses in trade summary
	UniqueAddresses // unique_addresses
	//TotalETHVolume is the enumerated fieldname for sum total volume in ETH
	TotalETHVolume // total_eth_volume
	//TotalUSDAmount is the enumerated fieldname for total total USD amount
	TotalUSDAmount // total_usd_amount
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
)

//countryStatsFields translates the stringer of reserveRate fields into its enumerated form
var countryStatsFields = map[string]FieldName{
	"time":                 Time,
	"unique_addresses":     UniqueAddresses,
	"total_eth_volume":     TotalETHVolume,
	"total_usd_amount":     TotalUSDAmount,
	"total_trade":          TotalTrade,
	"usd_per_trade":        USDPerTrade,
	"eth_per_trade":        ETHPerTrade,
	"total_burn_fee":       TotalBurnFee,
	"new_unique_addresses": NewUniqueAddresses,
	"kyced":                KYCedAddresses,
}
