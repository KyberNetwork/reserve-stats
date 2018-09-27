package common

// ETHToken is the token object represent ETH in the system
var ETHToken = Token{
	ID:       "ETH",
	Name:     "Ethereum",
	Address:  "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
	Decimals: 18,
	Active:   true,
	Internal: true,
}

const (
	//InfuraEndpoint: url for infura node
	InfuraEndpoint = "https://mainnet.infura.io"
	NetworkAddrV1  = "0x964F35fAe36d75B1e72770e244F6595B68508CF5"
	NetworkAddrV2  = "0x818E6FECD516Ecc3849DAf6845e3EC868087B755"

	WrapperAddrV1 = "0x533e6d1ffa2b96cf9c157475c76c38d1b13bc584"
	WrapperAddrV2 = "0x6172AFC8c00c46E0D07ce3AF203828198194620a"

	ReserveAddr = "0x63825c174ab367968EC60f061753D3bbD36A0D8F"

	StartingBlockV2 = 5926056
)
