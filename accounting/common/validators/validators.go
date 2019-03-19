package validators

import (
	"reflect"
	"strings"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)

// isValidCexName is a validator.Func function that returns true if given field
// is a valid cex name address.
func isValidCexName(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	cexNameInput := strings.ToLower(field.String())
	return common.IsValidCexName(cexNameInput)
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var validators = []struct {
			name string
			fn   validator.Func
		}{
			{"isValidCexName", isValidCexName},
		}
		for _, val := range validators {
			if err := v.RegisterValidation(val.name, val.fn); err != nil {
				panic(err)
			}
		}
	}
}
