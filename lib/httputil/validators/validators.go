package validators

import (
	"reflect"
	"strings"

	tradelog "github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/binding"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gopkg.in/go-playground/validator.v8"
)

// isEthereumAddress is a validator.Func function that returns true if given field
// is a valid Ethereum address.
func isEthereumAddress(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	address := field.String()
	if len(address) != 0 && !common.IsHexAddress(address) {
		return false
	}
	return true
}

// isEmail is a validator.Func function that returns true if given field
// is a valid email address.
func isEmail(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	if err := validation.Validate(field.String(), is.Email); err != nil {
		return false
	}
	return true
}

// isFreq is a validator.Func that returns true if given field is a valid request frequency
// m = minute, h = hour, d = day
func isFreq(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	freq := strings.ToLower(field.String())
	if freq == "h" || freq == "d" {
		return true
	}
	return false
}

// isSupportedTimezone is a validator.Func that returns true if given field is a supported timezone
// supported time range is from -11 to 14
func isSupportedTimezone(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	timezone := field.Int()
	if timezone < -11 || timezone > 14 {
		return false
	}
	return true
}

// isValidCountryCode is  a validator.Func that returns true if given field
// is a valid country code
func isValidCountryCode(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	country := field.String()
	if country == tradelog.UnknownCountry {
		return true
	}
	if err := validation.Validate(country, is.CountryCode2); err != nil {
		return false
	}
	return true
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isAddress", isEthereumAddress)
		v.RegisterValidation("isEmail", isEmail)
		v.RegisterValidation("isFreq", isFreq)
		v.RegisterValidation("isSupportedTimezone", isSupportedTimezone)
		v.RegisterValidation("isValidCountryCode", isValidCountryCode)
	}
}
