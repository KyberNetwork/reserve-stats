package userkyced

type userKycedReply struct {
	Kyced  bool   `json:"kyced"`
	Reason string `json:"reason"`
}
