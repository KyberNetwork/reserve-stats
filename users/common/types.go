package common

var (
	nonKYCCap = &UserCap{
		DailyLimit: 15000.0,
		TxLimit:    3000.0,
	}
	kycCap = &UserCap{
		DailyLimit: 200000.0,
		TxLimit:    6000.0,
	}
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
	Email string `json:"email" binding:"required,isemail" db:"email"`
	//UserInfo user info include
	UserInfo []Info `json:"user_info" binding:"required,dive,required"`
}

//UserCap is users transaction cap.
type UserCap struct {
	// DailyLimit is the USD amount if the user is considered rich
	// and will receive different rates.
	DailyLimit float64 `json:"daily_limit"`
	// TxLimit is the maximum value in USD of a transaction an user
	// is allowed to make.
	TxLimit float64 `json:"tx_limit"`
}

// NewUserCap returns user cap based on KYC status.
func NewUserCap(kyced bool) *UserCap {
	if kyced {
		return kycCap
	}
	return nonKYCCap
}
