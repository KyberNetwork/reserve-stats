package common

//tokenSymbol is the mapping of tokenID to token Symbol
var tokenSymbols = map[string](map[string]string){
	"coingecko": map[string]string{
		"kyber-network": "knc"},
}

//GetTokenSymbolFromProviderNameTokenID return tokenSymbol if avaialbe, else it return the ID itself
func GetTokenSymbolFromProviderNameTokenID(providerName, tokenID string) string {
	provider, ok := tokenSymbols[providerName]
	if !ok {
		return tokenID
	}
	symbol, ok := provider[tokenID]
	if !ok {
		return tokenID
	}
	return symbol
}
