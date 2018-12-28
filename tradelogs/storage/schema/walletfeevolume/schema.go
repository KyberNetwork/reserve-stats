package walletfeevolume

// FieldName define a list of field names for a walletfeevolume record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//SumAmount is the enumerated field for log index
	SumAmount // sum_amount
	//ReserveAddress is the enumerated field for reserve Address
	ReserveAddress //reserve_addr
	//WalletAddress is t he enumerated field for wallet Address
	WalletAddress //wallet_address
)

//walletFeeVolumeFields translates the stringer of walletfeevolume fields into its enumerated form
var walletFeeVolumeFields = map[string]FieldName{
	"time":         Time,
	"sum_amount":   SumAmount,
	"wallet_addr":  WalletAddress,
	"reserve_addr": ReserveAddress,
}
