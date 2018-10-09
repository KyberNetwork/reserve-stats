package common

import (
	"time"
)

//Info an information of an user
type Info struct {
	//Address add of user
	Address string `json:"address" binding:"required,isAddress"`
	//Timestamp return timestamp of adding address
	Timestamp int64 `json:"timestamp" binding:"required"`
}

//UserData user data post through post request to store in stats database
type UserData struct {
	//Email user email - unique
	Email string `json:"email" binding:"required,isemail"`
	//UserInfo user info include
	UserInfo []Info `json:"user_info" binding:"required,dive,required"`
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
