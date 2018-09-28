package common

//UserResponse response for /users query
type UserResponse struct {
	KYC bool   `json:"kyc"`
	Cap uint64 `json:"cap"`
	//TODO: get user stats by day and return rich field, now default is false
	Rich bool `json:"rich"`
}

//UserAddress an infomation of an user
type UserAddress struct {
	Address   string `json:"address"`
	Timestamp uint64 `json:"timestamp"`
}

//User information for an user
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Address   string `json:"address" sql:",unique"`
	Timestamp uint64 `json:"timestamp"`
}

//UserCap cap for user
type UserCap struct {
	DailyLimit float64
	TxLimit    float64
}
