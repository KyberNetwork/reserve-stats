package tokenratefetcher

//tokenSymbol is the mapping of tokenID to token Symbol
var tokenSymbols = map[string]string{
	"kyber-network": "knc",
}

//getTokenSymbolFromTokenID return tokenSymbol if avaialbe, else it return the ID itself
func getTokenSymbolFromTokenID(tokenID string) string {
	symbol, ok := tokenSymbols[tokenID]
	if !ok {
		return tokenID
	}
	return symbol
}
