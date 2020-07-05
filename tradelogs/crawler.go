package tradelogs

import (
	"bytes"
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
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/contracts"
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
)

var defaultTimeout = 10 * time.Second
var errUnknownLogTopic = errors.New("unknown log topic")

type tradeLogFetcher func(*big.Int, *big.Int, time.Duration) (*common.CrawlResult, error)

// NewCrawler create a new Crawler instance.
func NewCrawler(sugar *zap.SugaredLogger,
	client *ethclient.Client,
	broadcastClient broadcast.Interface,
	rateProvider tokenrate.ETHUSDRateProvider,
	addresses []ethereum.Address,
	sb deployment.VersionedStartingBlocks,
	etherscanClient *etherscan.Client,
	volumeExcludedReserves []ethereum.Address,
	networkProxy, kyberStorage, kyberFeeHandler, kyberNetwork ethereum.Address) (*Crawler, error) {
	resolver, err := blockchain.NewBlockTimeResolver(sugar, client)
	if err != nil {
		return nil, err
	}

	kyberStorageContract, err := contracts.NewKyberStorage(kyberStorage, client)
	if err != nil {
		return nil, err
	}

	kyberFeeHandlerContract, err := contracts.NewKyberFeeHandler(kyberFeeHandler, client)
	if err != nil {
		return nil, err
	}

	kyberNetworkContract, err := contracts.NewKyberNetwork(kyberNetwork, client)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		sugar:                   sugar,
		ethClient:               client,
		txTime:                  resolver,
		broadcastClient:         broadcastClient,
		rateProvider:            rateProvider,
		addresses:               addresses,
		startingBlocks:          sb,
		etherscanClient:         etherscanClient,
		volumeExludedReserves:   volumeExcludedReserves,
		networkProxy:            networkProxy,
		kyberStorageContract:    kyberStorageContract,
		kyberFeeHandlerContract: kyberFeeHandlerContract,
		kyberNetworkContract:    kyberNetworkContract,
	}, nil
}

