package common

// BinanceWithdrawStatus is type of binance withdraw order states.
//go:generate stringer -type=BinaneWithdrawStatus -linecomment
type BinanceWithdrawStatus int

const (
	// EmailSent ís initial status of binance withdraw
	EmailSent BinanceWithdrawStatus = iota // email-sent
	// Cancelled withdraw status
	Cancelled // cancelled
	// AwaitingApproval withdraw status
	AwaitingApproval // awaiting-approval
	// Rejected withdraw status
	Rejected // rejected
	// Processing withdraw status
	Processing // processing
	// Failure withdraw status
	Failure // failure
	// Completed withdraw status
	Completed // completed
)

// WithdrawStatuses map a string to BinanceWithdrawStatus
var WithdrawStatuses = map[string]BinanceWithdrawStatus{
	"Email Sent":        EmailSent,
	"Cancelled":         Cancelled,
	"Awaiting Approval": AwaitingApproval,
	"Rejected":          Rejected,
	"Processing":        Processing,
	"Failure":           Failure,
	"Completed":         Completed,
}
