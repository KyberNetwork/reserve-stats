package tradelogs

import (
	"context"
	"errors"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	ether "github.com/ethereum/go-ethereum"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/KyberNetwork/reserve-stats/common"
)

const (
	// tradeEvent is the topic of event
	// ExecuteTrade(address indexed sender, ERC20 src, ERC20 dest, uint actualSrcAmount, uint actualDestAmount).
	tradeEvent = "0x1849bd6a030a1bca28b83437fd3de96f3d27a5d172fa7e9c78e7b61468928a39"
	// etherReceivalEvent is the topic of event EtherReceival(address indexed sender, uint amount).
	etherReceivalEvent = "0x75f33ed68675112c77094e7c5b073890598be1d23e27cd7f6907b4a7d98ac619"

	// address of pricing contract
	pricingAddr = "0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B"
	// address of network contract
	networkAddr = "0x818E6FECD516Ecc3849DAf6845e3EC868087B755"

	ethDecimals int64  = 18
	ethAddress  string = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

// TradeLogCrawler gets trade logs on KyberNetwork on blockchain, adding the
// information about USD equivalent on each trade.
type TradeLogCrawler struct {
	ethClient *ethclient.Client
	ethRate   EthUSDRate
}

func logDataToTradeParams(data []byte) (ethereum.Address, ethereum.Address, ethereum.Hash, ethereum.Hash) {
	srcAddr := ethereum.BytesToAddress(data[0:32])
	desAddr := ethereum.BytesToAddress(data[32:64])
	srcAmount := ethereum.BytesToHash(data[64:96])
	desAmount := ethereum.BytesToHash(data[96:128])
	return srcAddr, desAddr, srcAmount, desAmount
}

func logDataToEtherReceivalParams(data []byte) ethereum.Hash {
	amount := ethereum.BytesToHash(data[0:32])
	return amount
}

func updateTradeLogs(allLogs []common.KNLog, logItem types.Log, ts uint64) ([]common.KNLog, error) {
	var (
		tradeLog      common.TradeLog
		updateLastLog = false
	)

	if len(logItem.Topics) == 0 {
		return allLogs, errors.New("log item has no topic")
	}

	// if the transaction hash is the same and last log is a TradeLog, update it,
	// otherwise append new one
	if len(allLogs) > 0 && allLogs[len(allLogs)-1].TxHash() == logItem.TxHash {
		var ok bool
		if tradeLog, ok = allLogs[len(allLogs)-1].(common.TradeLog); ok {
			updateLastLog = true
		}
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
	case etherReceivalEvent:
		amount := logDataToEtherReceivalParams(logItem.Data)
		tradeLog.EtherReceivalSender = ethereum.BytesToAddress(logItem.Topics[1].Bytes())
		tradeLog.EtherReceivalAmount = amount.Big()
	case tradeEvent:
		srcAddr, destAddr, srcAmount, destAmount := logDataToTradeParams(logItem.Data)
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
func NewTradeLogCrawler(nodeURL string, ethRate EthUSDRate) (*TradeLogCrawler, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, err
	}
	return &TradeLogCrawler{client, ethRate}, nil
}

// GetTradeLogs returns trade logs from KyberNetwork.
func (crawler *TradeLogCrawler) GetTradeLogs(fromBlock uint64, toBlock uint64) ([]common.KNLog, error) {
	var result []common.KNLog

	addresses := []ethereum.Address{
		ethereum.HexToAddress(pricingAddr), // pricing
		ethereum.HexToAddress(networkAddr), // network
	}

	topics := [][]ethereum.Hash{
		{
			ethereum.HexToHash(tradeEvent),
			ethereum.HexToHash(etherReceivalEvent),
		},
	}

	query := ether.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: addresses,
		Topics:    topics,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logs, err := crawler.ethClient.FilterLogs(ctx, query)
	if err != nil {
		return result, err
	}

	for _, logItem := range logs {
		if logItem.Removed {
			continue // Removed due to chain reorg
		}

		ts, err := crawler.InterpretTimestamp(
			logItem.BlockNumber,
			logItem.Index,
		)

		if err != nil {
			return result, err
		}

		topic := logItem.Topics[0]
		switch topic.Hex() {
		case etherReceivalEvent, tradeEvent:
			// add logItem to result
			log.Printf("Trade Event at %d.", ts)
			if result, err = updateTradeLogs(result, logItem, ts); err != nil {
				return result, err
			}
		default:
			log.Println("Unknow topic.")
		}
	}

	for i, logItem := range result {
		tradeLog, ok := logItem.(common.TradeLog)
		if !ok {
			continue
		}

		ethRate := crawler.ethRate.GetUSDRate(tradeLog.Timestamp / 1000000)
		if ethRate != 0 {
			result[i] = calculateFiatAmount(tradeLog, ethRate)
		}
	}

	return result, nil
}

// CachedBlockno is block number of a block which was cached it's header.
var CachedBlockno uint64

// CachedBlockHeader is the cached block header
var CachedBlockHeader *types.Header

// InterpretTimestamp returns timestamp from block number and transaction index.
// It cached block number and block header to reduces the number of request
// to node.
func (crawler *TradeLogCrawler) InterpretTimestamp(blockno uint64, txindex uint) (uint64, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var block *types.Header
	var err error
	if CachedBlockno == blockno {
		block = CachedBlockHeader
	} else {
		block, err = crawler.ethClient.HeaderByNumber(timeout, big.NewInt(int64(blockno)))
	}
	if err != nil {
		if block == nil {
			return uint64(0), err
		}

		// error because parity and geth are not compatible in mix hash
		// so we ignore it
		CachedBlockno = blockno
		CachedBlockHeader = block
		err = nil
	}
	unixSecond := block.Time.Uint64()
	unixNano := uint64(time.Unix(int64(unixSecond), 0).UnixNano())
	result := unixNano + uint64(txindex)
	return result, nil
}
