package fetcher

import (
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
	etherscan "github.com/nanmu42/etherscan-api"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
)

// WalletFetcher is an implementation of TransactionFetcher that uses Etherscan API to fetch wallet's erc transfer
type WalletFetcher struct {
	sugar  *zap.SugaredLogger
	client *etherscan.Client
}

// NewWalletFetcher returns a new Wallet Fetcher instance.
func NewWalletFetcher(sugar *zap.SugaredLogger, client *etherscan.Client) *WalletFetcher {
	return &WalletFetcher{sugar: sugar, client: client}
}

// Fetch return all ERC20transfer from input or error if occured.
func (wf *WalletFetcher) Fetch(walletAddress ethereum.Address, from, to *big.Int) ([]common.ERC20Transfer, error) {
	const offset = 500
	var (
		logger = wf.sugar.With(
			"func",
			"accounting/reserve-transaction-fetcher/fetcher.EtherscanTransactionFetcher.fetch",
			"wallt address", walletAddress.Hex(),
			"offset", offset,
		)
		wallet  = walletAddress.String()
		results []common.ERC20Transfer
	)

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
		txs, err := wf.client.ERC20Transfers(nil, &wallet, startBlock, endBlock, page, offset)
		if blockchain.IsEtherscanNotransactionFound(err) {
			logger.Debugw("all transaction fetched", "page", page)
			break
		} else if err != nil {
			return nil, err
		}
		for _, tx := range txs {
			results = append(results, common.EtherscanERC20TransferToCommon(tx))
		}
	}

	return results, nil
}
