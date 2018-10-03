package blockchain

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TxTime is a helper to get transaction timestamp from block number and
// block header, it has a cache for one block
type TxTime struct {
	cachedBlockNo     uint64            // cache 1 block number
	cachedBlockHeader *types.Header     // cache 1 block header
	EthClient         *ethclient.Client // eth client
}

// InterpretTimestamp returns timestamp from block number and transaction index.
// It cached block number and block header to reduces the number of request
// to node.
func (txTime *TxTime) InterpretTimestamp(blockno uint64, txindex uint) (time.Time, error) {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var block *types.Header
	var err error
	if txTime.cachedBlockNo == blockno {
		block = txTime.cachedBlockHeader
	} else {
		block, err = txTime.EthClient.HeaderByNumber(timeout, big.NewInt(int64(blockno)))
	}
	if err != nil {
		if block == nil {
			return time.Unix(0, 0), err
		}

		// error because parity and geth are not compatible in mix hash
		// so we ignore it
		txTime.cachedBlockNo = blockno
		txTime.cachedBlockHeader = block
		err = nil
	}

	unixSecond := block.Time.Int64()
	unixNano := int64(txindex)
	return time.Unix(unixSecond, unixNano).UTC(), nil
}
