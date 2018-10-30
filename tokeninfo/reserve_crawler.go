package tokeninfo

import (
	"bytes"
	"encoding/json"
	"math/big"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

const emptyErrMsg = "abi: unmarshalling empty output"

var (
	reserveNames = map[common.Address]string{}
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
	sugar                 *zap.SugaredLogger
	internalNetworkClient *contracts.InternalNetwork
}

// NewReserveCrawler creates a new ReserveCrawler instance.
func NewReserveCrawler(sugar *zap.SugaredLogger, internalNetworkClient *contracts.InternalNetwork) (*ReserveCrawler, error) {
	return &ReserveCrawler{
		sugar:                 sugar,
		internalNetworkClient: internalNetworkClient,
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

	for _, token := range tokens {
		var reserveAddrs = make(map[common.Address]bool)
		result[token.Name] = []*ReserveInfo{}

		f.sugar.Infow("fetching reserve info",
			"token", token.Name)

		for i := 0; ; i++ {
			var reserveAddr common.Address
			reserveAddr, err = f.internalNetworkClient.ReservesPerTokenSrc(
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
