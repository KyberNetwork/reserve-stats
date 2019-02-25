package kyced

// FieldName define a list of field names for a kyced record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for time
	Time FieldName = iota //time
	//UserAddress is the field name for user address
	UserAddress //user_addr
	//Country is the field name for country
	Country //country
	//KYCed is the field name for Kyc status
	KYCed //kyced
	//WalletAddress is the field name for wallet address
	WalletAddress //wallet_addr
)

//kycedSchemaFields translates the stringer of kyc fields into its enumerated form
var kycedSchemaFields = map[string]FieldName{
	"time":        Time,
	"user_addr":   UserAddress,
	"country":     Country,
	"kyced":       KYCed,
	"wallet_addr": WalletAddress,
}
