package common

//UsersResponse response for /users query
type UsersResponse struct {
	KYC  bool   `json:"kyc"`
	Cap  uint64 `json:"cap"`
	Rich bool   `json:"rich"`
}

//UserInfo an infomation of an user
type UserInfo struct {
	Address   string `json:"address"`
	Timestamp uint64 `json:"timestamp"`
}

//UserUpdate information for an user
type UserUpdate struct {
	Email string     `json:"email"`
	Infos []UserInfo `json:"info"`
}
