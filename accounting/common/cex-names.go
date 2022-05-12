package common

// CEXName is type of Centralized exchange names.
//go:generate stringer -type=CEXName -linecomment
type CEXName int

const (
	// Binance is exchange name for Binance
	Binance CEXName = iota // binance
	// Huobi is exchange name for Huobi
	Huobi // huobi
)

// CEXNames map a string to CEXName
var CEXNames = map[string]CEXName{
	"binance": Binance,
	"huobi":   Huobi,
}

// IsValidCEXName return if the name is a valid cexname
func IsValidCEXName(name string) bool {
	_, ok := CEXNames[name]
	return ok
}
