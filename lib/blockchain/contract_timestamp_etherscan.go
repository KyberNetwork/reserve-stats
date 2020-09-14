package blockchain

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/nanmu42/etherscan-api"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
)

// IsEtherscanNotransactionFound returns true if given error represents etherscan no transaction found error.
func IsEtherscanNotransactionFound(err error) bool {
	return err != nil && err.Error() == "etherscan server: No transactions found"
}

func IsEtherscanRateLimit(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "response status 403 403 Forbidden")
}

// EtherscanContractTimestampResolver is an implementation of EtherscanContractTimestampResolver
// that uses Etherscan API.
type EtherscanContractTimestampResolver struct {
	sugar  *zap.SugaredLogger
	client *etherscan.Client
}

// NewEtherscanContractTimestampResolver creates a new EtherscanContractTimestampResolver from given client.
func NewEtherscanContractTimestampResolver(sugar *zap.SugaredLogger, client *etherscan.Client) *EtherscanContractTimestampResolver {
	return &EtherscanContractTimestampResolver{sugar: sugar, client: client}
}

func (r *EtherscanContractTimestampResolver) resolveUsingInternalTx(address common.Address) (time.Time, error) {
	var logger = r.sugar.With(
		"func", caller.GetCurrentFunctionName(),
		"address", address.String(),
	)

	logger.Debug("fetching internal transactions from Etherscan")
	txs, err := r.client.InternalTxByAddress(address.String(), nil, nil, 1, 1, false)
	if err != nil {
		// etherscan package does not export error for this, have to compare error message
		if IsEtherscanNotransactionFound(err) {
			return time.Time{}, ErrNotAvailable
		}
	}
	logger.Debugw("got transactions from Etherscan", "txs", len(txs))

	// with current implementation of etherscan-api, the client will return an error with
	// message errMsg if no transaction found for given address.
	// Following codes should never be reached, just add for safe guard for implementation changes.
	if len(txs) == 0 {
		return time.Time{}, ErrNotAvailable
	}

	firstTx := txs[0]
	logger.Debugw("got first transaction", "first_tx", firstTx.Hash)

	// first transaction is not a contract creation transaction, given address is not a contract.
	if len(firstTx.ContractAddress) == 0 {
		logger.Errorw("Contract does not exist", "contract", address.Hex())
		return time.Time{}, ErrNotAvailable
	}

	return firstTx.TimeStamp.Time(), nil
}

// Resolve returns the creation timestamp of given contract address using Etherscan API.
func (r *EtherscanContractTimestampResolver) Resolve(address common.Address) (time.Time, error) {
	var logger = r.sugar.With(
		"func", caller.GetCurrentFunctionName(),
		"address", address.String(),
	)

	logger.Debug("fetching normal transactions from Etherscan")
	txs, err := r.client.NormalTxByAddress(address.String(), nil, nil, 1, 1, false)
	if err != nil {
		// etherscan package does not export error for this, have to compare error message
		if IsEtherscanNotransactionFound(err) {
			// fallback check internal tx
			logger.Debugw("fallback calling internal tx due to error", "error", err.Error())
			contractTimeStamp, err := r.resolveUsingInternalTx(address)
			if err != nil {
				return time.Time{}, err
			}
			return contractTimeStamp, nil
		}
		return time.Time{}, err
	}
	logger.Debugw("got transactions from Etherscan", "txs", len(txs))

	// with current implementation of etherscan-api, the client will return an error with
	// message errMsg if no transaction found for given address.
	// Following codes should never be reached, just add for safe guard for implementation changes.
	if len(txs) == 0 {
		return time.Time{}, ErrNotAvailable
	}

	firstTx := txs[0]
	logger.Debugw("got first transaction", "first_tx", firstTx.Hash)

	// first transaction is not a contract creation transaction, given address is not a contract.
	if len(firstTx.ContractAddress) == 0 {
		// fallback check internal tx
		logger.Debug("fallback calling internal tx due to tx contract address is empty")
		contractTimestamp, err := r.resolveUsingInternalTx(address)
		if err != nil {
			logger.Info("Check erc20 txs, get timestamp as the first tx available")
			contractTimestamp = firstTx.TimeStamp.Time()
			accountAddress := address.String()
			txs, err := r.client.ERC20Transfers(nil, &accountAddress, nil, nil, 1, 200) //
			if err != nil {
				// if could not get txs, return  timestamp as first normal txs
				logger.Debug("failed to get erc20 transfers txs", "error", err)
				return contractTimestamp, nil
			}
			if len(txs) == 0 {
				logger.Debug("failed to get erc20 tranfers txs, length txs equal to 0")
				// if could not get txs, return  timestamp as first normal txs
				return contractTimestamp, nil
			}
			firstTransferTx := txs[0]
			logger.Debug("got first erc20 transfer tx", "tx", firstTransferTx)
			if firstTransferTx.TimeStamp.Time().Before(contractTimestamp) {
				// if firstTransfer is first tx then account timestamp is first tranfer timestamp
				contractTimestamp = firstTransferTx.TimeStamp.Time()
			}
		}
		return contractTimestamp, nil
	}

	return firstTx.TimeStamp.Time(), nil
}
