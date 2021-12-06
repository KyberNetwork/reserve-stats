package common

import "math/big"

//DefaultDB is default db for user postgres database
const DefaultDB = "users"

//Info an information of an user
type Info struct {
	//Address add of user
	Address string `json:"address" binding:"required"`
	//Timestamp return timestamp of adding address
	Timestamp int64 `json:"timestamp" binding:"required"`
}

//UserResponse is reponse to user api
type UserResponse struct {
	Cap  *big.Int `json:"cap"`
	Rich bool     `json:"rich"`
}

//UserData user data post through post request to store in stats database
type UserData struct {
	//Email user email - unique
	Email string `json:"email" binding:"required" db:"email"`
	//UserInfo user info include
	UserInfo []Info `json:"user_info" binding:"required,dive"`
}
