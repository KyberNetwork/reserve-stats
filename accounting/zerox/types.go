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
