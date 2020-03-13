package postgres

import (
	"database/sql"
	"errors"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

const (
	schema = `
	CREATE TABLE IF NOT EXISTS "reserve_rates" (
		id SERIAL,
		reserve TEXT NOT NULL,
		pair TEXT NOT NULL,
		buy_rate FLOAT NOT NULL,
		sell_rate FLOAT NOT NULL,
		buy_sanity_rate FLOAT NOT NULL,
		sell_sanity_rate FLOAT NOT NULL,
		from_block INTEGER NOT NULL,
		to_block INTEGER NOT NULL,
		timestamp TIMESTAMP,
		PRIMARY KEY (reserve, pair, from_block)
	);
	`
)

//Storage is postgres storage
type Storage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	blkTimeRsv blockchain.BlockTimeResolverInterface
}

// NewPostgresStorage return new storage
func NewPostgresStorage(db *sqlx.DB, sugar *zap.SugaredLogger, blkTimeRsv blockchain.BlockTimeResolverInterface) (*Storage, error) {
	if _, err := db.Exec(schema); err != nil {
		sugar.Errorw("failed to init database", "error", err)
		return nil, err
	}
	return &Storage{
		db:         db,
		sugar:      sugar,
		blkTimeRsv: blkTimeRsv,
	}, nil
}

// UpdateRatesRecords save rate records to db
func (s *Storage) UpdateRatesRecords(blockNumber uint64, rateRecords map[string]map[string]common.ReserveRateEntry) error {
	var (
		logger = s.sugar.With(
			"func", caller.GetCurrentFunctionName(),
		)
		reserves, pairs                                      []string
		buyRates, sellRates, buySanityRates, sellSanityRates []float64
		fromBlocks, toBlocks                                 []uint64
		timestamps                                           []time.Time
	)
	query := `INSERT INTO reserve_rates 
	(reserve, pair, buy_rate, sell_rate, 
		buy_sanity_rate, sell_sanity_rate, from_block, to_block, timestamp)
	VALUES(
		UNNEST($1::TEXT[]), 
		UNNEST($2::TEXT[]),
		UNNEST($3::FLOAT[]),
		UNNEST($4::FLOAT[]),
		UNNEST($5::FLOAT[]),
		UNNEST($6::FLOAT[]),
		UNNEST($7::INTEGER[]),
		UNNEST($8::INTEGER[]),
		UNNEST($9::TIMESTAMP[])
	) ON CONFLICT (reserve, pair, from_block) DO UPDATE SET from_block = EXCLUDED.from_block, 
	to_block = EXCLUDED.to_block, timestamp = EXCLUDED.timestamp;`
	if s.blkTimeRsv == nil {
		return errors.New("block time resolver is not available")
	}
	for rsvAddr, rateRecord := range rateRecords {
		lastRates, fErr := s.lastRates(rsvAddr)
		if fErr != nil {
			return fErr
		}
		for pair, rate := range rateRecord {
			var (
				fromBlock uint64
				toBlock   uint64
			)
			lastRate := lastRates[pair].Rate
			lastFromBlock := lastRates[pair].FromBlock
			lastToBlock := lastRates[pair].ToBlock
			switch {
			case lastRate == nil:
				logger.Debugw("no last rate available",
					"reserve_addr", rsvAddr,
					"pair", pair)
				fromBlock = blockNumber
				toBlock = fromBlock + 1
			case *lastRate != rate:
				logger.Debugw("rate changed, starting new rate group",
					"reserve_addr", rsvAddr,
					"last_to_block", lastToBlock,
					"pair", pair,
					"last_rate", lastRate,
					"rate", rate,
				)
				fromBlock = blockNumber
				toBlock = fromBlock + 1
			default:
				logger.Debugw("rate is remain the same as last stored record",
					"reserve_addr", rsvAddr,
					"last_to_block", lastToBlock,
					"pair", pair)
				fromBlock = lastFromBlock
				toBlock = lastToBlock + 1
			}
			timestamp, err := s.blkTimeRsv.Resolve(fromBlock)
			if err != nil {
				return err
			}

			// append a record
			reserves = append(reserves, rsvAddr)
			pairs = append(pairs, pair)
			buyRates = append(buyRates, rate.BuyReserveRate)
			sellRates = append(sellRates, rate.SellReserveRate)
			buySanityRates = append(buySanityRates, rate.BuySanityRate)
			sellSanityRates = append(sellSanityRates, rate.SellSanityRate)
			fromBlocks = append(fromBlocks, fromBlock)
			toBlocks = append(toBlocks, toBlock)
			timestamps = append(timestamps, timestamp)
			logger.Debugw("fromBlock, toBlock has the same rate", "from block", fromBlock, "to block", toBlock)
		}
	}
	if _, err := s.db.Exec(query, pq.StringArray(reserves), pq.StringArray(pairs), pq.Array(buyRates),
		pq.Array(sellRates), pq.Array(buySanityRates), pq.Array(sellSanityRates), pq.Array(fromBlocks), pq.Array(toBlocks), pq.Array(timestamps)); err != nil {
		return err
	}
	return nil
}

