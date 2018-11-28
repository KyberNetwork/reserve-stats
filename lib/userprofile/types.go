package userprofile

//UserProfile contain infos of a user
type UserProfile struct {
	UserName  string `json:"name"`
	ProfileID int64  `json:"profile_id"`
}

type userClientReply struct {
	Success bool        `json:"success"`
	Data    UserProfile `json:"data"`
	Reason  []string    `json:"reason"`
}
