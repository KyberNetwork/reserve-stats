package listedtoken

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

//Fetcher to get token listed in a reserve
type Fetcher struct {
	ethClient                 *ethclient.Client
	contractTimestampResolver *blockchain.EtherscanContractTimestampResolver
	sugar                     *zap.SugaredLogger
}

//NewListedTokenFetcher return new fetcher for listed token
func NewListedTokenFetcher(ethClient *ethclient.Client, contractTimestampResolver *blockchain.EtherscanContractTimestampResolver,
	sugar *zap.SugaredLogger) *Fetcher {
	return &Fetcher{
		ethClient:                 ethClient,
		contractTimestampResolver: contractTimestampResolver,
		sugar:                     sugar,
	}
}

func updateListedToken(listedToken map[string]common.ListedToken, symbol, name string, address ethereum.Address, timestamp time.Time) map[string]common.ListedToken {
	timestamps := timeutil.TimeToTimestampMs(timestamp)
	key := fmt.Sprintf("%s-%s", symbol, name)
	if token, existed := listedToken[key]; existed {
		if token.Timestamp > timestamps {
			token.Old = append(token.Old, common.OldListedToken{
				Address:   token.Address,
				Timestamp: token.Timestamp,
			})
			token.Address = address.Hex()
			token.Timestamp = timestamps
		} else {
			token.Old = append(token.Old, common.OldListedToken{
				Address:   address.Hex(),
				Timestamp: timestamps,
			})
		}
		return listedToken
	}
	listedToken[key] = common.ListedToken{
		Name:      name,
		Address:   address.Hex(),
		Symbol:    symbol,
		Timestamp: timestamps,
	}
	return listedToken
}

//GetListedToken return listed token for a reserve address
func (f *Fetcher) GetListedToken(block *big.Int, reserveAddr ethereum.Address,
	tokenSymbol *blockchain.TokenInfoGetter) error {
	var (
		logger = f.sugar.With("func", "accounting/cmd/accounting-listed-token-fetcher")
		result = make(map[string]common.ListedToken)
	)
	// step 1: get conversionRatesContract address
	logger.Infow("reserve address", "reserve", reserveAddr)
	reserveContractClient, err := contracts.NewReserve(reserveAddr, f.ethClient)
	if err != nil {
		return err
	}
	callOpts := &bind.CallOpts{BlockNumber: block}
	conversionRatesContract, err := reserveContractClient.ConversionRatesContract(callOpts)
	if err != nil {
		return err
	}

	// step 2: get listedTokens from conversionRatesContract
	conversionRateContractClient, err := contracts.NewConversionRates(conversionRatesContract, f.ethClient)
	if err != nil {
		return err
	}
	listedTokens, err := conversionRateContractClient.GetListedTokens(callOpts)
	if err != nil {
		return err
	}

	for _, address := range listedTokens {
		symbol, err := tokenSymbol.Symbol(address)
		if err != nil {
			return err
		}
		name, err := tokenSymbol.Name(address)
		if err != nil {
			return err
		}
		timestamp, err := f.contractTimestampResolver.Resolve(address)
		if err != nil {
			return err
		}
		result = updateListedToken(result, symbol, name, address, timestamp)
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}

	// currently print out to cli, save to storage later
	logger.Debugw("result listed token", "result", string(resultJSON))
	return nil
}
