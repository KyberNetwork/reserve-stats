package common

// TransactionType define a list of transactionType of reserve transactions
//go:generate stringer -type=TransactionType -linecomment
type TransactionType int

const (
	//ERC20 is the type for erc20 transfer
	ERC20 TransactionType = iota //erc20
	//Normal is the type for ethereum standard transfer
	Normal //normal
	//Internal is the type for internal transfer
	Internal //internal
)

//TransactionTypes contain map of valid string to TransactionType
var TransactionTypes = map[string]TransactionType{
	`erc20`:    ERC20,
	`normal`:   Normal,
	`internal`: Internal,
}
