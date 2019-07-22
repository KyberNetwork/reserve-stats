package tradelogs

import (
	"context"
	"math/big"
	"time"

	ether "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	appname "github.com/KyberNetwork/reserve-stats/app-names"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/broadcast"
	"github.com/KyberNetwork/reserve-stats/lib/deployment"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/tokenrate"
)

const (
	// feeToWalletEvent is the topic of event AssignFeeToWallet(address reserve, address wallet, uint walletFee).
	feeToWalletEvent = "0x366bc34352215bf0bd3b527cfd6718605e1f5938777e42bcd8ed92f578368f52"

	// burnFeeEvent is the topic of event AssignBurnFees(address reserve, uint burnFee).
	burnFeeEvent = "0xf838f6ddc89706878e3c3e698e9b5cbfbf2c0e3d3dcd0bd2e00f1ccf313e0185"

	// executeTradeEvent is the topic of event
	// ExecuteTrade(address indexed sender, ERC20 src, ERC20 dest, uint actualSrcAmount, uint actualDestAmount).
	executeTradeEvent = "0x1849bd6a030a1bca28b83437fd3de96f3d27a5d172fa7e9c78e7b61468928a39"

	// etherReceivalEvent is the topic of event EtherReceival(address indexed sender, uint amount).
	etherReceivalEvent = "0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619"

	// kyberTradeEvent is the topic of event
	// KyberTrade (address trader, address src, address dest, uint256 srcAmount, uint256 dstAmount, address destAddress, uint256 ethWeiValue, address reserve1, address reserve2, bytes hint)
	// ETHwei = ETHAMOUNT
	// rsv1 ==0 if eth -> token
	// rsv2 ==0 if token -> eth
	kyberTradeEvent = "0xd30ca399cb43507ecec6a629a35cf45eb98cda550c27696dcb0d8c4a3873ce6c"

	// tradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress)
	// use for crawler v1 and v2
	tradeExecuteEvent = "0xea9415385bae08fe9f6dc457b02577166790cde83bb18cc340aac6cb81b824de"

	// tradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress, bool useInternalInventory)
	// overload function with useInternalInventory field
	tradeExecuteWithInternalEvent = "0x3e5691bfb48667c037c1014f0b320bbf4095c37ad5cd1f0b2573366756f54ac6"
)

var defaultTimeout = 10 * time.Second
var errUnknownLogTopic = errors.New("unknown log topic")

type tradeLogFetcher func(*big.Int, *big.Int, time.Duration) ([]common.TradeLog, error)

// NewCrawler create a new Crawler instance.
func NewCrawler(
	sugar *zap.SugaredLogger,
	client *ethclient.Client,
	broadcastClient broadcast.Interface,
	rateProvider tokenrate.ETHUSDRateProvider,
	addresses []ethereum.Address,
	sb deployment.VersionedStartingBlocks,
	etherscanClient *etherscan.Client) (*Crawler, error) {
	resolver, err := blockchain.NewBlockTimeResolver(sugar, client)
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
		etherscanClient: etherscanClient,
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
	startingBlocks  deployment.VersionedStartingBlocks

	etherscanClient *etherscan.Client
}

