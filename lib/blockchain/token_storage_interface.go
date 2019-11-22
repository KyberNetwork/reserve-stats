package blockchain

type tokenStorageInterface interface {
	GetTokenSymbol(address string) (string, error)
	UpdateTokens(addresses, symbol []string) error
}
