package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/KyberNetwork/reserve-stats/tradelogs"
)

func main() {
	// trade-logs-crawler --addresses=0xABCDEF,0xDEFGHI --from-block=100 -to-block=200

	fromBlock := uint64(6395977)
	toBlock := uint64(6396307)

	//set client & endpoint
	endpoint := "https://mainnet.infura.io"
	client, err := rpc.Dial(endpoint)
	if err != nil {
		panic(err)
	}
	ethClient := ethclient.NewClient(client)
	ethRate := tradelogs.NewCMCEthUSDRate()

	crawler := tradelogs.NewTradeLogCrawler(ethClient, ethRate)
	tradeLogs, err := crawler.GetTradeLogs(fromBlock, toBlock)

	if err != nil {
		log.Fatal(err)
	}

	for _, tradeLog := range tradeLogs {
		fmt.Printf("%+v\n", tradeLog)
	}
}
