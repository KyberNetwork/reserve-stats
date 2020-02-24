package tradelogs

import (
	"fmt"

	ethereum "github.com/ethereum/go-ethereum/common"
)

var reserves = map[ethereum.Address]string{
	ethereum.HexToAddress("0x9d27a2d71ac44e075f764d5612581e9afc1964fd"): "Orderbook reserve",
	ethereum.HexToAddress("0xba92981e049a79de1b79c2396d48063e02f47239"): "Bancor hybrid reserve",
	ethereum.HexToAddress("0x44aef3101432a64d1aa16388f4b9b352b09f42a9"): "Oasis hybrid reserve",
	ethereum.HexToAddress("0x5d154c145db2ca90b8ab5e8fe3e716afa4ab7ff0"): "Uniswap hybrid reserve",
	ethereum.HexToAddress("0x6f50e41885fdc44dbdf7797df0393779a9c0a3a6"): "Olympus reserve",
	ethereum.HexToAddress("0x04A487aFd662c4F9DEAcC07A7B10cFb686B682A4"): "Oasis hybrid reserve 2",
	ethereum.HexToAddress("0xcb57809435c66006d16db062c285be9e890c96fc"): "Virgil Capital reserve",
	ethereum.HexToAddress("0xd6000fda0b38f4bff4cfab188e0bd18e8725a5e7"): "DutchX hybrid reserve",
	ethereum.HexToAddress("0x45eb33d008801d547990caf3b63b4f8ae596ea57"): "REN APR rerserve",
	ethereum.HexToAddress("0x57f8160e1c59d16c01bbe181fd94db4e56b60495"): "WETH reserve",
	ethereum.HexToAddress("0x3e9FFBA3C3eB91f501817b031031a71de2d3163B"): "ABYSS APR reserve",
	ethereum.HexToAddress("0xa33c7c22d0bb673c2aea2c048bb883b679fa1be9"): "MLN APR reserve",
	ethereum.HexToAddress("0x13032deb2d37556cf49301f713e9d7e1d1a8b169"): "Uniswap hybrid reserve 2",
	ethereum.HexToAddress("0x5b756435bf2c8895bab3e3898dd7ed2ba073d7b9"): "Bancor hybrid reserve 2",
	ethereum.HexToAddress("0xa9312cb86d1e532b7c21881ce03a1a9d52f6adb1"): "TTC reserve",
	ethereum.HexToAddress("0x8463fDa3567D9228D6Bc2A9b6219fC85a19b89aa"): "Oasis hybrid reserve 3",
	ethereum.HexToAddress("0x2295fc6BC32cD12fdBb852cFf4014cEAc6d79C10"): "PT reserve",
	ethereum.HexToAddress("0x63825c174ab367968ec60f061753d3bbd36a0d8f"): "Kyber reserve",
	ethereum.HexToAddress("0x35183769bbbf63d2b4cac32ef593f4ad08104fba"): "KCC reserve",
	ethereum.HexToAddress("0x21433dec9cb634a23c6a4bbcce08c83f5ac2ec18"): "Prycto reserve",
	ethereum.HexToAddress("0xfe4474d73be9307ebb5b5519dca19e8109286acb"): "Tomo Reserve",
	ethereum.HexToAddress("0x2631a5222522156dfafaa5ca8480223d6465782d"): "Dether reserve",
	ethereum.HexToAddress("0x494696162d3c21b4b8ee08a7fcecc9b4a1dd1566"): "Tvnd reserve",
	ethereum.HexToAddress("0xe0e1f00a2537eccdbb993929a4265658353affc6"): "Mossland reserve",
	ethereum.HexToAddress("0x91be8fa21dc21cff073e07bae365669e154d6ee1"): "BBO APR reserve",
	ethereum.HexToAddress("0xc97094dced8b43be3d275e725f41e63eba2d4cb6"): "Snap reserve",
	ethereum.HexToAddress("0xb50b0d0ed29603c66c65c0582cf9e49b6a9e9da5"): "DCC reserve",
	ethereum.HexToAddress("0x56e37b6b79d4e895618b8bb287748702848ae8c0"): "Midas reserve",
	ethereum.HexToAddress("0x2aab2b157a03915c8a73adae735d0cf51c872f31"): "Prycto reserve 2",
	ethereum.HexToAddress("0x742e8bb8e6bde9cb2df5449f8de7510798727fb1"): "Mossland reserve 2",
	ethereum.HexToAddress("0xc935cad589bebd8673104073d5a5eccfe67fb7b1"): "CoinFi reserve",
	ethereum.HexToAddress("0x582ea0af091ae0d98fdf08216cb2846711a65f6a"): "Kyber reserve 2",
	ethereum.HexToAddress("0xe1213e46efcb8785b47ae0620a51f490f747f1da"): "Prycto reserve 3",
	ethereum.HexToAddress("0x4d864b5b4f866f65f53cbaad32eb9574760865e6"): "Snap reserve 2",
	ethereum.HexToAddress("0x5337d1df2d450945392d60b35f562b92fd96b6b6"): "ABYSS APR reserve 2",
	ethereum.HexToAddress("0x9e2b650f890236ab49609c5a6b00cddb4e61f408"): "MKR, DAI reserve",
	ethereum.HexToAddress("0x8bf5c569ecfd167f96fae6d9610e17571568a6a1"): "DAI reserve",
	ethereum.HexToAddress("0x148332cd398321989f37803188b9a69fa32b133c"): "Kyber reserve 3",
	ethereum.HexToAddress("0xA467b88BBF9706622be2784aF724C4B44a9d26F4"): "KNC APR reserve",
	ethereum.HexToAddress("0x607d7751d9F4845C5a1dE9eeD39c56f4fC0F855d"): "KNC APR reserve 2",
	ethereum.HexToAddress("0x1c802020eea688e2b05936cdb98b8e6894acc1c2"): "ABYSS APR reserve 3",
	ethereum.HexToAddress("0x1670dfb52806de7789d5cf7d5c005cf7083f9a5d"): "USDC APR reserve",
	ethereum.HexToAddress("0x485c4ec93d18ebd16623d455567886475ae28d04"): "WBTC APR reserve",
	ethereum.HexToAddress("0x95f1f428485Bd41729938D620af61718Ea9B1F9E"): "Axe Capital",
	ethereum.HexToAddress("0xa107dfa919c3f084a7893a260b99586981beb528"): "SNX APR reserve",
	ethereum.HexToAddress("0xcf1394c5e2e879969fdb1f464ce1487147863dcb"): "Oasis bridge reserve - v2",
	ethereum.HexToAddress("0xAA14DCAA0AdbE79cBF00edC6cC4ED17ed39240AC"): "DAO stack APR reserve",
	ethereum.HexToAddress("0xb45C8956a080d336934cEE52A35D4dbABF025b6F"): "MKR APR reserve",
	ethereum.HexToAddress("0x05461124c86c0ad7c5d8e012e1499fd9109ffb7d"): "GNO APR reserve",
	ethereum.HexToAddress("0x4Cb01bd05E4652CbB9F312aE604f4549D2bf2C99"): "Synth USD APR reserve",

	ethereum.HexToAddress("0x54A4a1167B004b004520c605E3f01906f683413d"): "Uniswap bridge reserve v3",
	ethereum.HexToAddress("0x3480e12b6c2438e02319e34b4c23770679169190"): "TKN APR reserve",
	ethereum.HexToAddress("0x08030715560a146e306b87ca93fd618bb2a80363"): "BTU APR reserve",
	ethereum.HexToAddress("0x751eea622edd1e3d768c18afbcaec7dce7750c65"): "RAE APR reserve",
	ethereum.HexToAddress("0x1833ad67362249823515b59a8aa8b4f6b4358d1b"): "MYB APR reserve",

	ethereum.HexToAddress("0x053aa84fcc676113a57e0ebb0bd1913839874be4"): "Bancor Reserve",
	ethereum.HexToAddress("0xa9742ee9a5407f4c2f8a49f65e3a440f3694960a"): "Santiment Reserve",
	ethereum.HexToAddress("0x7e2fd015616263add31a2acc2a437557cee80fc4"): "UPP Reserve",
	ethereum.HexToAddress("0xc6c8bce5e9383df025f982d6bbd84163957a6979"): "Nexxo Reserve",
	ethereum.HexToAddress("0x6b84dbd29643294703dbabf8ed97cdef74edd227"): "Sapien",

	ethereum.HexToAddress("0x1fe867bfe9cbe0045467605b959a355223e3885d"): "Bancor's Bridge Reserve",
	ethereum.HexToAddress("0x31e085afd48a1d6e51cc193153d625e8f0514c7f"): "Uniswap Bridge Reserve V4",

	ethereum.HexToAddress("0x1e158c0e93c30d24e918ef83d1e0be23595c3c0f"): "Oasis Bridge Reserve V3",
	ethereum.HexToAddress("0x4f32BbE8dFc9efD54345Fc936f9fEF1048746fCF"): "OneBit Quant",
}

// ReserveAddressToName return reserve name by its address
func ReserveAddressToName(address ethereum.Address) (string, error) {
	if name, existed := reserves[address]; existed {
		return name, nil
	}
	return address.Hex(), fmt.Errorf("address does not have a name: %s", address.Hex())
}
