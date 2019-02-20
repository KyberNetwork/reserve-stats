package userkyced

type userKycedReply struct {
	Kyced bool  `json:"kyced"`
	Error error `json:"error"`
}
