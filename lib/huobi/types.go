package huobi

//Account represent a huobi account
type Account struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	State  string `json:"state"`
	UserID string `json:"user-id"`
}

//AccountResponse response for accoutn api
type AccountResponse struct {
	Data []Account `json:"data"`
}

//TradeHistory is a history of a trade in huobi
type TradeHistory struct {
	ID              int64  `json:"id"`
	Symbol          string `json:"symbol"`
	AccountID       int64  `json:"account-id"`
	Amount          string `json:"amount"`
	Price           string `json:"price"`
	CreateAt        uint64 `json:"created-at"`
	Type            string `json:"type"`
	FieldAmount     string `json:"field-amount"`
	FieldCashAmount string `json:"field-cash-amount"`
	FieldFee        string `json:"field-fees"`
	FinishedAt      uint64 `json:"finished-at"`
	UserID          int64  `json:"user-id"`
	Source          string `json:"source"`
	State           string `json:"state"`
	CanceledAt      uint64 `json:"canceled-at"`
	Exchange        string `json:"exchange"`
	Batch           string `json:"batch"`
}

//TradeHistoryList is a list of trade history
type TradeHistoryList struct {
	Data []TradeHistory `json:"data"`
}

//WithdrawHistory is history of a withdraw
type WithdrawHistory struct {
	ID                uint64  `json:"id"`
	TransactionID     uint64  `json:"transaction-id"`
	CreatedAt         uint64  `json:"created-at"`
	UpdatedAt         uint64  `json:"updated-at"`
	CandidateCurrency string  `json:"candiate-currency"`
	Currency          string  `json:"currency"`
	Type              string  `json:"type"`
	Direction         string  `json:"direction"`
	Amount            float64 `json:"amount"`
	State             string  `json:"state"`
	Fees              string  `json:"fees"`
	ErrorCode         string  `json:"error-code"`
	ErrorMsg          string  `json:"error-msg"`
	ToAddress         string  `json:"to-address"`
	ToAddrTag         string  `json:"to-addr-tag"`
	TxHash            string  `json:"tx-hash"`
	Chain             string  `json:"chain"`
	Extra             string  `json:"extra"`
}

//WithdrawHistoryList is a list of withdraw history
type WithdrawHistoryList struct {
	Data []WithdrawHistory `json:"data"`
}

//SymbolsReply hold huobi's reply data and status
type SymbolsReply struct {
	Status string   `json:"status"`
	Data   []Symbol `json:"data"`
}

//Symbol struct hold huobi's reply.
type Symbol struct {
	Base            string `json:"base-currency"`
	Quote           string `json:"quote-currency"`
	PricePrecision  int    `json:"price-precision"`
	AmountPrecision int    `json:"amount-precision"`
	SymBolPartition string `json:"symbol-partition"`
	SymBol          string `json:"symbol"`
}

//ExtrasTradeHistoryParams hold extras params for queries
//it included: From - type string, the ID of the trade log to query from
// 			   Size - type string, the limit number of request to return
//			   Direct - type string, the direction to query extra tradelog
type ExtrasTradeHistoryParams struct {
	From   string
	Size   string
	Direct string
}

// ReplyStatus define a list of status code for huobi reply
//go:generate stringer -type=ReplyStatus -linecomment
type ReplyStatus int

const (
	//StatusOK indict that the reply contains correct data
	StatusOK ReplyStatus = iota //ok
	//StatusError indict that there was error in data reply
	StatusError //error
)

//CurrenciesReply hold huobi's reply on currencies endpoint data and status
type CurrenciesReply struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}
