package tradelogs

import (
	"errors"
	"github.com/KyberNetwork/tokenrate"
	"math/big"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"

	appname "github.com/KyberNetwork/reserve-stats/app-names"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
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
	// kyberTradeEvent is the topic of event
	// KyberTrade (address trader, address src, address dest, uint256 srcAmount, uint256 dstAmount, address destAddress, uint256 ethWeiValue, address reserve1, address reserve2, bytes hint)
	// ETHwei = ETHAMOUNT
	// rsv1 ==0 if eth -> token
	// rsv2 ==0 if token -> eth
	kyberTradeEvent = "0xd30ca399cb43507ecec6a629a35cf45eb98cda550c27696dcb0d8c4a3873ce6c"
)

// NewCrawler create a new Crawler instance.
func NewCrawler(
	sugar *zap.SugaredLogger,
	client *ethclient.Client,
	broadcastClient broadcast.Interface,
	rateProvider tokenrate.ETHUSDRateProvider,
	addresses []ethereum.Address,
	sb app.VersionedStartingBlocks) (*Crawler, error) {
	resolver, err := blockchain.NewBlockTimeResolver(sugar, client)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &Crawler{
		sugar:           sugar,
		ethClient:       client,
		txTime:          resolver,
		broadcastClient: broadcastClient,
		rateProvider:    rateProvider,
		addresses:       addresses,
		startingBlocks:  sb,
	}, nil
}

// Crawler gets trade logs on KyberNetwork on blockchain, adding the
// information about USD equivalent on each trade.
type Crawler struct {
	sugar           *zap.SugaredLogger
	ethClient       *ethclient.Client
	txTime          *blockchain.BlockTimeResolver
	broadcastClient broadcast.Interface
	rateProvider    tokenrate.ETHUSDRateProvider
	addresses       []ethereum.Address
	startingBlocks  app.VersionedStartingBlocks
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

func (crawler *Crawler) assembleTradeLogs(eventLogs []types.Log) ([]common.TradeLog, error) {
	var (
		result   []common.TradeLog
		tradeLog common.TradeLog
	)

	for _, log := range eventLogs {
		if log.Removed {
			continue // Removed due to chain reorg
		}

		if len(log.Topics) == 0 {
			return result, errors.New("log item has no topic")
		}

		topic := log.Topics[0]
		switch topic.Hex() {
		case feeToWalletEvent:
			reserveAddr, walletAddr, fee, err := logDataToFeeWalletParams(log.Data)
			if err != nil {
				return nil, err
			}

			walletFee := common.WalletFee{
				ReserveAddress: reserveAddr,
				WalletAddress:  walletAddr,
				Amount:         fee.Big(),
				Index:          log.Index,
			}
			tradeLog.WalletFees = append(tradeLog.WalletFees, walletFee)
		case burnFeeEvent:
			reserveAddr, fee, err := logDataToBurnFeeParams(log.Data)
			if err != nil {
				return nil, err
			}

			burnFee := common.BurnFee{
				ReserveAddress: reserveAddr,
				Amount:         fee.Big(),
				Index:          log.Index,
			}
			tradeLog.BurnFees = append(tradeLog.BurnFees, burnFee)
		default:
			done, err := crawler.handleEventLogWithBlockNumber(log, &tradeLog)
			if err != nil {
				return nil, err
			}
			if done {
				result = append(result, tradeLog)
				tradeLog = common.TradeLog{}
			}
		}
	}

	return result, nil
}

// GetTradeLogs returns trade logs from KyberNetwork.
func (crawler *Crawler) GetTradeLogs(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
	var (
		result []common.TradeLog
	)
	if fromBlock.Uint64() < crawler.startingBlocks.V3 {
		if toBlock.Uint64() < crawler.startingBlocks.V3 {
			tradeLogs, err := crawler.fetchTradeLogV2(fromBlock, toBlock, timeout)
			if err != nil {
				return result, err
			}
			result = append(result, tradeLogs...)
			return result, nil
		}

		tradeLogs, err := crawler.fetchTradeLogV2(fromBlock, big.NewInt(int64(crawler.startingBlocks.V3-1)), timeout)
		if err != nil {
			return result, err
		}
		result = append(result, tradeLogs...)
		tradeLogs, err = crawler.fetchTradeLogV3(big.NewInt(int64(crawler.startingBlocks.V3)), toBlock, timeout)
		if err != nil {
			return result, err
		}
		result = append(result, tradeLogs...)
		return result, nil
	}

	result, err := crawler.fetchTradeLogV3(fromBlock, toBlock, timeout)
	if err != nil {
		return result, err
	}

	for i, tradeLog := range result {
		var (
			ip, country string
		)

		ip, country, err = crawler.broadcastClient.GetTxInfo(tradeLog.TransactionHash.Hex())
		if err != nil {
			return result, err
		}
		result[i].IP = ip
		result[i].Country = country
		if tradeLog.IsKyberSwap() {
			result[i].IntegrationApp = appname.KyberSwapAppName
		} else {
			result[i].IntegrationApp = appname.ThirdPartyAppName
		}
		rate, err := crawler.rateProvider.USDRate(tradeLog.Timestamp)
		if err != nil {
			return nil, err
		}
		result[i].ETHUSDProvider = crawler.rateProvider.Name()
		result[i].ETHUSDRate = rate
	}
	return result, nil
}
