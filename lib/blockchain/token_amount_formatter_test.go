package blockchain

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-stats/lib/testutil"
)

func TestTokenAmountFormatter(t *testing.T) {
	testutil.SkipExternal(t)
	client := testutil.MustNewDevelopmentwEthereumClient()
	f, err := NewTokenAmountFormatter(client)
	require.NoError(t, err)

	var tests = []struct {
		address  common.Address
		decimals int64
	}{
		{ETHAddr, 18},
		{ETHAddr, 18},
		{common.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"), 18}, // BNB
		{common.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"), 18}, // BNB
		{common.HexToAddress("0x6f259637dcd74c767781e37bc6133cd6a68aa161"), 18}, // HuobiToken
		{common.HexToAddress("0x6f259637dcd74c767781e37bc6133cd6a68aa161"), 18}, // HuobiToken
	}

	for _, tc := range tests {
		decimals, fErr := f.GetDecimals(tc.address)
		require.NoError(t, fErr)
		assert.Equal(t, tc.decimals, decimals)
	}
}
