package schema

// TradeLogFieldsRegistrar store a map of RateSchemaFieldName to its index in the result from influx db
type TradeLogFieldsRegistrar map[TradeLogSchemaFieldName]int

// NewTradeLogFieldsRegistrar returns a FieldsRegistrar from a list of column name and error if occurs
func NewTradeLogFieldsRegistrar(colums []string) (TradeLogFieldsRegistrar, error) {
	result := make(TradeLogFieldsRegistrar)
	numFields := 0
	for idx, fieldNameStr := range colums {
		fieldName, ok := tradeLogSchemaFields[fieldNameStr]
		// if a column doesn't map to a field, we wont consider it as an error.
		if ok {
			result[fieldName] = idx
			numFields++
		}
	}
	return result, nil
}
