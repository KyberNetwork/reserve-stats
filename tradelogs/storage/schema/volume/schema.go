package volume

// FieldName define a list of field names for a volume record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for Time
	Time FieldName = iota //time
	//ETHVolume is enumerated field name for ETHVolume
	ETHVolume //eth_volume
	//TokenVolume is enumerated field name for TokenVolume
	TokenVolume //token_volume
	//USDVolume is the enumerated field name for USDVolume
	USDVolume //usd_volume

)

//tradeLogSchemaFields translates the stringer of reserveRate fields into its enumerated form
var voluemFieldSchema = map[string]FieldName{
	"time":         Time,
	"eth_volume":   ETHVolume,
	"token_volume": TokenVolume,
	"usd_volume":   USDVolume,
}
