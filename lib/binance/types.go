package binance

//APIResponse response a basic response from binance
type APIResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
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

//ExchangeInfo is info of binance
type ExchangeInfo struct {
	Timezone   string      `json:"timezone"`
	ServerTime uint64      `json:"serverTime"`
	RateLimits []RateLimit `json:"rateLimits"`
	Symbols    []Symbol    `json:"symbols"`
}

//RateLimit is a rate limit type
type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Internval     string `json:"interval"`
	IntervalNum   int    `json:"internvalNum"`
	Limit         int    `json:"limit"`
}

//Symbol is token symbol from binance
type Symbol struct {
	Symbol              string             `json:"symbol"`
	Status              string             `json:"status"`
	BaseAsset           string             `json:"baseAsset"`
	BaseAssetPrecision  int                `json:"baseAssetPrecision"`
	QuoteAsset          string             `json:"quoteAsset"`
	QuoteAssetPrecision int                `json:"quoteAssetPrecision"`
	OrderTypes          []string           `json:"orderTypes"`
	IcebergAllowed      bool               `json:"icebergAllowed"`
	Filters             []SymbolFilterType `json:"filters"`
}

//SymbolFilterType is a
type SymbolFilterType struct {
	FilterType       string `json:"filterType"`
	MinPrice         string `json:"minPrice,omitempty"`
	MaxPrice         string `json:"maxPrice,omitempty"`
	TickSize         string `json:"tickSize,omitempty"`
	MultiplierUp     string `json:"multiplierUp,omitempty"`
	MultiplierDown   string `json:"multiplierDown,omitempty"`
	AvgPriceMins     int    `json:"avgPriceMins,omitempty"`
	MinQty           string `json:"minQty,omitempty"`
	MaxQty           string `json:"maxQty,omitempty"`
	StepSize         string `json:"stepSize,omitempty"`
	MinNotional      string `json:"minNotional,omitempty"`
	ApplytoMarket    bool   `json:"applyToMarket,omitempty"`
	Limit            int    `json:"limit,omitempty"`
	MaxNumAlgoOrders int    `json:"MaxNumAlgoOrders,omitempty"`
}