// Crawler gets trade logs on KyberNetwork on blockchain, adding the
// information about USD equivalent on each trade.
type Crawler struct {
	sugar                 *zap.SugaredLogger
	ethClient             *ethclient.Client
	txTime                *blockchain.BlockTimeResolver
	broadcastClient       broadcast.Interface
	rateProvider          tokenrate.ETHUSDRateProvider
	addresses             []ethereum.Address
	startingBlocks        deployment.VersionedStartingBlocks
	volumeExludedReserves []ethereum.Address

	kyberStorageContract    *contracts.KyberStorage
	kyberFeeHandlerContract *contracts.KyberFeeHandler
	kyberNetworkContract    *contracts.KyberNetwork

	etherscanClient *etherscan.Client
	networkProxy    ethereum.Address
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

func fillExecuteTrade(tradeLog common.TradelogV4, logItem types.Log) (common.TradelogV4, error) {
	srcAddr, destAddr, srcAmount, destAmount, err := logDataToExecuteTradeParams(logItem.Data)
	if err != nil {
		return common.TradelogV4{}, err
	}
	tradeLog.TokenInfo.SrcAddress = srcAddr
	tradeLog.TokenInfo.DestAddress = destAddr
	tradeLog.SrcAmount = srcAmount.Big()
	tradeLog.DestAmount = destAmount.Big()

	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index
	tradeLog.User.UserAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	tradeLog.BlockNumber = logItem.BlockNumber

	return tradeLog, nil
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

func fillEtherReceival(tradeLog common.TradelogV4, logItem types.Log) (common.TradelogV4, error) {
	amount, err := logDataToEtherReceivalParams(logItem.Data)
	if err != nil {
		return tradeLog, err
	}
	// tradeLog.T2EReserves = append(tradeLog.T2EReserves, ethereum.BytesToAddress(logItem.Topics[1].Bytes()))
	tradeLog.SrcReserveAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	tradeLog.EthAmount = amount.Big()
	tradeLog.OriginalEthAmount = amount.Big()
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

func fillWalletFees(tradeLog common.TradelogV4, logItem types.Log) (common.TradelogV4, error) {
	reserveAddr, walletAddr, fee, err := logDataToFeeWalletParams(logItem.Data)
	if err != nil {
		return common.TradelogV4{}, err
	}

	tradelogFee := common.TradelogFee{
		WalletFee:      fee.Big(),
		PlatformWallet: walletAddr,
		WalletName:     WalletAddrToName(walletAddr),
		ReserveAddr:    reserveAddr,
		Index:          logItem.Index,
	}
	tradeLog.Fees = append(tradeLog.Fees, tradelogFee)
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

func fillBurnFees(tradeLog common.TradelogV4, logItem types.Log) (common.TradelogV4, error) {
	reserveAddr, fee, err := logDataToBurnFeeParams(logItem.Data)
	if err != nil {
		return common.TradelogV4{}, err
	}

	burnFee := common.TradelogFee{
		ReserveAddr: reserveAddr,
		Burn:        fee.Big(),
		Index:       logItem.Index,
	}
	tradeLog.Fees = append(tradeLog.Fees, burnFee)
	if common.LengthBurnFees(tradeLog) == 1 {
		tradeLog.SrcReserveAddress = reserveAddr
	} else {
		tradeLog.DstReserveAddress = reserveAddr
	}
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

func isVolumeExcludedReserves(address ethereum.Address, volumeExcludedReserves []ethereum.Address) bool {
	for _, a := range volumeExcludedReserves {
		if address == a {
			return true
		}
	}
	return false
}

func fillKyberTradeV3(tradeLog common.TradelogV4, logItem types.Log, volumeExcludedReserves []ethereum.Address) (common.TradelogV4, error) {
	srcAddress, destAddress, srcReserve, dstReserve, receiverAddress, srcAmount, destAmount, ethAmount, err := logDataToKyberTradeV3Params(logItem.Data)
	if err != nil {
		return common.TradelogV4{}, err
	}
	tradeLog.TokenInfo = common.TradeTokenInfo{
		SrcAddress:  srcAddress,
		DestAddress: destAddress,
	}
	tradeLog.SrcAmount = srcAmount.Big()
	tradeLog.DestAmount = destAmount.Big()
	tradeLog.OriginalEthAmount = ethAmount.Big()
	tradeLog.SrcReserveAddress = srcReserve
	tradeLog.DstReserveAddress = dstReserve

	// update logic base on reserve instead of number of burn fee event
	defaultRatio := 2
	if isVolumeExcludedReserves(srcReserve, volumeExcludedReserves) {
		defaultRatio--
	}
	if isVolumeExcludedReserves(dstReserve, volumeExcludedReserves) {
		defaultRatio--
	}
	tradeLog.EthAmount = big.NewInt(1).Mul(ethAmount.Big(), big.NewInt(int64(defaultRatio)))

	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index
	tradeLog.User.UserAddress = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
	tradeLog.BlockNumber = logItem.BlockNumber
	tradeLog.ReceiverAddress = receiverAddress

	return tradeLog, nil
}

func (crawler *Crawler) calculateTradeAmount(t2ESrcAmounts, e2TSrcAmounts, t2ERates, e2TRates []*big.Int, srcToken, destToken ethereum.Address) (*big.Int, *big.Int) {
	srcAmount := big.NewInt(0)
	dstAmount := big.NewInt(0)
	if len(t2ESrcAmounts) != 0 {
		for _, amount := range t2ESrcAmounts {
			srcAmount = srcAmount.Add(srcAmount, amount)
		}
	} else {
		for _, amount := range e2TSrcAmounts {
			srcAmount = srcAmount.Add(srcAmount, amount)
		}
	}
	if len(e2TSrcAmounts) != 0 {
		for i, amount := range e2TSrcAmounts {
			dstAmount = dstAmount.Add(dstAmount, amount.Mul(amount, e2TRates[i]))
		}
	} else {
		for i, amount := range t2ESrcAmounts {
			dstAmount = dstAmount.Add(dstAmount, amount.Mul(amount, t2ERates[i]))
		}
	}
	return srcAmount, dstAmount
}

func (crawler *Crawler) fillKyberTradeV4(tradelog common.TradelogV4, logItem types.Log, volumeExcludedReserves []ethereum.Address) (common.TradelogV4, error) {
	var (
		logger = crawler.sugar.With("func", caller.GetCurrentFunctionName())
	)
	trade, err := crawler.kyberNetworkContract.ParseKyberTrade(logItem)
	if err != nil {
		logger.Errorw("failed to parse kyber trade event", "error", err)
		return tradelog, err
	}

	tradelog.TransactionHash = logItem.TxHash
	tradelog.BlockNumber = logItem.BlockNumber

	tradelog.TokenInfo = common.TradeTokenInfo{
		SrcAddress:  trade.Src,
		DestAddress: trade.Dest,
	}

	tradelog.T2EReserves = trade.T2eIds
	tradelog.E2TReserves = trade.E2tIds
	tradelog.T2ERates = trade.T2eRates
	tradelog.E2TRates = trade.E2tRates

	srcAmount, dstAmount := crawler.calculateTradeAmount(trade.T2eSrcAmounts, trade.E2tSrcAmounts, trade.T2eRates, trade.E2tRates, trade.Src, trade.Dest)
	tradelog.SrcAmount = srcAmount
	tradelog.DestAmount = dstAmount

	tradelog.EthAmount = trade.EthWeiValue
	tradelog.Index = logItem.Index

	return tradelog, nil
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

func fillKyberTradeV2(tradeLog common.TradelogV4, logItem types.Log) (common.TradelogV4, error) {
	srcAddr, destAddr, userAddr, receiverAddr, srcAmount, destAmount, err := logDataToKyberTradeV2Params(logItem.Data)
	if err != nil {
		return common.TradelogV4{}, err
	}
	tradeLog.TokenInfo.SrcAddress = srcAddr
	tradeLog.TokenInfo.DestAddress = destAddr

	tradeLog.SrcAmount = srcAmount.Big()
	tradeLog.DestAmount = destAmount.Big()

	tradeLog.TransactionHash = logItem.TxHash
	tradeLog.Index = logItem.Index

	tradeLog.User.UserAddress = userAddr
	tradeLog.ReceiverAddress = receiverAddr

	tradeLog.BlockNumber = logItem.BlockNumber

	return tradeLog, nil
}

// func assembleTradeLogsReserveAddr(log common.TradelogV4, sugar *zap.SugaredLogger) common.TradelogV4 {
// 	switch {
// 	case blockchain.IsBurnable(log.TokenInfo.SrcAddress):
// 		if blockchain.IsBurnable(log.TokenInfo.DestAddress) {
// 			if common.LengthBurnFees(log) == 2 {
// 				log.SrcReserveAddress = log.BurnFees[0].ReserveAddress
// 				log.DstReserveAddress = log.BurnFees[1].ReserveAddress
// 			} else {
// 				sugar.Warnw("unexpected burn fees", "got", common.LengthBurnFees(log), "want", "2 burn fees (src-dst)")
// 			}
// 		} else {
// 			if common.LengthBurnFees(log) == 1 {
// 				log.SrcReserveAddress = log.BurnFees[0].ReserveAddress
// 			} else {
// 				sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (src)")
// 			}
// 		}
// 	case blockchain.IsBurnable(log.DestAddress):
// 		if len(log.BurnFees) == 1 {
// 			log.DstReserveAddress = log.BurnFees[0].ReserveAddress
// 		} else {
// 			sugar.Warnw("unexpected burn fees", "got", log.BurnFees, "want", "1 burn fees (dst)")
// 		}
// 	case common.LengthWalletFees(log) != 0:
// 		if common.LengthWalletFees(log) == 1 {
// 			log.SrcReserveAddress = log.WalletFees[0].ReserveAddress
// 		} else {
// 			log.SrcReserveAddress = log.WalletFees[0].ReserveAddress
// 			log.DstReserveAddress = log.WalletFees[1].ReserveAddress
// 		}
// 	}

// 	return log
// }

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

func (crawler *Crawler) updateBasicInfo(log types.Log, tradeLog common.TradelogV4, timeout time.Duration) (common.TradelogV4, error) {
	var txSender ethereum.Address
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	tx, _, err := crawler.ethClient.TransactionByHash(ctx, log.TxHash)
	if err != nil {
		return tradeLog, err
	}
	txSender, err = crawler.ethClient.TransactionSender(ctx, tx, log.BlockHash, log.TxIndex)
	tradeLog.TxDetail.TxSender = txSender
	tradeLog.TxDetail.GasPrice = tx.GasPrice()

	if common.LengthBurnFees(tradeLog) == 0 { // in case there's no fee, we try to get wallet addr from tradeWithHint input
		if tx.To() != nil && bytes.Equal(tx.To().Bytes(), crawler.networkProxy.Bytes()) { // try to fail early, tx must have dst == networkProxy
			tradeParam, err := decodeTradeInputParam(tx.Data())
			if err != nil {
				return tradeLog, errors.Wrapf(err, "failed to decode input param, tx %s", tx.Hash().String())
			}
			tradeLog.WalletAddress = tradeParam.WalletID
			tradeLog.WalletName = WalletAddrToName(tradeLog.WalletAddress)
		} else {
			crawler.sugar.Warnw("no walletFee but tx is not with dest is networkProxy, skip get wallet addr")
		}
	} else {
		for _, fee := range tradeLog.Fees {
			if fee.PlatformFee.Cmp(big.NewInt(0)) != 0 {
				tradeLog.WalletAddress = fee.PlatformWallet
				tradeLog.WalletName = WalletAddrToName(tradeLog.WalletAddress)
				break
			}
		}
	}

	return tradeLog, err
}

func (crawler *Crawler) updateBasicInfoV4(log types.Log, tradeLog common.TradelogV4, timeout time.Duration) (common.TradelogV4, error) {
	var txSender ethereum.Address
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	tx, _, err := crawler.ethClient.TransactionByHash(ctx, log.TxHash)
	if err != nil {
		return tradeLog, err
	}
	txSender, err = crawler.ethClient.TransactionSender(ctx, tx, log.BlockHash, log.TxIndex)
	tradeLog.TxDetail.TxSender = txSender
	tradeLog.TxDetail.GasPrice = tx.GasPrice()
	tradeLog.User.UserAddress = txSender

	return tradeLog, err
}

// GetTradeLogs returns trade logs from KyberNetwork.
func (crawler *Crawler) GetTradeLogs(fromBlock, toBlock *big.Int, timeout time.Duration) (*common.CrawlResult, error) {
	var (
		result  *common.CrawlResult
		fetchFn tradeLogFetcher
	)

	// fetchTradeLogV2 also works for v3 trades, so to keep it simple, we only use fetchTradeLogV3 if both
	// from, to blocks are >= starting block v3
	switch {
	case fromBlock.Uint64() >= crawler.startingBlocks.V4() && toBlock.Uint64() >= crawler.startingBlocks.V4():
		fetchFn = crawler.fetchTradeLogV4
	case fromBlock.Uint64() >= crawler.startingBlocks.V3() && toBlock.Uint64() >= crawler.startingBlocks.V3():
		fetchFn = crawler.fetchTradeLogV3
	case fromBlock.Uint64() >= crawler.startingBlocks.V2() && toBlock.Uint64() >= crawler.startingBlocks.V2():
		fetchFn = crawler.fetchTradeLogV2
	default:
		fetchFn = crawler.fetchTradeLogV1
	}

	result, err := fetchFn(fromBlock, toBlock, timeout)
	if err != nil {
		return result, errors.Wrapf(err, "failed to fetch trade logs fromBlock: %v toBlock:%v", fromBlock, toBlock)
	}
	if result == nil {
		return result, nil
	}
	for _, tradeLog := range result.Trades {
		var uid, ip, country string

		uid, ip, country, err = crawler.broadcastClient.GetTxInfo(tradeLog.TransactionHash.Hex())
		if err != nil {
			return result, err
		}
		tradeLog.User.IP = ip
		tradeLog.User.Country = country
		tradeLog.User.UID = uid

		if tradeLog.IsKyberSwap() {
			tradeLog.IntegrationApp = appname.KyberSwapAppName
		} else {
			tradeLog.IntegrationApp = appname.ThirdPartyAppName
		}

		rate, err := crawler.rateProvider.USDRate(tradeLog.Timestamp)
		if err != nil {
			return nil, err
		}
		tradeLog.ETHUSDProvider = crawler.rateProvider.Name()
		tradeLog.ETHUSDRate = rate
	}
	return result, nil
}
