package binance

//APIResponse response a basic response from binance
type APIResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//TradeHistory return a history of a trade
type TradeHistory struct {
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

//WithdrawHistory is a binance withdraw history
type WithdrawHistory struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Address   string  `json:"address"`
	Asset     string  `json:"asset"`
	TxID      string  `json:"txId"`
	ApplyTime uint64  `json:"applyTime"`
	Status    int64   `json:"status"`
}

//WithdrawHistoryList is a list of binance withdraw history
type WithdrawHistoryList struct {
	WithdrawList []WithdrawHistory `json:"withdrawList"`
}
