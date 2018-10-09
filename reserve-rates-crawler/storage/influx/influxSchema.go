package influx

// RateSchemaFieldName define a list of field names for a rate record in influxDB
//go:generate stringer -type=RateSchemaFieldName -linecomment
type RateSchemaFieldName int

const (
	Time           RateSchemaFieldName = iota //time
	Pair                                      //pair
	BuyRate                                   //buy_rate
	SellRate                                  //sell_rate
	BuySanityRate                             //buy_sanity_rate
	SellSanityRate                            //sell_sanity_rate
	BlockNumber                               //block_number
	Reserve                                   //reserve
)

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
