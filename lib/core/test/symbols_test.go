package test

import (
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/core"
)

type tokenIDResult struct {
	address string
	id      string
}

// TestTokenSymbols asserts that token symbols configured from Kyber Core has the same values as by calling the
// (optional) constant of the token contract directly.
func TestTokenSymbols(t *testing.T) {
	t.Skip("skip this test as it requires reserve-core")

	logger, err := zap.NewDevelopment()
	require.Nil(t, err, "logger should be initiated successfully")
	sugar := logger.Sugar()

	url, ok := os.LookupEnv("CORE_URL")
	assert.True(t, ok)

	secret, ok := os.LookupEnv("CORE_SIGNING_KEY")
	assert.True(t, ok)

	node, ok := os.LookupEnv("ETHEREUM_NODE")
	if !ok {
		node = defaultEthereumNode
	}

	c, err := core.NewClient(sugar, url, secret)
	require.NoError(t, err, "core client should be initiated successfully")

	listTokens, err := c.Tokens()
	require.NoError(t, err, "should get list token from core successfully")

	var (
		tokensFromCore       = make(map[string]string, len(listTokens)-1) // exclude eth
		tokensFromBlockchain = make(map[string]string, len(listTokens)-1) // exclude eth
	)

	for _, token := range listTokens {
		tokensFromCore[strings.ToLower(token.Address)] = token.ID
	}

	client, err := ethclient.Dial(node)
	require.NoError(t, err, "Ethereum client should init successfully")
	symbolLookup := blockchain.NewTokenSymbol(client)
	var (
		g           errgroup.Group
		resourcesCh = make(chan struct{}, 10) // resources limiter, thread need to acquire release resource
		resultCh    = make(chan tokenIDResult, len(listTokens))
	)

	for _, token := range listTokens {
		var address = common.HexToAddress(token.Address)
		g.Go(
			func() error {
				resourcesCh <- struct{}{}
				defer func() { <-resourcesCh }()
				symbol, err := symbolLookup.Symbol(address)
				if err != nil {
					return err
				}
				sugar.Debugw("Token", "token", address, "symbol", symbol)
				resultCh <- tokenIDResult{
					strings.ToLower(address.Hex()),
					strings.ToUpper(symbol)}
				return nil
			})
	}
	if err = g.Wait(); err != nil {
		require.NoError(t, err, "Could not get decimal from blockchain")
	}
	close(resultCh)
	for r := range resultCh {
		tokensFromBlockchain[r.address] = r.id
	}
	assert.Equal(t, len(tokensFromCore), len(tokensFromBlockchain), "Len")
	for token, symbol := range tokensFromCore {
		assert.Equal(t, symbol, tokensFromBlockchain[token], "Symbol of %s is not equal", token)
	}
}
