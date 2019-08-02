package common

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	//TradeExecute(address sender, address src, uint256 srcAmount, address destToken, uint256 destAmount, address destAddress)
	tradeExecuteEvent = "0xea9415385bae08fe9f6dc457b02577166790cde83bb18cc340aac6cb81b824de"
	//Transfer event
	transferEvent = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

	timeout = 10 * time.Second
	attempt = 4
)

// NormalTx holds info from normal tx query.
type NormalTx struct {
	BlockNumber int       `json:"blockNumber,string"`
	Timestamp   time.Time `json:"-"`
	Hash        string    `json:"hash"`
	BlockHash   string    `json:"blockHash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       *big.Int  `json:"value"`
	Gas         int       `json:"gas,string"`
	GasUsed     int       `json:"gasUsed,string"`
	GasPrice    *big.Int  `json:"gasPrice"`
	IsError     int       `json:"isError,string"`
}

// UnmarshalJSON is the custom unmarshaller that read timestamp in unix milliseconds.
func (tx *NormalTx) UnmarshalJSON(data []byte) error {
	type AliasNormalTx NormalTx
	decoded := new(struct {
		AliasNormalTx
		Timestamp uint64 `json:"timestamp"`
	})

	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	tx.BlockNumber = decoded.BlockNumber
	tx.Timestamp = timeutil.TimestampMsToTime(decoded.Timestamp).UTC()
	tx.Hash = decoded.Hash
	tx.BlockHash = decoded.BlockHash
	tx.From = decoded.From
	tx.To = decoded.To
	tx.Value = decoded.Value
	tx.Gas = decoded.Gas
	tx.GasUsed = decoded.GasUsed
	tx.GasPrice = decoded.GasPrice
	tx.IsError = decoded.IsError
	return nil
}

// MarshalJSON is the custom JSON marshaller that output timestamp in unix milliseconds.
func (tx NormalTx) MarshalJSON() ([]byte, error) {
	type AliasTNormalTx NormalTx
	return json.Marshal(struct {
		AliasTNormalTx
		Timestamp uint64 `json:"timestamp"`
	}{
		AliasTNormalTx: (AliasTNormalTx)(tx),
		Timestamp:      timeutil.TimeToTimestampMs(tx.Timestamp),
	})
}

// InternalTx holds info from internal tx query.
type InternalTx struct {
	BlockNumber int       `json:"blockNumber,string"`
	Timestamp   time.Time `json:"timestamp"`
	Hash        string    `json:"hash"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       *big.Int  `json:"value"`
	Gas         int       `json:"gas,string"`
	GasUsed     int       `json:"gasUsed,string"`
	IsError     int       `json:"isError,string"`
	IsTrade     bool      `json:"isTrade"`
}

// UnmarshalJSON is the custom unmarshaller that read timestamp in unix milliseconds.
func (tx *InternalTx) UnmarshalJSON(data []byte) error {
	type AliasInternalTx InternalTx
	decoded := new(struct {
		AliasInternalTx
		Timestamp uint64 `json:"timestamp"`
	})
	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	tx.BlockNumber = decoded.BlockNumber
	tx.Timestamp = timeutil.TimestampMsToTime(decoded.Timestamp).UTC()
	tx.Hash = decoded.Hash
	tx.From = decoded.From
	tx.To = decoded.To
	tx.Value = decoded.Value
	tx.Gas = decoded.Gas
	tx.GasUsed = decoded.GasUsed
	tx.IsError = decoded.IsError
	tx.IsTrade = decoded.IsTrade
	return nil
}

// MarshalJSON is the custom JSON marshaller that output timestamp in unix milliseconds.
func (tx InternalTx) MarshalJSON() ([]byte, error) {
	type AliasInternalTx InternalTx
	return json.Marshal(struct {
		AliasInternalTx
		Timestamp uint64 `json:"timestamp"`
	}{
		AliasInternalTx: (AliasInternalTx)(tx),
		Timestamp:       timeutil.TimeToTimestampMs(tx.Timestamp),
	})
}

// ERC20Transfer holds info from ERC20 token transfer event query.
type ERC20Transfer struct {
	BlockNumber     int              `json:"blockNumber,string"`
	Timestamp       time.Time        `json:"-"`
	Hash            ethereum.Hash    `json:"-"`
	From            ethereum.Address `json:"-"`
	ContractAddress ethereum.Address `json:"-"`
	To              ethereum.Address `json:"-"`
	Value           *big.Int         `json:"value"`
	Gas             int              `json:"gas,string"`
	GasUsed         int              `json:"gasUsed,string"`
	GasPrice        *big.Int         `json:"gasPrice"`
	IsTrade         bool             `json:"isTrade"`
}

