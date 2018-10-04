package tokeninfo

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const emptyErrMsg = "abi: unmarshalling empty output"

var (
	reserveNames = map[common.Address]string{
		common.HexToAddress("0x63825c174ab367968ec60f061753d3bbd36a0d8f"): "KN",
		common.HexToAddress("0x21433dec9cb634a23c6a4bbcce08c83f5ac2ec18"): "Prycto",
		common.HexToAddress("0x6f50e41885fdc44dbdf7797df0393779a9c0a3a6"): "MOT",
		common.HexToAddress("0x4d864b5b4f866f65f53cbaad32eb9574760865e6"): "SNAP",
		common.HexToAddress("0x91be8fa21dc21cff073e07bae365669e154d6ee1"): "BigBom",
		common.HexToAddress("0xc935cad589bebd8673104073d5a5eccfe67fb7b1"): "CoinFi",
		common.HexToAddress("0x742e8bb8e6bde9cb2df5449f8de7510798727fb1"): "Moss Coin",
		common.HexToAddress("0x8bf5c569ecfd167f96fae6d9610e17571568a6a1"): "Oasis Integration (KN)",
		common.HexToAddress("0xcb57809435c66006d16db062c285be9e890c96fc"): "Virgil Capital",
		common.HexToAddress("0x56e37b6b79d4E895618B8Bb287748702848Ae8c0"): "Midas Protocol",
	}
)

type tokenInfo struct {
	Name    string
	Address common.Address
}

// ReserveInfo is the information of a KyberNetwork reserve.
type ReserveInfo struct {
	Name    string
	Address common.Address
}

// ReserveCrawler gets the tokeninfo reserve mapping information from blockchain.
type ReserveCrawler struct {
	sugar  *zap.SugaredLogger
	client *ethclient.Client
}

// NewReserveCrawler creates a new ReserveCrawler instance.
func NewReserveCrawler(sugar *zap.SugaredLogger, nodeURL string) (*ReserveCrawler, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, err
	}
	return &ReserveCrawler{
		sugar:  sugar,
		client: client,
	}, nil
}

// Fetch returns the reserve information of all tokens.
func (f *ReserveCrawler) Fetch() (map[string][]*ReserveInfo, error) {
	var (
		tokens []*tokenInfo
		result = make(map[string][]*ReserveInfo)
	)

	err := json.NewDecoder(bytes.NewReader([]byte(tokenData))).Decode(&tokens)
	if err != nil {
		return nil, err
	}

	internalNetworkClient, err := contracts.NewInternalNetwork(
		common.HexToAddress(contracts.InternalNetworkContractAddress),
		f.client)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		var reserveAddrs = make(map[common.Address]bool)
		result[token.Name] = []*ReserveInfo{}

		f.sugar.Infow("fetching reserve info",
			"token", token.Name)

		for i := 0; ; i++ {
			var reserveAddr common.Address
			reserveAddr, err = internalNetworkClient.ReservesPerTokenSrc(
				nil,
				token.Address,
				big.NewInt(int64(i)))

			if err != nil {
				if err.Error() == emptyErrMsg {
					break
				}
				return nil, err
			}
			reserveAddrs[reserveAddr] = true
		}

		for reserveAddr := range reserveAddrs {
			result[token.Name] = append(result[token.Name], &ReserveInfo{Name: reserveNames[reserveAddr], Address: reserveAddr})
		}
	}
	return result, nil
}
