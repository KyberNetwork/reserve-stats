package blockchain

import (
	"context"
	"errors"
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

func (bc *Blockchain) GetReserveRates(atBlock uint64, rsvAddress ethereum.Address, tokens []common.Token) (*common.ReserveRates, error) {
	result := common.ReserveRates{}
	rates := common.ReserveRates{}
	rates.Timestamp = common.GetTimepoint()

	input, err := bc.wrapper.GetReserveRatesInput(atBlock, rsvAddress, tokens)
	if err != nil {
		return nil, err
	}
	wrapperAddr := contracts.WrapperContractAddr(atBlock)
	msg := geth.CallMsg{
		To:   &wrapperAddr,
		Data: input,
	}
	output, err := bc.client.CallContract(context.Background(), msg, big.NewInt(int64(atBlock)))
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return nil, errEmptyOutput
	}
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	parsed := &[]interface{}{
		ret0,
		ret1,
	}
	bc.wrapper.ABI.Unpack(parsed, "getReserveRate", output)

	return &result, nil
}
