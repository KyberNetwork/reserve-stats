package tradedetect

import (
	"context"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	//TradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress)
	tradeExecuteEvent = "0xea9415385bae08fe9f6dc457b02577166790cde83bb18cc340aac6cb81b824de"
	timeout           = 10 * time.Second
)

//DetectTradeTransaction detect if a provided txHash is belong to a trade transaction or not
func DetectTradeTransaction(txHash ethereum.Hash, ethClient *ethclient.Client) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	receipt, err := ethClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return false, err
	}
	for _, log := range receipt.Logs {
		for _, topic := range log.Topics {
			if topic == ethereum.HexToHash(tradeExecuteEvent) {
				return true, nil
			}
		}
	}
	return false, nil
}
