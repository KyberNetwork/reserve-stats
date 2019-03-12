package fetcher

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/nanmu42/etherscan-api"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-transaction-fetcher/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
)

// EtherscanTransactionFetcher is an implementation of TransactionFetcher that uses Etherscan API.
type EtherscanTransactionFetcher struct {
	sugar  *zap.SugaredLogger
	client *etherscan.Client
}

// NewEtherscanTransactionFetcher returns a new EtherscanTransactionFetcher instance.
func NewEtherscanTransactionFetcher(sugar *zap.SugaredLogger, client *etherscan.Client) *EtherscanTransactionFetcher {
	return &EtherscanTransactionFetcher{sugar: sugar, client: client}
}

type fetchFn struct {
	name  string
	fetch func(address string, startBlock *int, endBlock *int, page int, offset int) ([]interface{}, error)
}

func newFetchFunction(name string, fetch func(address string, startBlock *int, endBlock *int, page int, offset int) ([]interface{}, error)) *fetchFn {
	return &fetchFn{name: name, fetch: fetch}
}

func (f *EtherscanTransactionFetcher) fetch(fn *fetchFn, addr ethereum.Address, from, to *big.Int) ([]interface{}, error) {
	// maximum number of transactions to return in a page.
	// Too small value will increase the fetching time, too big value will result in a timed out response.
	const offset = 500
	var (
		logger = f.sugar.With(
			"func",
			"accounting/reserve-transaction-fetcher/fetcher.EtherscanTransactionFetcher.fetch",
			"fetch_function", fn.name,
			"address", addr.String(),
			"offset", offset,
		)
		results []interface{}
	)

	var (
		startBlock *int
		endBlock   *int
	)

	if from != nil {
		logger = logger.With("start_block", from.String())
		fromVal := int(from.Int64())
		startBlock = &fromVal
	}

	if to != nil {
		// Ethereum API includes the transactions of to block
		to.Sub(to, big.NewInt(1))

		logger = logger.With("endBlock", to.String())
		toVal := int(to.Int64())
		endBlock = &toVal
	}

	logger.Info("fetching transactions")

	// Etherscan paging starts with index=1
	for page := 1; ; page++ {
		logger.Debugw("fetching a page of transactions", "page", page)
		txs, err := fn.fetch(addr.String(), startBlock, endBlock, page, offset)
		if blockchain.IsEtherscanNotransactionFound(err) {
			logger.Debugw("all transaction fetched", "page", page)
			break
		} else if err != nil {
			return nil, err
		}
		results = append(results, txs...)
	}
	return results, nil
}

func etherscanNormalTxToCommon(tx etherscan.NormalTx) common.NormalTx {
	return common.NormalTx{
		BlockNumber: tx.BlockNumber,
		TimeStamp:   tx.TimeStamp.Time(),
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

func (f *EtherscanTransactionFetcher) NormalTx(addr ethereum.Address, from, to *big.Int) ([]common.NormalTx, error) {
	fn := newFetchFunction("normal", func(address string, startBlock *int, endBlock *int, page int, offset int) ([]interface{}, error) {
		normalTxs, err := f.client.NormalTxByAddress(address, startBlock, endBlock, page, offset, false)
		if err != nil {
			return nil, err
		}
		results := make([]interface{}, len(normalTxs))
		for i, v := range normalTxs {
			results[i] = v
		}
		return results, nil
	})

	results, err := f.fetch(fn, addr, from, to)
	if err != nil {
		return nil, err
	}

	txs := make([]common.NormalTx, len(results))
	for i, v := range results {
		tx := v.(etherscan.NormalTx)
		txs[i] = etherscanNormalTxToCommon(tx)
	}
	return txs, nil
}

func etherscanInternalTxToCommon(tx etherscan.InternalTx) common.InternalTx {
	return common.InternalTx{
		BlockNumber: tx.BlockNumber,
		TimeStamp:   tx.TimeStamp.Time(),
		Hash:        tx.Hash,
		From:        tx.From,
		To:          tx.To,
		Value:       tx.Value.Int(),
		Gas:         tx.Gas,
		GasUsed:     tx.GasUsed,
		IsError:     tx.IsError,
	}
}

func (f *EtherscanTransactionFetcher) InternalTx(addr ethereum.Address, from, to *big.Int) ([]common.InternalTx, error) {
	fn := newFetchFunction("internal", func(address string, startBlock *int, endBlock *int, page int, offset int) ([]interface{}, error) {
		internalTxs, err := f.client.InternalTxByAddress(address, startBlock, endBlock, page, offset, false)
		if err != nil {
			return nil, err
		}
		results := make([]interface{}, len(internalTxs))
		for i, v := range internalTxs {
			results[i] = v
		}
		return results, nil
	})

	results, err := f.fetch(fn, addr, from, to)
	if err != nil {
		return nil, err
	}

	txs := make([]common.InternalTx, len(results))
	for i, v := range results {
		tx := v.(etherscan.InternalTx)
		txs[i] = etherscanInternalTxToCommon(tx)
	}
	return txs, nil
}

func etherscanERC20TransferToCommon(tx etherscan.ERC20Transfer) common.ERC20Transfer {
	return common.ERC20Transfer{
		BlockNumber:     tx.BlockNumber,
		TimeStamp:       tx.TimeStamp.Time(),
		Hash:            tx.Hash,
		From:            tx.From,
		ContractAddress: tx.ContractAddress,
		To:              tx.To,
		Value:           tx.Value.Int(),
		Gas:             tx.Gas,
		GasUsed:         tx.GasUsed,
		GasPrice:        tx.GasPrice.Int(),
	}
}

func (f *EtherscanTransactionFetcher) ERC20Transfer(addr ethereum.Address, from, to *big.Int) ([]common.ERC20Transfer, error) {
	fn := newFetchFunction("transfer", func(address string, startBlock *int, endBlock *int, page int, offset int) ([]interface{}, error) {
		transfers, err := f.client.ERC20Transfers(nil, &address, startBlock, endBlock, page, offset)
		if err != nil {
			return nil, err
		}
		results := make([]interface{}, len(transfers))
		for i, v := range transfers {
			results[i] = v
		}
		return results, nil
	})

	results, err := f.fetch(fn, addr, from, to)
	if err != nil {
		return nil, err
	}

	transfers := make([]common.ERC20Transfer, len(results))
	for i, v := range results {
		transfer := v.(etherscan.ERC20Transfer)
		transfers[i] = etherscanERC20TransferToCommon(transfer)
	}
	return transfers, nil
}
