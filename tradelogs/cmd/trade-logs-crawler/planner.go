package main

import (
	"context"
	"errors"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage"
)

// crawlPlanner returns next from/to block to use in trade logs crawler.
type crawlPlanner struct {
	sugar *zap.SugaredLogger

	ethClient *ethclient.Client
	st        storage.Interface

	fromBlock *big.Int
	toBlock   *big.Int

	confirmations int64 // number of block confirmations before fetching
	delay         time.Duration

	completed bool
}

// newCrawlerPlanner returns new crawler planner instance with given context.
func newCrawlerPlanner(sugar *zap.SugaredLogger, c *cli.Context, st storage.Interface) (*crawlPlanner, error) {
	var (
		fromBlock *big.Int
		toBlock   *big.Int
	)

	ethClient, err := blockchain.NewEthereumClientFromFlag(c)
	if err != nil {
		return nil, err
	}

	if c.String(fromBlockFlag) != "" {
		if fromBlock, err = app.ParseBigIntFlag(c, fromBlockFlag); err != nil {
			return nil, err
		}
	}

	if c.String(toBlockFlag) != "" {
		if toBlock, err = app.ParseBigIntFlag(c, toBlockFlag); err != nil {
			return nil, err
		}
	}
	return &crawlPlanner{
		sugar:         sugar,
		ethClient:     ethClient,
		st:            st,
		fromBlock:     fromBlock,
		toBlock:       toBlock,
		confirmations: c.Int64(blockConfirmationsFlag),
		delay:         c.Duration(delayFlag),
		completed:     false,
	}, nil
}

// Next returns next from/to block to fetch trade logs.
// This method will return io.EOF if there is no next block range to fetch.
func (p *crawlPlanner) Next() (*big.Int, *big.Int, error) {
	var (
		logger   = p.sugar.With("func", "tradelogs/cmd/trade-logs-crawler/crawlPlanner.Next")
		err      error
		dstBlock *big.Int
	)
	if p.completed {
		return nil, nil, io.EOF
	}

	if p.fromBlock == nil {
		var lastBlock int64
		if lastBlock, err = p.st.LastBlock(); err != nil {
			return nil, nil, err
		}
		if lastBlock == 0 {
			logger.Infow("using default from block number", "from_block", defaultFromBlock)
			lastBlock = defaultFromBlock
		}
		p.fromBlock = big.NewInt(lastBlock)
	}

	if p.fromBlock != nil && p.toBlock != nil && p.fromBlock.Cmp(p.toBlock) >= 0 {
		return nil, nil, errors.New("fromBlock is bigger than toBlock")
	}

	// to block is configured, terminate after first run
	if p.toBlock != nil {
		p.completed = true
		return p.fromBlock, p.toBlock, nil
	}

	// loop until there is block with enough confirmations
	for {
		var currentHeader *types.Header
		if currentHeader, err = p.ethClient.HeaderByNumber(context.Background(), nil); err != nil {
			return nil, nil, err
		}
		dstBlock = big.NewInt(0).Sub(currentHeader.Number, big.NewInt(p.confirmations))
		if p.fromBlock.Cmp(dstBlock) < 0 {
			break
		}
		logger.Infow("nothing to fetch, sleeping...",
			"from_block", p.fromBlock,
			"latest_block", currentHeader.Number,
			"confirmations", p.confirmations,
			"sleep", p.delay,
		)
		time.Sleep(p.delay)
	}

	logger.Infow("fetching trade logs up to latest known block number",
		"from_block", p.fromBlock,
		"to_block", dstBlock,
		"confirmations", p.confirmations,
	)
	fromBlock := big.NewInt(0).Set(p.fromBlock)
	p.fromBlock = dstBlock
	return fromBlock, dstBlock, nil
}
