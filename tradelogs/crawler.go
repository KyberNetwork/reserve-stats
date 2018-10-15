package tradelogs

import (
	"context"
	"errors"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/tokenrate"
	ether "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

const (
	// feeToWalletEvent is the topic of event AssignFeeToWallet(address reserve, address wallet, uint walletFee).
	feeToWalletEvent = "0x366bc34352215bf0bd3b527cfd6718605e1f5938777e42bcd8ed92f578368f52"
	// burnFeeEvent is the topic of event AssignBurnFees(address reserve, uint burnFee).
	burnFeeEvent = "0xf838f6ddc89706878e3c3e698e9b5cbfbf2c0e3d3dcd0bd2e00f1ccf313e0185"
	// tradeEvent is the topic of event
	// ExecuteTrade(address indexed sender, ERC20 src, ERC20 dest, uint actualSrcAmount, uint actualDestAmount).
	// tradeEvent is the topic of event
	// ExecuteTrade(address indexed sender, ERC20 src, ERC20 dest, uint actualSrcAmount, uint actualDestAmount).
	tradeEvent = "0x1849bd6a030a1bca28b83437fd3de96f3d27a5d172fa7e9c78e7b61468928a39"
	// etherReceivalEvent is the topic of event EtherReceival(address indexed sender, uint amount).
	etherReceivalEvent = "0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619"

	// address of pricing contract
	pricingAddr = "0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B"
	// address of network contract
	networkAddr = "0x818E6FECD516Ecc3849DAf6845e3EC868087B755"
	// address of bunner contract
	burnerAddr = "0xed4f53268bfdFF39B36E8786247bA3A02Cf34B04"
	// address of internal network contract
	internalNetworkAddr = "0x91a502C678605fbCe581eae053319747482276b9"

	ethDecimals int64  = 18
	ethAddress  string = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

// TradeLogCrawler gets trade logs on KyberNetwork on blockchain, adding the
// information about USD equivalent on each trade.
type TradeLogCrawler struct {
	sugar        *zap.SugaredLogger
	ethClient    *ethclient.Client
	rateProvider tokenrate.ETHUSDRateProvider
	txTime       *blockchain.BlockTimeResolver
}

func logDataToTradeParams(data []byte) (ethereum.Address, ethereum.Address, ethereum.Hash, ethereum.Hash, error) {
	var srcAddr, desAddr ethereum.Address
	var srcAmount, desAmount ethereum.Hash

	if len(data) != 128 {
		err := errors.New("invalid trade data")
		return srcAddr, desAddr, srcAmount, desAmount, err
	}

	srcAddr = ethereum.BytesToAddress(data[0:32])
	desAddr = ethereum.BytesToAddress(data[32:64])
	srcAmount = ethereum.BytesToHash(data[64:96])
	desAmount = ethereum.BytesToHash(data[96:128])
	return srcAddr, desAddr, srcAmount, desAmount, nil
}

func logDataToEtherReceivalParams(data []byte) (ethereum.Hash, error) {
	var amount ethereum.Hash

	if len(data) != 32 {
		err := errors.New("invalid eth receival data")
		return amount, err
	}

	amount = ethereum.BytesToHash(data[0:32])
	return amount, nil
}

func logDataToFeeWalletParams(data []byte) (ethereum.Address, ethereum.Address, ethereum.Hash, error) {
	var reserveAddr, walletAddr ethereum.Address
	var walletFee ethereum.Hash

	if len(data) != 96 {
		err := errors.New("invalid fee wallet data")
		return reserveAddr, walletAddr, walletFee, err
	}

	reserveAddr = ethereum.BytesToAddress(data[0:32])
	walletAddr = ethereum.BytesToAddress(data[32:64])
	walletFee = ethereum.BytesToHash(data[64:96])
	return reserveAddr, walletAddr, walletFee, nil
}

func logDataToBurnFeeParams(data []byte) (ethereum.Address, ethereum.Hash, error) {
	var reserveAddr ethereum.Address
	var burnFees ethereum.Hash

	if len(data) != 64 {
		err := errors.New("invalid burn fee data")
		return reserveAddr, burnFees, err
	}

	reserveAddr = ethereum.BytesToAddress(data[0:32])
	burnFees = ethereum.BytesToHash(data[32:64])
	return reserveAddr, burnFees, nil
}

func updateTradeLogs(allLogs []common.TradeLog, logItem types.Log, ts time.Time) ([]common.TradeLog, error) {
	var (
		tradeLog      common.TradeLog
		updateLastLog = false
	)

	if len(logItem.Topics) == 0 {
		return allLogs, errors.New("log item has no topic")
	}

	// if the transaction hash is the same with last log is a TradeLog, update it,
	// otherwise append new one
	if len(allLogs) > 0 && allLogs[len(allLogs)-1].TransactionHash == logItem.TxHash {
		tradeLog = allLogs[len(allLogs)-1]
		updateLastLog = true
	}

	if !updateLastLog {
		tradeLog = common.TradeLog{
			BlockNumber:     logItem.BlockNumber,
			TransactionHash: logItem.TxHash,
			Index:           logItem.Index,
			Timestamp:       ts,
		}
	}

	switch logItem.Topics[0].Hex() {
	case feeToWalletEvent:
		reserveAddr, walletAddr, fee, err := logDataToFeeWalletParams(logItem.Data)
		if err != nil {
			return allLogs, err
		}

		walletFee := common.WalletFee{
			ReserveAddress: reserveAddr,
			WalletAddress:  walletAddr,
			Amount:         fee.Big(),
		}
		tradeLog.WalletFees = append(tradeLog.WalletFees, walletFee)
	case burnFeeEvent:
		reserveAddr, fee, err := logDataToBurnFeeParams(logItem.Data)
		if err != nil {
			return allLogs, err
		}
		burnFee := common.BurnFee{
			ReserveAddress: reserveAddr,
			Amount:         fee.Big(),
		}
		tradeLog.BurnFees = append(tradeLog.BurnFees, burnFee)
	case etherReceivalEvent:
		amount, err := logDataToEtherReceivalParams(logItem.Data)
		if err != nil {
			return allLogs, err
		}
		tradeLog.EtherReceivalSender = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
		tradeLog.EtherReceivalAmount = amount.Big()
	case tradeEvent:
		srcAddr, destAddr, srcAmount, destAmount, err := logDataToTradeParams(logItem.Data)
		if err != nil {
			return allLogs, err
		}
		tradeLog.SrcAddress = srcAddr
		tradeLog.DestAddress = destAddr
		tradeLog.SrcAmount = srcAmount.Big()
		tradeLog.DestAmount = destAmount.Big()
		tradeLog.UserAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	}

	if updateLastLog {
		allLogs[len(allLogs)-1] = tradeLog
	} else {
		allLogs = append(allLogs, tradeLog)
	}

	return allLogs, nil
}

// calculateFiatAmount returns new TradeLog with fiat amount calculated.
// * For ETH-Token or Token-ETH conversions, the ETH amount is taken from ExecuteTrade event.
// * For Token-Token, the ETH amount is reading from EtherReceival event.
func calculateFiatAmount(tradeLog common.TradeLog, rate float64) common.TradeLog {
	ethAmount := new(big.Float)

	if strings.ToLower(ethAddress) == strings.ToLower(tradeLog.SrcAddress.String()) {
		// ETH-Token
		ethAmount.SetInt(tradeLog.SrcAmount)
	} else if strings.ToLower(ethAddress) == strings.ToLower(tradeLog.DestAddress.String()) {
		// Token-ETH
		ethAmount.SetInt(tradeLog.DestAmount)
	} else if tradeLog.EtherReceivalAmount != nil {
		// Token-Token
		ethAmount.SetInt(tradeLog.EtherReceivalAmount)
	}

	// fiat amount = ETH amount * rate
	ethAmount = ethAmount.Mul(ethAmount, new(big.Float).SetFloat64(rate))
	ethAmount.Quo(ethAmount, new(big.Float).SetFloat64(math.Pow10(int(ethDecimals))))
	tradeLog.FiatAmount, _ = ethAmount.Float64()

	return tradeLog
}

// NewTradeLogCrawler create a new TradeLogCrawler instance.
func NewTradeLogCrawler(sugar *zap.SugaredLogger, nodeURL string, rateProvider tokenrate.ETHUSDRateProvider) (*TradeLogCrawler, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, err
	}
	resolver, err := blockchain.NewBlockTimeResolver(sugar, client)
	return &TradeLogCrawler{sugar, client, rateProvider, resolver}, nil
}

// GetTradeLogs returns trade logs from KyberNetwork.
func (crawler *TradeLogCrawler) GetTradeLogs(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
	var result []common.TradeLog

	addresses := []ethereum.Address{
		ethereum.HexToAddress(pricingAddr),         // pricing
		ethereum.HexToAddress(networkAddr),         // network
		ethereum.HexToAddress(burnerAddr),          // burner
		ethereum.HexToAddress(internalNetworkAddr), // internal network
	}

	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(tradeEvent),
			ethereum.HexToHash(burnFeeEvent),
			ethereum.HexToHash(feeToWalletEvent),
			ethereum.HexToHash(etherReceivalEvent),
		},
	}

	query := ether.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: addresses,
		Topics:    topics,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logs, err := crawler.ethClient.FilterLogs(ctx, query)
	if err != nil {
		return result, err
	}

	for _, logItem := range logs {
		if logItem.Removed {
			continue // Removed due to chain reorg
		}

		ts, err := crawler.txTime.Resolve(logItem.BlockNumber)

		if err != nil {
			return result, err
		}

		topic := logItem.Topics[0]
		switch topic.Hex() {
		case feeToWalletEvent, burnFeeEvent, etherReceivalEvent, tradeEvent:
			// add logItem to result
			if result, err = updateTradeLogs(result, logItem, ts); err != nil {
				return result, err
			}
		default:
			crawler.sugar.Info("Unknown log topic.")
		}
	}

	for i, tradeLog := range result {
		var ethRate float64
		if ethRate, err = crawler.rateProvider.USDRate(tradeLog.Timestamp); err != nil {
			crawler.sugar.Errorw("failed to get ETH/USD rate",
				"timestamp", tradeLog.Timestamp.String())
			return nil, err
		}
		if ethRate != 0 {
			crawler.sugar.Debugw("got ETH/USD rate",
				"rate", ethRate,
				"timestamp", tradeLog.Timestamp.String())
			result[i] = calculateFiatAmount(tradeLog, ethRate)
		}
	}

	return result, nil
}
