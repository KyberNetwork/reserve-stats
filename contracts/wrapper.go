package contracts

import (
	"bytes"

	"github.com/KyberNetwork/reserve-stats/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethereum "github.com/ethereum/go-ethereum/common"
)

const wrapperABI = `[{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getInt8FromByte","outputs":[{"name":"","type":"int8"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"reserve","type":"address"},{"name":"tokens","type":"address[]"}],"name":"getBalances","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"ratesContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenIndicies","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"reserve","type":"address"},{"name":"srcs","type":"address[]"},{"name":"dests","type":"address[]"}],"name":"getReserveRate","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"x","type":"bytes14"},{"name":"byteInd","type":"uint256"}],"name":"getByteFromBytes14","outputs":[{"name":"","type":"bytes1"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"network","type":"address"},{"name":"srcs","type":"address[]"},{"name":"dests","type":"address[]"},{"name":"qty","type":"uint256[]"}],"name":"getExpectedRates","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"ratesContract","type":"address"},{"name":"tokenList","type":"address[]"}],"name":"getTokenRates","outputs":[{"name":"","type":"uint256[]"},{"name":"","type":"uint256[]"},{"name":"","type":"int8[]"},{"name":"","type":"int8[]"},{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"}]`

type WrapperContract struct {
	ABI abi.ABI
}

func NewWrapperContract() (*WrapperContract, error) {
	contractABI, err := abi.JSON(bytes.NewReader([]byte(wrapperABI)))
	if err != nil {
		return nil, err
	}
	return &WrapperContract{
		ABI: contractABI,
	}, nil
}

func (wc *WrapperContract) GetReserveRatesInput(block uint64, rsvAddress ethereum.Address, tokens []common.Token) ([]byte, error) {
	method := "getReserveRate"
	ETH := common.ETHToken
	srcAddresses := []ethereum.Address{}
	destAddresses := []ethereum.Address{}
	for _, token := range tokens {
		srcAddresses = append(srcAddresses, ethereum.HexToAddress(token.Address), ethereum.HexToAddress(ETH.Address))
		destAddresses = append(destAddresses, ethereum.HexToAddress(ETH.Address), ethereum.HexToAddress(token.Address))
	}
	return wc.ABI.Pack(method, rsvAddress, srcAddresses, destAddresses)
}

// cWrapperontractAddrs returns the proper network, contract addresses to use with given block number.
func WrapperContractAddr(block uint64) (wrapperAddr ethereum.Address) {
	if block < common.StartingBlockV2 {
		wrapperAddr = ethereum.HexToAddress(common.WrapperAddrV1)
	} else {
		wrapperAddr = ethereum.HexToAddress(common.WrapperAddrV2)
	}
	return
}