type lastRatesResponse struct {
	Pair           string  `db:"pair"`
	BuyRate        float64 `db:"buy_rate"`
	SellRate       float64 `db:"sell_rate"`
	BuySanityRate  float64 `db:"buy_sanity_rate"`
	SellSanityRate float64 `db:"sell_sanity_rate"`
	FromBlock      uint64  `db:"from_block"`
	ToBlock        uint64  `db:"to_block"`
	Rk             uint64  `db:"rk"` // row number used for find latest rates for each pair of token
}

func (s *Storage) lastRates(rsvAddr string) (map[string]common.LastRate, error) {
	var (
		logger = s.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"reserve_addr", rsvAddr,
		)
		lastRates   = make(map[string]common.LastRate)
		lastRatesDB []lastRatesResponse
	)
	query := `
	WITH latest as (
		SELECT buy_rate, sell_rate, buy_sanity_rate, sell_sanity_rate, from_block, to_block, pair,
		       ROW_NUMBER() OVER (PARTITION BY pair ORDER BY timestamp DESC) as rk
		FROM reserve_rates
		WHERE reserve = $1
	) 
	SELECT latest.*
	FROM latest WHERE latest.rk = 1
	`
	logger.Infow("get last rates", "query", query)
	if err := s.db.Select(&lastRatesDB, query, rsvAddr); err != nil {
		return lastRates, err
	}

	for _, rate := range lastRatesDB {
		lastRates[rate.Pair] = common.LastRate{
			FromBlock: rate.FromBlock,
			ToBlock:   rate.ToBlock,
			Rate: &common.ReserveRateEntry{
				BuyReserveRate:  rate.BuyRate,
				SellReserveRate: rate.SellRate,
				BuySanityRate:   rate.BuySanityRate,
				SellSanityRate:  rate.SellSanityRate,
			},
		}
	}
	return lastRates, nil
}

type ratesQueryResponse struct {
	ID             uint64    `db:"id"`
	Reserve        string    `db:"reserve"`
	Pair           string    `db:"pair"`
	BuyRate        float64   `db:"buy_rate"`
	SellRate       float64   `db:"sell_rate"`
	BuySanityRate  float64   `db:"buy_sanity_rate"`
	SellSanityRate float64   `db:"sell_sanity_rate"`
	FromBlock      uint64    `db:"from_block"`
	ToBlock        uint64    `db:"to_block"`
	Timestamp      time.Time `db:"timestamp"`
}

// GetRatesByTimePoint return rates by from time and to time
func (s *Storage) GetRatesByTimePoint(addrs []ethereum.Address, fromTime, toTime uint64) (map[string]map[string][]common.ReserveRates, error) {
	var (
		result = make(map[string]map[string][]common.ReserveRates)
		logger = s.sugar.With(
			"func", caller.GetCurrentFunctionName(),
			"from", fromTime,
			"to", toTime,
		)
		reserves     []string
		rateResponse []ratesQueryResponse
	)
	for _, addr := range addrs {
		reserves = append(reserves, addr.Hex())
	}
	logger.With("reserve", reserves)
	query := `SELECT * FROM reserve_rates WHERE EXTRACT(EPOCH FROM timestamp)*1000 > $1 AND EXTRACT (EPOCH FROM timestamp)*1000 < $2 AND reserve = any($3::TEXT[])`
	logger.Infow("get rates by time point", "query", query)
	if err := s.db.Select(&rateResponse, query, fromTime, toTime, pq.StringArray(reserves)); err != nil {
		return result, err
	}
	for _, rate := range rateResponse {
		if len(result) == 0 {
			result[rate.Reserve] = make(map[string][]common.ReserveRates)
		}
		ratePair := result[rate.Reserve]
		if ratePair == nil {
			ratePair = make(map[string][]common.ReserveRates)
		}
		ratePair[rate.Pair] = append(ratePair[rate.Pair], common.ReserveRates{
			Timestamp: rate.Timestamp,
			FromBlock: rate.FromBlock,
			ToBlock:   rate.ToBlock,
			Rates: common.ReserveRateEntry{
				BuyReserveRate:  rate.BuyRate,
				SellReserveRate: rate.SellRate,
				BuySanityRate:   rate.BuySanityRate,
				SellSanityRate:  rate.SellSanityRate,
			},
		})
		result[rate.Reserve] = ratePair
	}
	return result, nil
}

// LastBlock return last block saved in db
func (s *Storage) LastBlock() (int64, error) {
	var (
		lastBlock int64
		logger    = s.sugar.With("func", caller.GetCallerFunctionName())
	)
	query := `SELECT to_block FROM reserve_rates ORDER BY timestamp DESC LIMIT 1;`
	logger.Infow("Getting last block stored in db", "query", query)
	if err := s.db.Get(&lastBlock, query); err != nil {
		if err == sql.ErrNoRows {
			return lastBlock, nil
		}
		return lastBlock, err
	}
	return lastBlock, nil
}
