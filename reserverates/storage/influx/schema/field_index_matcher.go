package schema

import "errors"

// FieldsRegistrar store a map of RateSchemaFieldName to its index in the result from influx db
type FieldsRegistrar map[RateSchemaFieldName]int

// NewFieldsRegistrar returns a FieldsRegistrar from a list of collumn name and error if occurs
func NewFieldsRegistrar(colums []string) (FieldsRegistrar, error) {
	result := make(FieldsRegistrar)
	numFields := 0
	for idx, fieldNameStr := range colums {
		fieldName, ok := rateSchemaFields[fieldNameStr]
		// if a column doesn't map to a field, we wont consider it as an error.
		if ok {
			result[fieldName] = idx
			numFields++
		}
	}
	//if a field doesn't map to a column, this mean error.
	if numFields < len(rateSchemaFields) {
		return nil, errors.New("Not enough columns compared to field name")
	}
	return result, nil
}
