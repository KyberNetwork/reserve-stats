package common

import (
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// BurnAssignedFeesEvent is the event emit from burner contract of KNC token.
// https://etherscan.io/txs?ea=0xed4f53268bfdff39b36e8786247ba3a02cf34b04&topic0=0x2f8d2d194cbe1816411754a2fc9478a11f0707da481b11cff7c69791eb877ee1
type BurnAssignedFeesEvent struct {
	BlockNumber uint64
	TxHash      ethereum.Hash
	Reserve     ethereum.Address
	Sender      ethereum.Address
	Quantity    *big.Int
}

func (bae *BurnAssignedFeesEvent) String() string {
	return fmt.Sprintf("BurnAssignedFeesEvent blocknumber=%d tx_hash=%s reserve=%s sender=%s quantity=%s",
		bae.BlockNumber,
		bae.TxHash.String(),
		bae.Reserve.String(),
		bae.Sender.String(),
		bae.Quantity.String())
}
