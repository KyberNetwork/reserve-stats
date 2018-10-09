package contracts

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// MockVersionedWrapper return mock result for testing purpose
// It won't make any actual call to blockchain
type MockVersionedWrapper struct {
}

const (
	//MockBuyRate is int64 form of mock buy rate in wei. This equal to 1 in float64 rate
	MockBuyRate = 1000000000000000000
	//MockSellRate is int64 form of mock sell rate in wei. This equal to 2 in float64 rate
	MockSellRate = 2000000000000000000
	//MockBuySanityRate is int64 form of mock sanity buy rate in wei. This equal to 3 in float64 rate
	MockBuySanityRate = 3000000000000000000
	//MockSellSanityRate is int64 form of mock sanity sell rate in wei. This equal to 4 in float64 rate
	MockSellSanityRate = 4000000000000000000
)

// GetReserveRate return mock reserveRate and SanityRate from blockchain
// it returned fixed value as dictated.
func (mvw *MockVersionedWrapper) GetReserveRate(block uint64, rsvAddr ethereum.Address, srcs, dest []ethereum.Address) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	for range srcs {
		*ret0 = append(*ret0, big.NewInt(MockSellRate))
		*ret0 = append(*ret0, big.NewInt(MockBuyRate))
		*ret1 = append(*ret1, big.NewInt(MockSellSanityRate))
		*ret1 = append(*ret1, big.NewInt(MockBuySanityRate))
	}
	return *ret0, *ret1, nil
}
