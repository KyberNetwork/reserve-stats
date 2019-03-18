package core

// Interface represents a client o interact with KyberNetwork core APIs.
type Interface interface {
	Tokens() ([]Token, error)
}
