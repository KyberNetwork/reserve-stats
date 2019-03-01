package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
)

const (
	blockFlag           = "block"
	reserveAddresesFlag = "reserve-address"
)

func main() {
	app := libapp.NewApp()
	app.Name = "accounting-listed-token-fetcher"
	app.Usage = "get listed token for provided block and reserve address"
	app.Action = run
	app.Flags = append(app.Flags,
		cli.StringFlag{
			Name:   blockFlag,
			EnvVar: "BLOCK",
			Usage:  "block to get listed token",
		},
		cli.StringFlag{
			Name:   reserveAddresesFlag,
			EnvVar: "RESERVE_ADDRESS",
			Usage:  "reserve address to get listed token",
		},
	)
	app.Flags = append(app.Flags, blockchain.NewEthereumNodeFlags())
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	var (
		block       *big.Int
		reserveAddr ethereum.Address
	)
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	if c.String(blockFlag) == "" {
		sugar.Info("no block number provided, get listed token from latest block")
	} else {
		block, err = libapp.ParseBigIntFlag(c, blockFlag)
		if err != nil {
			return err
		}
	}
	if c.String(reserveAddresesFlag) == "" {
		return fmt.Errorf("reserve address is required")
	}
	reserveAddrStr := c.String(reserveAddresesFlag)
	reserveAddr = ethereum.HexToAddress(reserveAddrStr)

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return err
	}
	tokenSymbol, err := blockchain.NewTokenSymbolFromContext(c)
	if err != nil {
		return err
	}

	return getListedToken(ethClient, block, reserveAddr, tokenSymbol, sugar)
}

func getListedToken(ethClient *ethclient.Client, block *big.Int, reserveAddr ethereum.Address,
	tokenSymbol *blockchain.TokenSymbol, sugar *zap.SugaredLogger) error {
	var (
		logger = sugar.With("func", "accounting/cmd/accounting-listed-token-fetcher")
		result []common.ListedToken
	)
	// step 1: get conversionRatesContract address
	logger.Infow("reserve address", "reserve", reserveAddr)
	reserveContractClient, err := contracts.NewReserve(reserveAddr, ethClient)
	if err != nil {
		return err
	}
	callOpts := &bind.CallOpts{BlockNumber: block}
	conversionRatesContract, err := reserveContractClient.ConversionRatesContract(callOpts)
	if err != nil {
		return err
	}
	// step 2: get listedTokens from conversionRatesContract
	conversionRateContractClient, err := contracts.NewConversionRates(conversionRatesContract, ethClient)
	if err != nil {
		return err
	}
	listedTokens, err := conversionRateContractClient.GetListedTokens(callOpts)
	if err != nil {
		return err
	}
	for _, v := range listedTokens {
		symbol, err := tokenSymbol.Symbol(v)
		if err != nil {
			return err
		}
		name, err := tokenSymbol.Name(v)
		if err != nil {
			return err
		}
		result = append(result, common.ListedToken{
			Address: v.Hex(),
			Symbol:  symbol,
			Name:    name,
		})
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}
	// currently print out to cli, save to storage later
	log.Printf("%s", resultJSON)
	// step 3: use api etherscan to get first transaction timestamp
	//TODO: wait for favadi task to finish
	return nil
}
