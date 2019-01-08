package tokenrate

// FieldName define a list of field names for a token rate record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for time
	Time FieldName = iota //time
	//Provider is enumerated field name for provider
	Provider //provider
	//Rate is the enumerated field for country rate
	Rate //rate
)

//tokenRateFieldNames translates the stringer of reserveRate fields into its enumerated form
var tokenRateFieldNames = map[string]FieldName{
	"time":     Time,
	"provider": Provider,
	"rate":     Rate,
}
