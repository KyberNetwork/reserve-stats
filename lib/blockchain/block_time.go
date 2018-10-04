package blockchain

import (
	"context"
	"sort"
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
	mu        *sync.RWMutex
	ethClient *ethclient.Client // eth client
	sugar     *zap.SugaredLogger

	cachedHeaders   map[uint64]*types.Header
	cachedIndices   []uint64
	maxCachedBlocks int
}

// NewBlockTimeResolver returns BlockTimeResolver instance given a ethereum client.
func NewBlockTimeResolver(sugar *zap.SugaredLogger, client *ethclient.Client) (*BlockTimeResolver, error) {
	// maxCachedBlocks is the maximum number of block headers in cache, must > 1.
	const maxCachedBlocks = 10

	return &BlockTimeResolver{
		mu:        &sync.RWMutex{},
		ethClient: client,
		sugar:     sugar,

		cachedHeaders:   make(map[uint64]*types.Header, maxCachedBlocks),
		maxCachedBlocks: maxCachedBlocks,
	}, nil
}

func sortUint64s(a []uint64) {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
}

// Resolve returns timestamp from block number.
// It cachedHeaders block number and block header to reduces the number of request
// to node.
func (btr *BlockTimeResolver) Resolve(blockNumber uint64) (time.Time, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// cache hit happy path
	btr.mu.RLock()
	header, ok := btr.cachedHeaders[blockNumber]
	if ok {
		btr.sugar.Debugw("block timestamp resolver cache hit", "block_number", blockNumber)
		ts := time.Unix(header.Time.Int64(), 0).UTC()
		btr.mu.RUnlock()
		return ts, nil
	}
	btr.mu.RUnlock()

	// cache miss
	btr.sugar.Debugw("block timestamp resolver cache miss", "block_number", blockNumber)
	btr.mu.Lock()
	defer btr.mu.Unlock()

	header, err := btr.ethClient.HeaderByNumber(timeout, big.NewInt(int64(blockNumber)))
	if err != nil && header == nil {
		return time.Unix(0, 0), err
	}

	// error because parity and geth are not compatible in mix hash
	// so we ignore it as we can still get time from block
	if err != nil && strings.Contains(err.Error(), "missing required field") {
		btr.sugar.Infof("ignore block header error: %s", "err", err)
	}

	if len(btr.cachedHeaders) >= btr.maxCachedBlocks {
		oldestBlockNumber := btr.cachedIndices[0]
		btr.sugar.Debugw("purging oldest cached header",
			"block_number", oldestBlockNumber,
			"max_cached_blocks", btr.maxCachedBlocks)
		btr.cachedIndices = btr.cachedIndices[1:]
		delete(btr.cachedHeaders, oldestBlockNumber)
	}
	btr.cachedIndices = append(btr.cachedIndices, blockNumber)
	sortUint64s(btr.cachedIndices)
	btr.cachedHeaders[blockNumber] = header

	return time.Unix(header.Time.Int64(), 0).UTC(), nil
}
