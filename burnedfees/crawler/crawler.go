package crawler

import (
	"context"
	"errors"
	"math/big"

	"github.com/KyberNetwork/reserve-stats/burnedfees/common"
	"github.com/KyberNetwork/reserve-stats/burnedfees/storage"
	"github.com/KyberNetwork/reserve-stats/lib/mathutil"
	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// burnAssignedFeesTopic is the topic of BurnAssignedFees (index_topic_1 address reserve, address sender, uint256 quantity)
const burnAssignedFeesTopic = "0x2f8d2d194cbe1816411754a2fc9478a11f0707da481b11cff7c69791eb877ee1"

// BurnedFeesCrawler is the crawler that tracks BurnAssignedFees events of burners contracts of KNC token.
type BurnedFeesCrawler struct {
	sugar     *zap.SugaredLogger
	ethClient *ethclient.Client
	st        storage.Interface
	burners   []ethcommon.Address
}

// NewBurnedFeesCrawler creates new instance of BurnedFeesCrawler.
func NewBurnedFeesCrawler(sugar *zap.SugaredLogger, ethClient *ethclient.Client, st storage.Interface, burners []ethcommon.Address) *BurnedFeesCrawler {
	return &BurnedFeesCrawler{
		sugar:     sugar,
		ethClient: ethClient,
		st:        st,
		burners:   burners,
	}
}

func logDataToBurnedFeesParams(data []byte) (ethcommon.Address, *big.Int, error) {
	var (
		sender ethcommon.Address
		qty    *big.Int
	)

	if len(data) != 64 {
		return sender, qty, errors.New("invalid assigned burned fees data")
	}
	sender = ethcommon.BytesToAddress(data[0:32])
	qty = ethcommon.BytesToHash(data[32:64]).Big()
	return sender, qty, nil
}

// crawl returns the BurnAssignedEvent from given block range.
func (c *BurnedFeesCrawler) crawl(fromBlock, toBlock uint64) ([]common.BurnAssignedFeesEvent, error) {
	var (
		logger = c.sugar.With(
			"func", "burnedfees/crawler/BurnedFeesCrawler.crawl",
			"from_block", fromBlock,
			"to_block", toBlock,
		)
		events []common.BurnAssignedFeesEvent
	)

	logger.Debugw("fetching BurnAssignedFees event logs",
		"topic", burnAssignedFeesTopic,
		"addresses", c.burners,
	)

	logItems, err := c.ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(0).SetUint64(fromBlock),
		ToBlock:   big.NewInt(0).SetUint64(toBlock),
		Topics:    [][]ethcommon.Hash{{ethcommon.HexToHash(burnAssignedFeesTopic)}},
		Addresses: c.burners,
	})

	if err != nil {
		return nil, err
	}

	for _, logItem := range logItems {
		sender, qty, fErr := logDataToBurnedFeesParams(logItem.Data)
		if fErr != nil {
			return nil, fErr
		}
		event := common.BurnAssignedFeesEvent{
			BlockNumber: logItem.BlockNumber,
			TxHash:      logItem.TxHash,
			Reserve:     ethcommon.BytesToAddress(logItem.Topics[1].Bytes()),
			Sender:      sender,
			Quantity:    qty,
		}
		events = append(events, event)
	}

	return events, nil
}

// Crawl is the same as crawl but split to each maxBlocks for each request.
func (c *BurnedFeesCrawler) Crawl(fromBlock, toBlock, maxBlocks uint64) error {
	var (
		logger = c.sugar.With(
			"func", "burnedfees/crawler/BurnedFeesCrawler.Crawl",
			"from_block", fromBlock,
			"to_block", toBlock,
		)
	)

	logger.Debugw("fetching BurnAssignedFees event logs")

	for i := fromBlock; i < toBlock; i += maxBlocks {
		events, err := c.crawl(i, mathutil.MintUint64(i+maxBlocks, toBlock))
		if err != nil {
			return err
		}
		if err = c.st.Store(events); err != nil {
			return err
		}
	}
	return nil
}
