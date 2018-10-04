package blockchain

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockTimeResolver is a helper to get transaction timestamp from block number.
// It has a cache for one block.
type BlockTimeResolver struct {
	mu                *sync.RWMutex
	cachedBlockNo     uint64            // cache 1 block number
	cachedBlockHeader *types.Header     // cache 1 block header
	ethClient         *ethclient.Client // eth client
	sugar             *zap.SugaredLogger
}

// NewTxTime returns TxTime instance given a ethereum client.
func NewTxTime(sugar *zap.SugaredLogger, client *ethclient.Client) *BlockTimeResolver {
	return &BlockTimeResolver{
		mu:        &sync.RWMutex{},
		ethClient: client,
		sugar:     sugar,
	}
}

// Resolve returns timestamp from block number.
// It cached block number and block header to reduces the number of request
// to node.
func (btr *BlockTimeResolver) Resolve(blockno uint64) (time.Time, error) {
	btr.mu.Lock()
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer func() {
		cancel()
		btr.mu.Unlock()
	}()

	var block *types.Header
	var err error
	if btr.cachedBlockNo == blockno {
		block = btr.cachedBlockHeader
	} else {
		block, err = btr.ethClient.HeaderByNumber(timeout, big.NewInt(int64(blockno)))

		if err != nil {
			if block == nil {
				return time.Unix(0, 0), err
			}

		// error because parity and geth are not compatible in mix hash
		// so we ignore it as we can still get time from block
		if strings.Contains(err.Error(), "missing required field") {
			btr.sugar.Infof("Ignore block header error: %s", "err", err)
		}

		btr.cachedBlockNo = blockno
		btr.cachedBlockHeader = block
	}

	return time.Unix(block.Time.Int64(), 0).UTC(), nil
}