//MarshalJSON return marshal form of erc20transfer
func (et ERC20Transfer) MarshalJSON() ([]byte, error) {
	type AliasErc20 ERC20Transfer
	var ts *uint64
	if !et.Timestamp.IsZero() {
		millis := timeutil.TimeToTimestampMs(et.Timestamp)
		ts = &millis
	}

	return json.Marshal(struct {
		Timestamp       *uint64 `json:"timestamp"`
		Hash            string  `json:"hash"`
		From            string  `json:"from"`
		ContractAddress string  `json:"contractAddress"`
		To              string  `json:"to"`
		AliasErc20
	}{
		Timestamp:       ts,
		Hash:            et.Hash.Hex(),
		From:            et.From.Hex(),
		ContractAddress: et.ContractAddress.Hex(),
		To:              et.To.Hex(),
		AliasErc20:      (AliasErc20)(et),
	})
}

//UnmarshalJSON return an ERC20Transfer object from JSON form
func (et *ERC20Transfer) UnmarshalJSON(data []byte) error {
	type AliasErc20 ERC20Transfer
	decoded := new(struct {
		Timestamp       *uint64 `json:"timestamp,omitempty"`
		Hash            string  `json:"hash"`
		From            string  `json:"from"`
		ContractAddress string  `json:"contractAddress"`
		To              string  `json:"to"`
		AliasErc20
	})
	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	if decoded.Timestamp != nil {
		et.Timestamp = timeutil.TimestampMsToTime(*decoded.Timestamp).UTC()
	}

	et.BlockNumber = decoded.BlockNumber

	et.Hash = ethereum.HexToHash(decoded.Hash)
	et.From = ethereum.HexToAddress(decoded.From)
	et.To = ethereum.HexToAddress(decoded.To)
	et.ContractAddress = ethereum.HexToAddress(decoded.ContractAddress)
	et.Value = decoded.Value
	et.Gas = decoded.Gas
	et.GasUsed = decoded.GasUsed
	et.GasPrice = decoded.GasPrice
	et.IsTrade = decoded.IsTrade

	return nil
}

//DetectTradeInternalTransaction detect if a provided txHash is belong to a trade transaction or not
func DetectTradeInternalTransaction(txHash ethereum.Hash, ethAmount *big.Int, ethClient *ethclient.Client) (bool, error) {
	var (
		err     error
		receipt *types.Receipt
	)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for i := 0; i < attempt; i++ {
		receipt, err = ethClient.TransactionReceipt(ctx, txHash)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return false, err
	}
	for _, log := range receipt.Logs {
		for _, topic := range log.Topics {
			if topic == ethereum.HexToHash(tradeExecuteEvent) {
				//check if the amount is match then it is belong to this trade
				destAddress := ethereum.BytesToAddress(log.Data[64:96])

				//only compare if the destAddress is ethereum
				if destAddress == blockchain.ETHAddr {
					destAmount := ethereum.BytesToHash(log.Data[96:128])
					if destAmount.Big().Cmp(ethAmount) == 0 {
						return true, nil
					}
				}

				//check if the amount is match then it is belong to this trade
				srcAddress := ethereum.BytesToAddress(log.Data[0:32])

				//only compare if the destAddress is ethereum
				if srcAddress == blockchain.ETHAddr {
					srcAmount := ethereum.BytesToHash(log.Data[32:64])
					if srcAmount.Big().Cmp(ethAmount) == 0 {
						return true, nil
					}
				}
			}
		}
	}
	return false, nil
}

//EtherscanInternalTxToCommon transforms etherScan.InternalTx to accounting's InternalTx
func EtherscanInternalTxToCommon(tx etherscan.InternalTx, ethClient *ethclient.Client, sugar *zap.SugaredLogger, throttle chan int) (InternalTx, error) {
	var (
		logger = sugar.With("func", "tx-types/EtherscanInternalTxToCommon")
	)
	defer func() {
		<-throttle
	}()
	// Detect if an internal tx is belong to a trade or not
	isTrade, err := DetectTradeInternalTransaction(ethereum.HexToHash(tx.Hash), tx.Value.Int(), ethClient)
	if err != nil {
		return InternalTx{}, err
	}
	logger.Infow("tx is trade detection", "tx", tx.Hash, "is trade", isTrade)
	return InternalTx{
		BlockNumber: tx.BlockNumber,
		Timestamp:   tx.TimeStamp.Time(),
		Hash:        tx.Hash,
		From:        tx.From,
		To:          tx.To,
		Value:       tx.Value.Int(),
		Gas:         tx.Gas,
		GasUsed:     tx.GasUsed,
		IsError:     tx.IsError,
		IsTrade:     isTrade,
	}, nil
}

