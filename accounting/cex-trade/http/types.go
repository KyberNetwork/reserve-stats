package http

// ConvertTrade ...
type ConvertTrade struct {
	Timestamp    int64   `json:"timestamp"`
	Rate         float64 `json:"rate"`
	AccountName  string  `json:"account_name"`
	Pair         string  `json:"pair"`
	Type         string  `json:"type"`
	Qty          float64 `json:"qty"` // amount of base token
	ETHChange    float64 `json:"eth_change"`
	TokenChange  float64 `json:"token_change"`
	Hash         string  `json:"hash"`
	TakerAddress string  `json:"taker_address"`
	PricingGood  bool    `json:"pricing_good"`
	PnLBPS       float64 `json:"pnl_bps"`
}
