package walletfeevolume

// FieldName define a list of field names for a walletfeevolume record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//SumAmount is the enumerated field for log index
	SumAmount // sum_amount
	//SrcReserveAddress is the enumerated field for reserve Address
	SrcReserveAddress //src_reserve_addr
	//DstReserveAddress is the enumerated field for reserve Address
	DstReserveAddress //dst_reserve_addr
	//WalletAddress is t he enumerated field for wallet Address
	WalletAddress //wallet_addr
)

//walletFeeVolumeFields translates the stringer of walletfeevolume fields into its enumerated form
var walletFeeVolumeFields = map[string]FieldName{
	"time":             Time,
	"sum_amount":       SumAmount,
	"wallet_addr":      WalletAddress,
	"src_reserve_addr": SrcReserveAddress,
	"dst_reserve_addr": DstReserveAddress,
}