//DetectTradeTransaction detect if a provided txHash is belong to a trade transaction or not
func DetectTradeTransaction(tx etherscan.ERC20Transfer, ethClient *ethclient.Client) (bool, error) {
	var (
		transferEventIndex uint
		receipt            *types.Receipt
		err                error
	)
	txHash := ethereum.HexToHash(tx.Hash)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for i := 0; i < attempt; i++ {
		receipt, err = ethClient.TransactionReceipt(ctx, txHash)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return false, err
	}
	for _, log := range receipt.Logs {
		for _, topic := range log.Topics {
			if topic == ethereum.HexToHash(transferEvent) {
				fromAddress := ethereum.HexToAddress(log.Topics[1].Hex())
				toAddress := ethereum.HexToAddress(log.Topics[2].Hex())
				value := ethereum.BytesToHash(log.Data)

				// detect event in logs
				if fromAddress == ethereum.HexToAddress(tx.From) &&
					toAddress == ethereum.HexToAddress(tx.To) && value.Big() == tx.Value.Int() {
					transferEventIndex = log.Index
				}
				// if there is another transfer into reserve
				// after the current transfer and before trade, then the current transfer is not belong to trade
				if toAddress == ethereum.HexToAddress(tx.To) && transferEventIndex != 0 && transferEventIndex < log.Index {
					return false, nil
				}
			}
			if topic == ethereum.HexToHash(tradeExecuteEvent) {
				// if transferEvent not yet reach then, it is definitely a trade
				if transferEventIndex == 0 {
					return true, nil
				}

				// if trade from ETH - Token,
				// the transfer event appear before the trade is not a trade transfer event
				srcToken := ethereum.BytesToAddress(log.Data[0:32])
				if srcToken == blockchain.ETHAddr {
					return false, nil
				}
				// if trade from Token - ETH,
				// then transfer token must be the same with srcToken for first trade event
				if ethereum.HexToAddress(tx.ContractAddress) != srcToken {
					return false, nil
				}
				return true, nil
			}
		}
	}
	return false, nil
}

//EtherscanERC20TransferToCommon transforms etherScan.ERC20Trasnfer to accounting's ERC20Transfer
func EtherscanERC20TransferToCommon(tx etherscan.ERC20Transfer, ethClient *ethclient.Client, sugar *zap.SugaredLogger, throttle chan int) (ERC20Transfer, error) {
	var (
		logger = sugar.With("func", "tx-types/EtherscanERC20TransferToCommon")
	)
	defer func() {
		<-throttle
	}()
	//Detect if a transfer transaction is a trade or not
	isTrade, err := DetectTradeTransaction(tx, ethClient)
	if err != nil {
		return ERC20Transfer{}, err
	}
	logger.Infow("tx is trade detected", "tx", tx.Hash, "is trade", isTrade)
	return ERC20Transfer{
		BlockNumber:     tx.BlockNumber,
		Timestamp:       tx.TimeStamp.Time(),
		Hash:            ethereum.HexToHash(tx.Hash),
		From:            ethereum.HexToAddress(tx.From),
		ContractAddress: ethereum.HexToAddress(tx.ContractAddress),
		To:              ethereum.HexToAddress(tx.To),
		Value:           tx.Value.Int(),
		Gas:             tx.Gas,
		GasUsed:         tx.GasUsed,
		GasPrice:        tx.GasPrice.Int(),
		IsTrade:         isTrade,
	}, nil
}

//EtherscanNormalTxToCommon transform etherScan.NormalTx to accounting's normalTx
func EtherscanNormalTxToCommon(tx etherscan.NormalTx) NormalTx {
	return NormalTx{
		BlockNumber: tx.BlockNumber,
		Timestamp:   tx.TimeStamp.Time(),
		Hash:        tx.Hash,
		BlockHash:   tx.BlockHash,
		From:        tx.From,
		To:          tx.To,
		Value:       tx.Value.Int(),
		Gas:         tx.Gas,
		GasUsed:     tx.GasUsed,
		GasPrice:    tx.GasPrice.Int(),
		IsError:     tx.IsError,
	}
}
