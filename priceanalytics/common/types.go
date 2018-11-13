package common

//PriceAnalytic represent a price analytic object
type PriceAnalytic struct {
	ID                   int                 `json:"-" db:"id"` // id for database item
	Timestamp            uint64              `json:"timestamp" binding:"required" db:"timestamp"`
	BlockExpiration      bool                `json:"block_expiration" binding:"required" db:"block_expiration"`
	TriggeringTokensList []PriceAnalyticData `json:"triggering_tokens_list" binding:"required" db:"triggering_token_list"`
}

//PriceAnalyticData represent data for price analytic
type PriceAnalyticData struct {
	ID              int     `json:"-" db:"id"`
	Token           string  `json:"token" binding:"required" db:"token"`
	AskPrice        float64 `json:"ask_price" binding:"required" db:"ask_price"`
	BidPrice        float64 `json:"bid_price" binding:"required" db:"bid_price"`
	MidAfpPrice     float64 `json:"mid_afp_price" binding:"required" db:"mid_afp_price"`
	MidAfpOldPrice  float64 `json:"mid_afp_old_price" binding:"required" db:"mid_afp_old_price"`
	MinSpread       float64 `json:"min_spread" binding:"required" db:"min_spread"`
	TriggerUpdate   bool    `json:"trigger_update" binding:"true" db:"trigger_update"`
	PriceAnalyticID int     `json:"-" db:"price_analytic_id"`
}

//OldPriceAnalyticFormat represent old object of price analytic
type OldPriceAnalyticFormat struct {
	Timestamp uint64 `json:"timestamp"`
	Value     struct {
		BlockExpiration      bool                `json:"block_expiration"`
		TriggeringTokensList []PriceAnalyticData `json:"triggering_tokens_list"`
	} `json:"value"`
}
