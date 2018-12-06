package common

// AppObject an app
type AppObject struct {
	ID        int64    `json:"id,omitempty"`
	AppName   string   `json:"app_name" binding:"required"`
	Addresses []string `json:"addresses" binding:"required"`
}