func logDataToExecuteTradeParams(data []byte) (ethereum.Address, ethereum.Address, ethereum.Hash, ethereum.Hash, error) {
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

func logDataToTradeExecuteParams(data []byte) (ethereum.Address, ethereum.Hash, ethereum.Address, ethereum.Hash, error) {
	var srcAddr, desAddr ethereum.Address
	var srcAmount, desAmount ethereum.Hash

	if len(data) < 128 {
		err := errors.New("invalid trade data")
		return srcAddr, srcAmount, desAddr, desAmount, err
	}

	srcAddr = ethereum.BytesToAddress(data[0:32])
	srcAmount = ethereum.BytesToHash(data[32:64])
	desAddr = ethereum.BytesToAddress(data[64:96])
	desAmount = ethereum.BytesToHash(data[96:128])
	return srcAddr, srcAmount, desAddr, desAmount, nil
}

func fillExecuteTrade(tradeLog common.TradeLog, logItem types.Log) (common.TradeLog, error) {
	srcAddr, destAddr, srcAmount, destAmount, err := logDataToExecuteTradeParams(logItem.Data)
	if err != nil {
		return common.TradeLog{}, err
	}
	tradeLog.SrcAddress = srcAddr
	tradeLog.DestAddress = destAddr
	tradeLog.SrcAmount = srcAmount.Big()
	tradeLog.DestAmount = destAmount.Big()

	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index
	tradeLog.UserAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	tradeLog.BlockNumber = logItem.BlockNumber

	return tradeLog, nil
}

//func logDataToEtherReceivalParams(data []byte) (ethereum.Hash, error) {
//	var amount ethereum.Hash
//
//	if len(data) != 32 {
//		err := errors.New("invalid eth receival data")
//		return amount, err
//	}
//
//	amount = ethereum.BytesToHash(data[0:32])
//	return amount, nil
//}

func fillEtherReceival(tradeLog common.TradeLog, logItem types.Log) (common.TradeLog, error) {
	tradeLog.SrcReserveAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	return tradeLog, nil
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

func fillWalletFees(tradeLog common.TradeLog, logItem types.Log) (common.TradeLog, error) {
	reserveAddr, walletAddr, fee, err := logDataToFeeWalletParams(logItem.Data)
	if err != nil {
		return common.TradeLog{}, err
	}

	walletFee := common.WalletFee{
		ReserveAddress: reserveAddr,
		WalletAddress:  walletAddr,
		Amount:         fee.Big(),
		Index:          logItem.Index,
	}
	tradeLog.WalletFees = append(tradeLog.WalletFees, walletFee)
	return tradeLog, nil
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

func fillBurnFees(tradeLog common.TradeLog, logItem types.Log) (common.TradeLog, error) {
	reserveAddr, fee, err := logDataToBurnFeeParams(logItem.Data)
	if err != nil {
		return common.TradeLog{}, err
	}

	burnFee := common.BurnFee{
		ReserveAddress: reserveAddr,
		Amount:         fee.Big(),
		Index:          logItem.Index,
	}
	tradeLog.BurnFees = append(tradeLog.BurnFees, burnFee)
	return tradeLog, nil
}

func logDataToKyberTradeV3Params(data []byte) (
	srcAddress, destAddress, srcReserve, dstReserve, receiverAddress ethereum.Address,
	srcAmount, destAmount, etherReceivalAmount ethereum.Hash,
	err error,
) {
	if len(data) < 256 {
		err = errors.New("invalid kyber trade v3 data")
		return
	}
	srcAddress = ethereum.BytesToAddress(data[0:32])
	destAddress = ethereum.BytesToAddress(data[32:64])
	srcReserve = ethereum.BytesToAddress(data[192:224])
	dstReserve = ethereum.BytesToAddress(data[224:256])
	receiverAddress = ethereum.BytesToAddress(data[128:160])
	srcAmount = ethereum.BytesToHash(data[64:96])
	destAmount = ethereum.BytesToHash(data[96:128])
	etherReceivalAmount = ethereum.BytesToHash(data[160:192])
	return
}

func fillKyberTradeV3(tradeLog common.TradeLog, logItem types.Log) (common.TradeLog, error) {
	srcAddress, destAddress, srcReserve, dstReserve, receiverAddress, srcAmount, destAmount, _, err := logDataToKyberTradeV3Params(logItem.Data)
	if err != nil {
		return common.TradeLog{}, err
	}
	tradeLog.SrcAddress = srcAddress
	tradeLog.DestAddress = destAddress
	tradeLog.SrcAmount = srcAmount.Big()
	tradeLog.DestAmount = destAmount.Big()
	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index
	tradeLog.UserAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	tradeLog.BlockNumber = logItem.BlockNumber
	tradeLog.SrcReserveAddress = srcReserve
	tradeLog.DstReserveAddress = dstReserve
	tradeLog.ReceiverAddress = receiverAddress

	return tradeLog, nil
}

func logDataToKyberTradeV2Params(data []byte) (ethereum.Address, ethereum.Address, ethereum.Address, ethereum.Address, ethereum.Hash, ethereum.Hash, error) {
	var srcAddr, desAddr, userAddr, receiverAddr ethereum.Address
	var srcAmount, desAmount ethereum.Hash

	if len(data) != 192 {
		err := errors.New("invalid kyber trade v2 data")
		return srcAddr, desAddr, userAddr, receiverAddr, srcAmount, desAmount, err
	}

	srcAddr = ethereum.BytesToAddress(data[32:64])
	desAddr = ethereum.BytesToAddress(data[128:160])
	srcAmount = ethereum.BytesToHash(data[64:96])
	desAmount = ethereum.BytesToHash(data[160:192])
	userAddr = ethereum.BytesToAddress(data[0:32])
	receiverAddr = ethereum.BytesToAddress(data[96:128])
	return srcAddr, desAddr, userAddr, receiverAddr, srcAmount, desAmount, nil
}

func fillKyberTradeV2(tradeLog common.TradeLog, logItem types.Log) (common.TradeLog, error) {
	srcAddr, destAddr, userAddr, receiverAddr, srcAmount, destAmount, err := logDataToKyberTradeV2Params(logItem.Data)
	if err != nil {
		return common.TradeLog{}, err
	}
	tradeLog.SrcAddress = srcAddr
	tradeLog.DestAddress = destAddr

	tradeLog.SrcAmount = srcAmount.Big()
	tradeLog.DestAmount = destAmount.Big()

	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index

	tradeLog.UserAddress = userAddr
	tradeLog.ReceiverAddress = receiverAddr

	tradeLog.BlockNumber = logItem.BlockNumber

	return tradeLog, nil
}

func assembleTradeLogsReserveAddr(log common.TradeLog, sugar *zap.SugaredLogger) common.TradeLog {
	switch {
	case blockchain.IsBurnable(log.SrcAddress):
		if blockchain.IsBurnable(log.DestAddress) {
			if len(log.BurnFees) == 2 {
				log.SrcReserveAddress = log.BurnFees[0].ReserveAddress
				log.DstReserveAddress = log.BurnFees[1].ReserveAddress
			} else {
				sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "2 burn fees (src-dst)")
			}
		} else {
			if len(log.BurnFees) == 1 {
				log.SrcReserveAddress = log.BurnFees[0].ReserveAddress
			} else {
				sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (src)")
			}
		}
	case blockchain.IsBurnable(log.DestAddress):
		if len(log.BurnFees) == 1 {
			log.DstReserveAddress = log.BurnFees[0].ReserveAddress
		} else {
			sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (dst)")
		}
	case len(log.WalletFees) != 0:
		if len(log.WalletFees) == 1 {
			log.SrcReserveAddress = log.WalletFees[0].ReserveAddress
		} else {
			log.SrcReserveAddress = log.WalletFees[0].ReserveAddress
			log.DstReserveAddress = log.WalletFees[1].ReserveAddress
		}
	}

	return log
}

func (crawler *Crawler) fetchLogsWithTopics(fromBlock, toBlock *big.Int, timeout time.Duration, topics [][]ethereum.Hash) ([]types.Log, error) {
	query := ether.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: crawler.addresses,
		Topics:    topics,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return crawler.ethClient.FilterLogs(ctx, query)

}

func (crawler *Crawler) getTxSender(log types.Log, timeout time.Duration) (ethereum.Address, error) {
	var txSender ethereum.Address
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	tx, _, err := crawler.ethClient.TransactionByHash(ctx, log.TxHash)
	if err != nil {
		return txSender, err
	}
	txSender, err = crawler.ethClient.TransactionSender(ctx, tx, log.BlockHash, log.TxIndex)
	return txSender, err
}

// function for v2, v3
func getEthAmountFromReceipt(receipt *types.Receipt) (*big.Int, error) {
	var ethAmount = big.NewInt(0)
	for index := uint64(len(receipt.Logs) - 1); index > 0; index-- {
		log := receipt.Logs[index]
		for _, topic := range log.Topics {
			if topic == ethereum.HexToHash(tradeExecuteEvent) || topic == ethereum.HexToHash(tradeExecuteWithInternalEvent) {
				srcAddr, srcAmount, dstAddr, dstAmount, err := logDataToTradeExecuteParams(log.Data)
				if err != nil {
					return nil, err
				}
				var addedEthAmount *big.Int
				if srcAddr == blockchain.ETHAddr {
					addedEthAmount = srcAmount.Big()
				} else {
					if dstAddr == blockchain.ETHAddr {
						addedEthAmount = dstAmount.Big()
					} else {
						return nil, errors.Errorf("failed to take eth amount src_addr = %v dst_addr = %v", srcAddr.Hex(), dstAddr.Hex())
					}
				}
				ethAmount.Add(ethAmount, addedEthAmount)
			}
		}
	}
	return ethAmount, nil
}

// GetTradeLogs returns trade logs from KyberNetwork.
func (crawler *Crawler) GetTradeLogs(fromBlock, toBlock *big.Int, timeout time.Duration) ([]common.TradeLog, error) {
	var (
		result  []common.TradeLog
		fetchFn tradeLogFetcher
	)

	// fetchTradeLogV2 also works for v3 trades, so to keep it simple, we only use fetchTradeLogV3 if both
	// from, to blocks are >= starting block v3
	switch {
	case fromBlock.Uint64() >= crawler.startingBlocks.V3() && toBlock.Uint64() >= crawler.startingBlocks.V3():
		fetchFn = crawler.fetchTradeLogV3
	case fromBlock.Uint64() >= crawler.startingBlocks.V2() && toBlock.Uint64() >= crawler.startingBlocks.V2():
		fetchFn = crawler.fetchTradeLogV2
	default:
		fetchFn = crawler.fetchTradeLogV1
	}

	result, err := fetchFn(fromBlock, toBlock, timeout)
	if err != nil {
		return result, err
	}

	for i, tradeLog := range result {
		var uid, ip, country string

		uid, ip, country, err = crawler.broadcastClient.GetTxInfo(tradeLog.TransactionHash.Hex())
		if err != nil {
			return result, err
		}
		result[i].IP = ip
		result[i].Country = country
		result[i].UID = uid

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
