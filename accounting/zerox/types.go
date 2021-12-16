package zerox

// Token ...
type Token struct {
	ID       string `json:"id"`
	Decimals string `json:"decimals"`
	Symbol   string `json:"symbol"`
}

// Transaction ...
type Transaction struct {
	ID string `json:"id"`
}

// Tradelog ...
type Tradelog struct {
	InputToken        Token       `json:"inputToken"`
	InputTokenAmount  string      `json:"inputTokenAmount"`
	OutputToken       Token       `json:"outputToken"`
	OutputTokenAmount string      `json:"outputTokenAmount"`
	Timestamp         string      `json:"timestamp"`
	Transaction       Transaction `json:"transaction"`
}

// TradelogsResponse ...
type TradelogsResponse struct {
	Data struct {
		Maker struct {
			NativeOrderFills []Tradelog `json:"nativeOrderFills"`
		}
	} `json:"data"`
}

// ConvertTrade ...
type ConvertTrade struct {
	Symbol    string  `db:"symbol"`
	Price     float64 `db:"price"`
	Timestamp uint64  `db:"timestamp"`
}

// ConvertTradeInfo ...
type ConvertTradeInfo struct {
	ETHRate        float64 `db:"eth_rate"`
	InToken        string  `db:"in_token"`
	InTokenAmount  float64 `db:"in_token_amount"`
	InTokenRate    float64 `db:"in_token_rate"`
	OutToken       string  `db:"out_token"`
	OutTokenAmount float64 `db:"out_token_amount"`
	OutTokenRate   float64 `db:"out_token_rate"`
	Timestamp      int64   `db:"timestamp"`
}
