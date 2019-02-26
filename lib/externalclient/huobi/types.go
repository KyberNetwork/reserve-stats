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
	CreateAt        uint64 `json:"create-at"`
	Type            string `json:"type"`
	FieldAmount     string `json:"field-amount"`
	FieldCashAmount string `json:"field-cash-amount"`
	FieldFee        string `json:"field-fee"`
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
	ID                uint64 `json:"id"`
	TransactionID     uint64 `json:"transaction-id"`
	CreatedAt         uint64 `json:"created-at"`
	UpdatedAt         uint64 `json:"updated-at"`
	CandidateCurrency string `json:"candiate-currency"`
	Currency          string `json:"currency"`
	Type              string `json:"type"`
	Direction         string `json:"direction"`
	Amount            string `json:"amount"`
	State             string `json:"state"`
	Fees              string `json:"fees"`
	ErrorCode         string `json:"error-code"`
	ErrorMsg          string `json:"error-msg"`
	ToAddress         string `json:"to-address"`
	ToAddrTag         string `json:"to-addr-tag"`
	TxHash            string `json:"tx-hash"`
	Chain             string `json:"chain"`
	Extra             string `json:"extra"`
}

//WithdrawHistoryList is a list of withdraw history
type WithdrawHistoryList struct {
	Data []WithdrawHistory `json:"data"`
}
