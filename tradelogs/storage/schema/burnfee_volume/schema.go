package burnfeevolume

// FieldName define a list of field names for a TradeLog record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//SrcReserveAddr is enumerated fieldname for src reserve Address
	SrcReserveAddr // src_reserve_addr
	//DstReserveAddr is enumerated fieldname for src reserve Address
	DstReserveAddr // dst_reserve_addr
	//SumAmount is the enumberated field for sum burnfee amount
	SumAmount // sum_amount
)

//burnFeeVolumeFields translates the stringer of reserveRate fields into its enumerated form
var burnFeeVolumeFields = map[string]FieldName{
	"time":             Time,
	"src_reserve_addr": SrcReserveAddr,
	"dst_reserve_addr": DstReserveAddr,
	"sum_amount":       SumAmount,
}
