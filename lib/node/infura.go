package node

const (
	infuraProjectID = "/v3/a056d94386a94b35a70d60f82745f600"
	// infuraEndpoint is url for infura node
	infuraEndpoint = "https://mainnet.infura.io" + infuraProjectID
)

// InfuraEndpoint returns configured Infura Ethereum node endpoint.
func InfuraEndpoint() string {
	return infuraEndpoint
}
