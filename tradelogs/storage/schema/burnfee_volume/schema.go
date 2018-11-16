package burnfeevolume

// FieldName define a list of field names for a TradeLog record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//ReserveAddr is enumerated fieldname for reserve Address
	ReserveAddr // reserve_addr
	//SumAmount is the enumberated field for sum burnfee amount
	SumAmount // sum_amount
)

//burnFeeVolumeFields translates the stringer of reserveRate fields into its enumerated form
var burnFeeVolumeFields = map[string]FieldName{
	"time":         Time,
	"reserve_addr": ReserveAddr,
	"sum_amount":   SumAmount,
}
