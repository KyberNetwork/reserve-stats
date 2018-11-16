package walletfee

// FieldName define a list of field names for a TradeLog record
//go:generate stringer -type=FieldName -linecomment
type FieldName int

const (
	//Time is enumerated field name for reserveRate.time
	Time FieldName = iota //time
	//TxHash is enumerated field name for TxHash
	TxHash //tx_hash
	//ReserveAddr is enumerated fieldname for reserve Address
	ReserveAddr // reserve_addr
	//WalletAddr is enumerated fieldname for wallet Address
	WalletAddr // wallet_addr
	//LogIndex is the enumerated field for log index
	LogIndex // log_index
	//TradeLogIndex is the enumerated field for tradelog index
	TradeLogIndex // trade_log_index
	//Amount is the enumerated field for burnAmount
	Amount //amount
	//Country is the enumerated field for country name
	Country //country
)

//walletFeeSchemaFields translates the stringer of reserveRate fields into its enumerated form
var walletFeeSchemaFields = map[string]FieldName{
	"time":            Time,
	"tx_hash":         TxHash,
	"reserve_addr":    ReserveAddr,
	"wallet_addr":     WalletAddr,
	"log_index":       LogIndex,
	"trade_log_index": TradeLogIndex,
	"amount":          Amount,
}
