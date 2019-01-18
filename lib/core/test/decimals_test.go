package test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/contracts"
	"github.com/KyberNetwork/reserve-stats/lib/core"
	"github.com/KyberNetwork/reserve-stats/lib/testutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const defaultEthereumNode = "https://mainnet.infura.io"

var timeout = 30 * time.Second

type tokenResult struct {
	address string
	decimal int64
}

// TestTokenDecimals asserts that token decimals configured from Kyber Core has the same values as by calling the
// (optional) constant of the token contract directly.
func TestTokenDecimals(t *testing.T) {
	testMode, err := testutil.GetTestMode()
	if err != nil {
		t.Skip("Can't get test mode. skip this external test")
	}
	if testMode != testutil.External {
		t.Skip("disable as this test require external resource")
	}
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
		tokensFromCore       = make(map[string]int64, len(listTokens)-1) // exclude eth
		tokensFromBlockchain = make(map[string]int64, len(listTokens)-1) // exclude eth
	)

	for _, token := range listTokens {
		if token.ID != "ETH" {
			tokensFromCore[strings.ToLower(token.Address)] = token.Decimals
		}
	}

	client, err := ethclient.Dial(node)
	require.NoError(t, err, "Ethereum client should init successfully")
	var (
		g           errgroup.Group
		resourcesCh = make(chan struct{}, 10)                   // resources limiter, thread need to acquire release resource
		resultCh    = make(chan tokenResult, len(listTokens)-1) // exclude eth
	)

	for _, token := range listTokens {
		if token.ID == "ETH" {
			continue
		}
		var address = common.HexToAddress(token.Address)
		g.Go(
			func() error {
				resourcesCh <- struct{}{}
				defer func() { <-resourcesCh }()

				tokenContract, err := contracts.NewERC20(address, client)
				if err != nil {
					sugar.Debugw("Could not create erc20 contract", "token", address, "error", err)
					return err
				}
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				decimals, err := tokenContract.Decimals(&bind.CallOpts{Context: ctx})
				if err != nil {
					sugar.Debugw("Could not get decimals", "token", address, "error", err)
					return err
				}
				resultCh <- tokenResult{strings.ToLower(address.Hex()), int64(decimals)}
				return nil
			})
	}
	if err = g.Wait(); err != nil {
		require.NoError(t, err, "Could not get decimal from blockchain")
	}
	close(resultCh)
	for r := range resultCh {
		tokensFromBlockchain[r.address] = r.decimal
	}
	assert.Equal(t, len(tokensFromCore), len(tokensFromBlockchain), "Len")
	for token, decimal := range tokensFromCore {
		assert.Equal(t, decimal, tokensFromBlockchain[token], "Decimal of %s is not equal", token)
	}
}
