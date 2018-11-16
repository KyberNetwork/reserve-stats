package tradelog

// FieldsRegistrar store a map of RateSchemaFieldName to its index in the result from influx db
type FieldsRegistrar map[FieldName]int

// NewFieldsRegistrar returns a FieldsRegistrar from a list of column name and error if occurs
func NewFieldsRegistrar(colums []string) (FieldsRegistrar, error) {
	result := make(FieldsRegistrar)
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
