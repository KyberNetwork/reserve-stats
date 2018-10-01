package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/KyberNetwork/reserve-stats/common"
	"github.com/KyberNetwork/reserve-stats/contracts"
	geth "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var errEmptyOutput = errors.New("blockchain return empty result")

type Blockchain struct {
	client  *ethclient.Client
	wrapper *contracts.WrapperContract
}

func NewBlockchain(endpoint string) (*Blockchain, error) {
	client, err := ethclient.Dial(common.InfuraEndpoint)
	if err != nil {
		return nil, err
	}
	wc, err := contracts.NewWrapperContract()
	if err != nil {
		return nil, err
	}
	return &Blockchain{client: client, wrapper: wc}, nil
}

// GetRatesFromOutput takes the output returned from call(), unpack it accordingly
// and do extra computation before returning the final ReserveRates or error if occurred.
func (bc *Blockchain) GetRatesFromOutput(output []byte, atBlock uint64, tokens []common.Token) (common.ReserveRates, error) {
	var (
		reserveRate = new([]*big.Int)
		sanityRate  = new([]*big.Int)
	)
	parsing := &[]interface{}{
		reserveRate,
		sanityRate,
	}

	rates := common.ReserveRates{}
	if err := bc.wrapper.ABI.Unpack(parsing, "getReserveRate", output); err != nil {
		return rates, err
	}

	rates.Timestamp = common.GetTimepoint()
	rates.BlockNumber = atBlock - 1
	rates.ToBlockNumber = atBlock
	rates.ReturnTime = common.GetTimepoint()
	result := common.ReserveTokenRateEntry{}

	for index, token := range tokens {
		rateEntry := common.ReserveRateEntry{}
		rateEntry.BuyReserveRate = common.BigToFloat((*reserveRate)[index*2+1], 18)
		rateEntry.BuySanityRate = common.BigToFloat((*sanityRate)[index*2+1], 18)
		rateEntry.SellReserveRate = common.BigToFloat((*reserveRate)[index*2], 18)
		rateEntry.SellSanityRate = common.BigToFloat((*sanityRate)[index*2], 18)
		result[fmt.Sprintf("ETH-%s", token.ID)] = rateEntry
	}
	rates.Data = result
	return rates, nil
}

func (bc *Blockchain) GetReserveRates(atBlock uint64, rsvAddress ethereum.Address, tokens []common.Token) (common.ReserveRates, error) {

	input, err := bc.wrapper.GetReserveRatesInput(atBlock, rsvAddress, tokens)
	if err != nil {
		return common.ReserveRates{}, err
	}
	wrapperAddr := contracts.WrapperContractAddr(atBlock)
	msg := geth.CallMsg{
		To:   &wrapperAddr,
		Data: input,
	}
	output, err := bc.client.CallContract(context.Background(), msg, big.NewInt(int64(atBlock)))
	if err != nil {
		return common.ReserveRates{}, err
	}
	if len(output) == 0 {
		return common.ReserveRates{}, errEmptyOutput
	}
	return bc.GetRatesFromOutput(output, atBlock, tokens)
}
