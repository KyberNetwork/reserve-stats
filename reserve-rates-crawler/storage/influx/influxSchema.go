package influx

// RateSchemaFieldName define a list of field names for a rate record in influxDB
//go:generate stringer -type=RateSchemaFieldName -linecomment
type RateSchemaFieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time RateSchemaFieldName = iota //time
	//Pair is enumerated field name for reserveRate.Pair
	Pair //pair
	//BuyRate is enumerated field name for reserveRate.BuyRate
	BuyRate //buy_rate
	//SellRate is enumerated field name for reserveRate.SellRate
	SellRate //sell_rate
	//BuySanityRate is enumerated field name for reserveRate.BuySanityRate
	BuySanityRate //buy_sanity_rate
	//SellSanityRate is enumerated field name for reserveRate.SellSanityRate
	SellSanityRate //sell_sanity_rate
	//BlockNumber is enumerated field name for reserveRate.BlockNumber
	BlockNumber //block_number
	//Reserve is enumerated field name for reserveRate.Reserve
	Reserve //reserve
)

//RateSchemaFields translates the stringer of reserveRate fields into its enumerated form
var RateSchemaFields = map[string]RateSchemaFieldName{
	"time":             Time,
	"pair":             Pair,
	"buy_rate":         BuyRate,
	"sell_rate":        SellRate,
	"buy_sanity_rate":  BuySanityRate,
	"sell_sanity_rate": SellSanityRate,
	"block_number":     BlockNumber,
	"reserve":          Reserve,
}
