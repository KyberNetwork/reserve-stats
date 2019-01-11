package tradelogspostprocessor

// FieldName define a list of field names for a TradeLog record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//ReserveAddr is reserve address
	ReserveAddr //reserve_addr
	//ETHVolume is volume of eth monthly
	ETHVolume //eth_volume
	//USDVolume is volume of usd monthly
	USDVolume //usd_volume
	//BurnFee is amount burn fee
	BurnFee //burn_fee
	//WalletFee is amount wallet fee
	WalletFee //wallet_fee
)

//tradelogsPostProcessor translates the stringer of reserveRate fields into its enumerated form
var tradelogsPostProcessor = map[string]FieldName{
	"time":         Time,
	"reserve_addr": ReserveAddr,
	"eth_volume":   ETHVolume,
	"usd_volume":   USDVolume,
	"burn_fee":     BurnFee,
	"wallet_fee":   WalletFee,
}
