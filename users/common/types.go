package common

import "time"

//UserResponse response for /users query
type UserResponse struct {
	//KYC return if user is kyc or not
	KYC bool `json:"kyc"`
	//Cap return transaction cap for user by wei
	Cap uint64 `json:"cap"`
	//Rich return true if user is exceed his daily cap
	Rich bool `json:"rich"`
}

//UserAddress an information of an user
type UserAddress struct {
	//Address add of user
	Address string `json:"address"`
	//Timestamp return timestamp of adding address
	Timestamp time.Time `json:"timestamp"`
}

//User information for an user
type User struct {
	//ID of user in postgres db
	ID int64 `json:"id"`
	//Email user email
	Email string `json:"email"`
	//Address of user
	Address string `json:"address" sql:",unique"`
	//Timestamp of user adding
	Timestamp time.Time `json:"timestamp"`
}

//UserCap cap for user
type UserCap struct {
	//DailyLimit for user
	DailyLimit float64 `json:"daily-limit"`
	//TxLimit cap for each transaction
	TxLimit float64 `json:"tx-limit"`
}
