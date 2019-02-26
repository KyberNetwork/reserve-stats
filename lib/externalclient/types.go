package externalclient

//BinanceResponse response a basic response from binance
type BinanceResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//BinanceTradeHistory return a history of a trade
type BinanceTradeHistory struct {
	Symbol          string `json:"symbol"`
	ID              int64  `json:"id"`
	OrderID         int64  `json:"orderId"`
	Price           string `json:"price"`
	Quantity        string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            uint64 `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}

//DetailOfAsset return detail of an asset
type DetailOfAsset struct {
	MinWithdrawAmount float64 `json:"minWithdrawAmount"`
	DepositStatus     bool    `json:"depositStatus"`
	WithdrawFee       float64 `json:"withdrawFee"`
	WithdrawStatus    bool    `json:"withdrawStatus"`
	DepositTip        string  `json:"depositTip"`
}

//AssetDetailResponse detail of an asset
type AssetDetailResponse struct {
	AssetDetail map[string]DetailOfAsset `json:"assetDetail"`
}

//BinanceWithdrawHistory is a binance withdraw history
type BinanceWithdrawHistory struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Address   string  `json:"address"`
	Asset     string  `json:"asset"`
	TxID      string  `json:"txId"`
	ApplyTime uint64  `json:"applyTime"`
	Status    int64   `json:"status"`
}

//BinanceWithdrawHistoryList is a list of binance withdraw history
type BinanceWithdrawHistoryList struct {
	WithdrawList []BinanceWithdrawHistory `json:"withdrawList"`
}

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

//HuobiTradeHistory is a history of a trade in huobi
type HuobiTradeHistory struct {
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

//HuobiTradeHistoryList is a list of trade history
type HuobiTradeHistoryList struct {
	Data []HuobiTradeHistory `json:"data"`
}

//HuobiWithdrawHistory is history of a withdraw
type HuobiWithdrawHistory struct {
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

//HuobiWithdrawHistoryList is a list of withdraw history
type HuobiWithdrawHistoryList struct {
	Data []HuobiWithdrawHistory `json:"data"`
}
