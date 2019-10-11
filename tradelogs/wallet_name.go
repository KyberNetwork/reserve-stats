package tradelogs

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

var walletNamingMap = map[ethereum.Address]string{
	ethereum.HexToAddress("0xf89220007d9280f97FA44C8bB82EfBdEcC39063A"): "Fulcrum",
	ethereum.HexToAddress("0xdE63aef60307655405835DA74BA02CE4dB1a42Fb"): "Enjin",
	ethereum.HexToAddress("0xb9E29984Fe50602E7A619662EBED4F90D93824C7"): "ImToken",
	ethereum.HexToAddress("0xb21090C8f6bAC1ba614A3F529aAe728eA92B6487"): "Multis",
	ethereum.HexToAddress("0xa7615CD307F323172331865181DC8b80a2834324"): "Easwap",
	ethereum.HexToAddress("0xa6bC6dF9Eba23abfF0d1eCD6C9847893D2B1643D"): "CoinManager",
	ethereum.HexToAddress("0xa5c603e1C27a96171487aea0649b01c56248d2e8"): "Argent",
	ethereum.HexToAddress("0xF7075e232b34E57Ca3bB91980b97C4f8a20d7ee4"): "CoinGecko",
	ethereum.HexToAddress("0xF257246627f7CB036AE40Aa6cFe8D8CE5F0EbA63"): "Fulcrum",
	ethereum.HexToAddress("0xF1AA99C69715F423086008eB9D06Dc1E35Cc504d"): "Trust",
	ethereum.HexToAddress("0xF12c4E73868a4A028382AC51b57482b627A323d2"): "Nuo",
	ethereum.HexToAddress("0xEC1e3dc16eE138991E105DfA3230F1c9D607A6d0"): "Fulcrum",
	ethereum.HexToAddress("0xEA1a7dE54a427342c8820185867cF49fc2f95d43"): "KyberSwap Non-EU",
	ethereum.HexToAddress("0xDECAF9CD2367cdbb726E904cD6397eDFcAe6068D"): "Myetherwallet",
	ethereum.HexToAddress("0xDD61803d4a56C597E0fc864F7a20eC7158c6cBA5"): "Cipher",
	ethereum.HexToAddress("0xC9D81352fBdb0294b091e51d774A0652ef776D99"): "Unknown Arbitrage Bot",
	ethereum.HexToAddress("0xB4700Da07508553877A81a9A2F40a872DE788cfE"): "Fulcrum",
	ethereum.HexToAddress("0x9a68f7330A3Fe9869FfAEe4c3cF3E6BBef1189Da"): "KyberSwap iOS",
	ethereum.HexToAddress("0x9E1c71c25111F4CA7B40C956c8a21B6AC2f02274"): "Fulcrum",
	ethereum.HexToAddress("0x92afB508a46494AC00A242627703d1f21CA2dF1B"): "Fulcrum",
	ethereum.HexToAddress("0x7A342739F58A55a3a01Efea152EFd95E8e96ef70"): "Fulcrum",
	ethereum.HexToAddress("0x7284a8451d9a0e7Dc62B3a71C0593eA2eC5c5638"): "Instadapp",
	ethereum.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F"): "Etherscan",
	ethereum.HexToAddress("0x673d26360Af6688fDD9d788677fD06f58aad5b4D"): "Midas",
	ethereum.HexToAddress("0x52D35e8f0Ffa18337B093Aec3DfFF40445d8f4f4"): "prod-limit-order",
	ethereum.HexToAddress("0x468fbBCB28E4D2699139c64551D6F0178760209F"): "prod-binance-deposit",
	ethereum.HexToAddress("0x4247951c2eb6d0bA38d233fe7d542c8c80c9d46A"): "KyberSwap EU",
	ethereum.HexToAddress("0x9a68f7330A3Fe9869FfAEe4c3cF3E6BBef1189Da"): "MEW",
	ethereum.HexToAddress("0x9a68f7330A3Fe9869FfAEe4c3cF3E6BBef1189Da"): "KyberSwap iOS",
	ethereum.HexToAddress("0x398d297BAB517770feC4d8Bb7a4127b486c244bB"): "Dex wallet",
	ethereum.HexToAddress("0x332D87209f7c8296389C307eAe170c2440830A47"): "Betoken",
	ethereum.HexToAddress("0x322d58b9E75a6918f7e7849AEe0fF09369977e08"): "CDP saver",
	ethereum.HexToAddress("0x25E3d9B98A4DeA9809B65045D1F007335032EDd4"): "Infinito (IBL)",
	ethereum.HexToAddress("0x21357B3dcb7AE07Da23A708DBbd9a2340001a3F4"): "LinkTime",
	ethereum.HexToAddress("0x1bF3e7EDE31dBB93826C2aF8686f80Ac53f9ed93"): "ipfswap.com",
	ethereum.HexToAddress("0x1a719375E9b8b056C5492Fdf7BAd9bf5A2F79cC2"): "Altitude games",
	ethereum.HexToAddress("0x13ddAC8d492E463073934E2a101e419481970299"): "Fulcrum",
	ethereum.HexToAddress("0x09227deaeE08a5Ba9D6Eb057F922aDfAd191c36c"): "OlympusLab",
	ethereum.HexToAddress("0x087aC7736469716D73498e479E09119A02D7A59D"): "Opyn",
	ethereum.HexToAddress("0x03E0635A77Ca3DbC23748aF10a568663964f4BAD"): "Fulcrum",
	ethereum.HexToAddress("0x3fFFF2F4f6C0831FAC59534694ACd14AC2Ea501b"): "KyberSwap Android",
}

// WalletAddrToName convert eth addr to name
func WalletAddrToName(addr ethereum.Address) string {
	if name, ok := walletNamingMap[addr]; ok {
		return name
	}
	return ""
}
