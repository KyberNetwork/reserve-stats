package firsttrade

// FieldName define a list of field names for a firsttrade record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//Traded is the enumerated field for log index
	Traded // traded
	//WalletAddress is the enumerated field for wallet Address
	WalletAddress //wallet_addr
	//Country is t he enumerated field for country
	Country //country
)
