package firsttrade

// FieldName define a list of field names for a firsttrade record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//Traded is the enumerated field for log index
	Traded // traded
)

//tradeLogSchemaFields translates the stringer of reserveRate fields into its enumerated form
var tradeLogSchemaFields = map[string]FieldName{
	"time":   Time,
	"traded": Traded,
}
