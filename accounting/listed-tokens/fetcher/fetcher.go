package fetcher

import (
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

func updateListedToken(listedToken map[string]common.ListedToken, symbol, name string, address ethereum.Address, timestamp time.Time,
	decimals uint8) map[string]common.ListedToken {
	key := fmt.Sprintf("%s-%s", symbol, name)
	if token, existed := listedToken[key]; existed {
		if token.Timestamp.After(timestamp) {
			token.Old = append(token.Old, common.OldListedToken{
				Address:   token.Address,
				Timestamp: token.Timestamp,
				Decimals:  token.Decimals,
			})
			token.Address = address
			token.Timestamp = timestamp
			token.Decimals = decimals
			listedToken[key] = token
		} else {
			token.Old = append(token.Old, common.OldListedToken{
				Address:   address,
				Timestamp: timestamp,
				Decimals:  decimals,
			})
			listedToken[key] = token
		}
		return listedToken
	}
	listedToken[key] = common.ListedToken{
		Name:      name,
		Address:   address,
		Symbol:    symbol,
		Timestamp: timestamp,
		Decimals:  decimals,
	}
	return listedToken
}

//GetListedToken return listed token for a reserve address
func (f *Fetcher) GetListedToken(block *big.Int, reserveAddr ethereum.Address,
	tokenSymbol *blockchain.TokenInfoGetter) ([]common.ListedToken, error) {
	var (
		logger       = f.sugar.With("func", "accounting/cmd/accounting-listed-token-fetcher")
		result       = make(map[string]common.ListedToken)
		returnResult []common.ListedToken
	)
	// step 1: get conversionRatesContract address
	logger.Infow("reserve address", "reserve", reserveAddr)
	reserveContractClient, err := contracts.NewReserve(reserveAddr, f.ethClient)
	if err != nil {
		return nil, err
	}
	callOpts := &bind.CallOpts{BlockNumber: block}
	conversionRatesContract, err := reserveContractClient.ConversionRatesContract(callOpts)
	if err != nil {
		return nil, err
	}

	// step 2: get listedTokens from conversionRatesContract
	conversionRateContractClient, err := contracts.NewConversionRates(conversionRatesContract, f.ethClient)
	if err != nil {
		return nil, err
	}
	listedTokens, err := conversionRateContractClient.GetListedTokens(callOpts)
	if err != nil {
		return nil, err
	}

	for _, address := range listedTokens {
		symbol, err := tokenSymbol.Symbol(address)
		if err != nil {
			return nil, err
		}
		name, err := tokenSymbol.Name(address)
		if err != nil {
			return nil, err
		}
		timestamp, err := f.contractTimestampResolver.Resolve(address)
		if err != nil {
			return nil, err
		}
		decimals, err := tokenSymbol.Decimals(address)
		if err != nil {
			return nil, err
		}
		result = updateListedToken(result, symbol, name, address, timestamp, decimals)
	}
	for _, token := range result {
		returnResult = append(returnResult, token)
	}

	return returnResult, nil
}
