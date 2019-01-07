package userkyced

type userKycedReply struct {
	Success bool   `json:"success"`
	Kyced   bool   `json:"kyced"`
	Reason  string `json:"reason"`
}
