package influx

import (
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenUpdate(t *testing.T) {
	const (
		dbName = "testInfluxTokenUpdate"
	)
	var (
		ethAddress = ethereum.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee").Hex()
	)
	is, err := newTestInfluxStorage(dbName)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, is.tearDown())
	}()
	require.NoError(t, loadTestData(dbName))

	symbol, err := is.GetTokenSymbol(ethAddress)
	assert.NoError(t, err)
	assert.Equal(t, "", symbol)

	err = is.UpdateTokens([]string{ethAddress}, []string{"ETH"})
	assert.NoError(t, err)

	symbol, err = is.GetTokenSymbol(ethAddress)
	assert.NoError(t, err)
	assert.Equal(t, "ETH", symbol)
}
