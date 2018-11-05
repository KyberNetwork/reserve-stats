package common

//PriceAnalytic represent a price analytic object
type PriceAnalytic struct {
	ID                   int                 `json:"id"` // id for database item
	Timestamp            uint64              `json:"timestamp" binding:"required"`
	BlockExpiration      bool                `json:"block_expiration" binding:"required"`
	TriggeringTokensList []PriceAnalyticData `json:"triggering_tokens_list" binding:"required"`
}

//PriceAnalyticData represent data for price analytic
type PriceAnalyticData struct {
	Token          string  `json:"token" binding:"required"`
	AskPrice       float64 `json:"ask_price" binding:"required"`
	BidPrice       float64 `json:"bid_price" binding:"required"`
	MinAfpPrice    float64 `json:"min_afp_price" binding:"required"`
	MinAfpOldPrice float64 `json:"min_afp_old_price" binding:"required"`
	MinSpread      float64 `json:"min_spread" binding:"required"`
	TriggerUpdate  bool    `json:"trigger_update" binding:"true"`
}
