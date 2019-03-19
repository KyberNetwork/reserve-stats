package common

// CexName is type of Centralized exchange names.
//go:generate stringer -type=CexName -linecomment
type CexName int

const (
	// Binance is exchange name for Binance
	Binance CexName = iota //binance
	//Huobi is exchange name for Huobi
	Huobi // huobi
)

//CexNames map a string to CexName
var CexNames = map[string]CexName{
	"binance": Binance,
	"huobi":   Huobi,
}

//IsValidCexName return if the name is a valid cexname
func IsValidCexName(name string) bool {
	_, ok := CexNames[name]
	return ok
}
