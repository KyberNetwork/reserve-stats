package fetcher

import (
	"fmt"
	"math/big"
	"sync"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
)

const (
	maxConcurrentRequest = 50
)

// EtherscanTransactionFetcher is an implementation of TransactionFetcher that uses Etherscan API.
type EtherscanTransactionFetcher struct {
	sugar     *zap.SugaredLogger
	client    *etherscan.Client
	ethClient *ethclient.Client
	attempt   int
}

// NewEtherscanTransactionFetcher returns a new EtherscanTransactionFetcher instance.
func NewEtherscanTransactionFetcher(sugar *zap.SugaredLogger, client *etherscan.Client, ethClient *ethclient.Client, attempt int) *EtherscanTransactionFetcher {
	return &EtherscanTransactionFetcher{sugar: sugar, client: client, ethClient: ethClient, attempt: attempt}
}

type fetchFn struct {
	name  string
	fetch func(address string, startBlock *int, endBlock *int, page int, offset int) ([]interface{}, error)
}

func newFetchFunction(name string, fetch func(address string, startBlock *int, endBlock *int, page int, offset int) ([]interface{}, error)) *fetchFn {
	return &fetchFn{name: name, fetch: fetch}
}

func (f *EtherscanTransactionFetcher) fetchWithRetry(fn *fetchFn, addr ethereum.Address, startBlock, endBlock *int, page, offset int) ([]interface{}, error) {
	var (
		txs    []interface{}
		err    error
		logger = f.sugar.With("func", caller.GetCurrentFunctionName())
	)
	for i := 0; i < f.attempt; i++ {
		txs, err = fn.fetch(addr.String(), startBlock, endBlock, page, offset)
		if blockchain.IsEtherscanNotransactionFound(err) {
			// the fetcher reach the end of result
			break
		} else if err != nil {
			logger.Warnw("Fetch data from Etherscan api failed, retry", "error", err, "offset", offset, "attempt", i)
			continue
		}
		logger.Infow("fetch data success", "page", page, "attempt", i, "offset", offset)
		break
	}
	return txs, err
}

func (f *EtherscanTransactionFetcher) fetch(fn *fetchFn, addr ethereum.Address, from, to *big.Int, offset int) ([]interface{}, error) {
	// maximum number of transactions to return in a page.
	// Too small value will increase the fetching time, too big value will result in a timed out response.
	var (
		logger = f.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"fetch_function", fn.name,
			"address", addr.String(),
			"offset", offset,
		)
		results []interface{}
	)
	// clone from, to value to avoid changing
	if from != nil {
		from = big.NewInt(0).Set(from)
	}
	if to != nil {
		to = big.NewInt(0).Set(to)
	}

	var (
		startBlock *int
		endBlock   *int
	)

	if from != nil {
		logger = logger.With("start_block", from.String())
		if !from.IsInt64() {
			return nil, fmt.Errorf("unsupported block: number=%s", from.String())
		}
		fromVal := int(from.Int64())
		startBlock = &fromVal
	}

	if to != nil {
		// Ethereum API includes the transactions of to block
		to.Sub(to, big.NewInt(1))
		if !to.IsInt64() {
			return nil, fmt.Errorf("unsupported block: number=%s", to.String())
		}
		logger = logger.With("endBlock", to.String())
		toVal := int(to.Int64())
		endBlock = &toVal
	}

	logger.Info("fetching transactions")

	// Etherscan paging starts with index=1
	for page := 1; ; page++ {
		logger.Debugw("fetching a page of transactions", "page", page)
		txs, err := f.fetchWithRetry(fn, addr, startBlock, endBlock, page, offset)
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

// NormalTx returns all normal Ethereum transaction of given address between block range.
func (f *EtherscanTransactionFetcher) NormalTx(addr ethereum.Address, from, to *big.Int, offset int) ([]common.NormalTx, error) {
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

	results, err := f.fetch(fn, addr, from, to, offset)
	if err != nil {
		return nil, err
	}

	txs := make([]common.NormalTx, len(results))
	for i, v := range results {
		tx := v.(etherscan.NormalTx)
		txs[i] = common.EtherscanNormalTxToCommon(tx)
	}
	return txs, nil
}

// InternalTx returns all internal transaction of given address between block range.
func (f *EtherscanTransactionFetcher) InternalTx(addr ethereum.Address, from, to *big.Int, offset int) ([]common.InternalTx, error) {
	var (
		logger   = f.sugar.With("func", caller.GetCurrentFunctionName())
		temp     sync.Map
		g        errgroup.Group
		throttle = make(chan int, maxConcurrentRequest)
	)
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

	results, err := f.fetch(fn, addr, from, to, offset)
	if err != nil {
		return nil, err
	}

	txs := make([]common.InternalTx, len(results))
	for i, v := range results {
		throttle <- 1
		tx := v.(etherscan.InternalTx)
		g.Go(func(index int) func() error {
			return func() error {
				internalTx, err := common.EtherscanInternalTxToCommon(tx, f.ethClient, f.sugar, throttle)
				if err != nil {
					logger.Errorw("convert from Etherscan internal tx to common failed", "error", err, "tx", tx.Hash)
					return err
				}
				temp.Store(index, internalTx)
				return nil
			}
		}(i),
		)
	}
	if err = g.Wait(); err != nil {
		return txs, err
	}

	temp.Range(func(index, value interface{}) bool {
		i, ok := index.(int)
		if !ok {
			err = fmt.Errorf("index (%vl) cannot be asserted to integer", index)
			return false
		}
		tx, ok := value.(common.InternalTx)
		if !ok {
			err = fmt.Errorf("value (%v) cannot be asserted to common.InternalTx", value)
			return false
		}
		txs[i] = tx
		return true
	})
	if err != nil {
		return txs, err
	}
	return txs, nil
}

// ERC20Transfer returns all ERC20 transfers of given address between given block range.
func (f *EtherscanTransactionFetcher) ERC20Transfer(addr ethereum.Address, from, to *big.Int, offset int) ([]common.ERC20Transfer, error) {
	var (
		logger   = f.sugar.With("func", caller.GetCurrentFunctionName())
		temp     sync.Map
		g        errgroup.Group
		throttle = make(chan int, maxConcurrentRequest)
	)
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

	results, err := f.fetch(fn, addr, from, to, offset)
	if err != nil {
		return nil, err
	}

	transfers := make([]common.ERC20Transfer, len(results))
	for i, v := range results {
		transfer := v.(etherscan.ERC20Transfer)
		throttle <- 1
		g.Go(func(index int) func() error {
			return func() error {
				transferTx, err := common.EtherscanERC20TransferToCommon(transfer, f.ethClient, f.sugar, throttle)
				if err != nil {
					logger.Errorw("convert from Etherscan erc20 tx to common failed", "error", err, "tx", transfer.Hash)
					return err
				}
				temp.Store(index, transferTx)
				return nil
			}
		}(i),
		)
	}
	if err = g.Wait(); err != nil {
		return transfers, err
	}
	temp.Range(func(index, value interface{}) bool {
		i, ok := index.(int)
		if !ok {
			err = fmt.Errorf("index (%v) cannot be asserted to integer", index)
			return false
		}
		tx, ok := value.(common.ERC20Transfer)
		if !ok {
			err = fmt.Errorf("value (%v) cannot be asserted to common.InternalTx", value)
			return false
		}
		transfers[i] = tx
		return true
	})
	if err != nil {
		return transfers, err
	}
	return transfers, nil
}
