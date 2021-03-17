package postgres

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// TradeLogDB is storage of tradelog data
type TradeLogDB struct {
	sugar                *zap.SugaredLogger
	db                   *sqlx.DB
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface

	// used for calculate burn amount
	// as different environment have different knc address
	kncAddr ethereum.Address
}

//NewTradeLogDB create a new instance of TradeLogDB
func NewTradeLogDB(sugar *zap.SugaredLogger, db *sqlx.DB, tokenAmountFormatter blockchain.TokenAmountFormatterInterface, kncAddr ethereum.Address) (*TradeLogDB, error) {
	var logger = sugar.With("func", caller.GetCurrentFunctionName())
	var err error
	logger.Debug("initializing database schema")
	if _, err = db.Exec(schema.TradeLogsSchema); err != nil {
		return nil, err
	}
	logger.Debug("database schema initialized successfully")

	return &TradeLogDB{
		sugar:                sugar,
		db:                   db,
		tokenAmountFormatter: tokenAmountFormatter,
		kncAddr:              kncAddr,
	}, err
}

// LastBlock returns last stored trade log block number from database.
func (tldb *TradeLogDB) LastBlock() (int64, error) {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		result sql.NullInt64
	)
	stmt := fmt.Sprintf(`SELECT MAX("block_number") FROM "%v"`, schema.TradeLogsTableName)
	logger = logger.With("query", stmt)
	logger.Debug("Start query")
	err := tldb.db.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		logger.Errorw("Get error ", "error", err)
		return 0, err
	}
	return result.Int64, nil
}

type tradeLogDBData struct {
	ID                uint64         `db:"id"`
	Timestamp         time.Time      `db:"timestamp"`
	BlockNumber       uint64         `db:"block_number"`
	EthAmount         float64        `db:"eth_amount"`
	OriginalEthAmount float64        `db:"original_eth_amount"`
	EthUsdRate        float64        `db:"eth_usd_rate"`
	UserAddress       pq.StringArray `db:"user_address"`
	SrcAddress        pq.StringArray `db:"src_address"`
	DstAddress        pq.StringArray `db:"dst_address"`
	SrcAmount         float64        `db:"src_amount"`
	DstAmount         float64        `db:"dst_amount"`
	LogIndex          uint           `db:"index"`
	TxHash            string         `db:"tx_hash"`
	TxSender          string         `db:"tx_sender"`
	ReceiverAddr      string         `db:"receiver_address"`
	GasUsed           uint64         `db:"gas_used"`
	GasPrice          float64        `db:"gas_price"`
	TransactionFee    float64        `db:"transaction_fee"`
	Version           uint           `db:"version"`
}

func (tldb *TradeLogDB) tradeLogFromDBData(r tradeLogDBData) (common.Tradelog, error) {
	var (
		tradeLog common.Tradelog
		err      error

		ethAmountInWei                     *big.Int
		srcAmountInWei                     *big.Int
		dstAmountInWei                     *big.Int
		originalEthAmountInWei             *big.Int
		gasPriceInWei, transactionFeeInWei *big.Int

		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	)

	if ethAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.EthAmount); err != nil {
		logger.Debugw("failed to parse eth amount", "error", err)
		return tradeLog, err
	}

	if originalEthAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.OriginalEthAmount); err != nil {
		logger.Debugw("failed to parse original eth amount", "error", err)
		return tradeLog, err
	}
	SrcAddress := ethereum.HexToAddress(r.SrcAddress[0])
	if srcAmountInWei, err = tldb.tokenAmountFormatter.ToWei(SrcAddress, r.SrcAmount); err != nil {
		logger.Debugw("failed to parse src amount", "error", err)
		return tradeLog, err
	}
	DstAddress := ethereum.HexToAddress(r.DstAddress[0])
	if dstAmountInWei, err = tldb.tokenAmountFormatter.ToWei(DstAddress, r.DstAmount); err != nil {
		logger.Debugw("failed to parse dst amount", "error", err)
		return tradeLog, err
	}

	// these conversion below is from Gwei to wei which is used ^18 method, so I used ToWei function with ETHAddr - which have decimals of 18
	if gasPriceInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.GasPrice); err != nil {
		logger.Debugw("failed to parse gas price", "error", err)
		return tradeLog, err
	}

	if transactionFeeInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.TransactionFee); err != nil {
		logger.Debugw("failed to parse transaction fee", "error", err)
		return tradeLog, err
	}
	tradeLog = common.Tradelog{
		TransactionHash:    ethereum.HexToHash(r.TxHash),
		Index:              r.LogIndex,
		Timestamp:          r.Timestamp,
		BlockNumber:        r.BlockNumber,
		USDTAmount:         ethAmountInWei,
		OriginalUSDTAmount: originalEthAmountInWei,
		User: common.KyberUserInfo{
			UserAddress: ethereum.HexToAddress(r.UserAddress[0]),
		},
		TokenInfo: common.TradeTokenInfo{
			SrcAddress:  SrcAddress,
			DestAddress: DstAddress,
		},
		SrcAmount:       srcAmountInWei,
		DestAmount:      dstAmountInWei,
		FiatAmount:      r.EthAmount * r.EthUsdRate,
		ReceiverAddress: ethereum.HexToAddress(r.ReceiverAddr),
		ETHUSDRate:      r.EthUsdRate,
		TxDetail: common.TxDetail{
			GasUsed:        r.GasUsed,
			GasPrice:       gasPriceInWei,
			TransactionFee: transactionFeeInWei,
			TxSender:       ethereum.HexToAddress(r.TxSender),
		},
		Version: r.Version,
	}
	return tradeLog, nil
}

// LoadTradeLogsByTxHash get list of tradelogs by tx hash
func (tldb *TradeLogDB) LoadTradeLogsByTxHash(tx ethereum.Hash) ([]common.Tradelog, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult []tradeLogDBData
		result      = make([]common.Tradelog, 0)
	)
	err := tldb.db.Select(&queryResult, selectTradeLogsWithTxHashQuery, tx.Hex())
	if err != nil {
		logger.Errorw("failed to get tradelog from database", "error", err)
		return nil, err
	}

	if len(queryResult) == 0 {
		logger.Debugw("empty result returned", "query", selectTradeLogsWithTxHashQuery)
		return result, nil
	}

	for _, r := range queryResult {
		tradeLog, err := tldb.tradeLogFromDBData(r)
		if err != nil {
			logger.Errorw("cannot parse db data to trade log", "error", err)
			return nil, err
		}
		result = append(result, tradeLog)
	}
	return result, nil
}

// LoadTradeLogs get list of tradelogs by timestamp from time to time
func (tldb *TradeLogDB) LoadTradeLogs(from, to time.Time) ([]common.Tradelog, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult []tradeLogDBData
		result      = make([]common.Tradelog, 0)
	)
	err := tldb.db.Select(&queryResult, selectTradeLogsQuery, from, to)
	if err != nil {
		return nil, err
	}

	if len(queryResult) == 0 {
		logger.Debugw("empty result returned", "query", selectTradeLogsQuery)
		return result, nil
	}

	for _, r := range queryResult {
		tradeLog, err := tldb.tradeLogFromDBData(r)
		if err != nil {
			logger.Errorw("cannot parse db data to trade log", "error", err)
			return nil, err
		}
		result = append(result, tradeLog)
	}
	return result, nil
}

const insertionAddressTemplate = `INSERT INTO %[1]s(
	address
) VALUES(
	unnest($1::TEXT[])
)
ON CONFLICT ON CONSTRAINT %[1]s_address_key DO NOTHING`

const insertionUserTemplate string = `
INSERT INTO users(
	address,
	timestamp
) VALUES (
	:user_address,
	:timestamp
)
ON CONFLICT (address) 
DO NOTHING;`

const selectTradeLogsQuery = `
SELECT a.id, a.timestamp AS timestamp, a.block_number, a.eth_amount, original_eth_amount, eth_usd_rate, 
ARRAY_AGG(d.address) AS user_address,
ARRAY_AGG(e.address) AS src_address, 
ARRAY_AGG(f.address) AS dst_address,
a.src_amount, 
a.dst_amount, 
ip, country, integration_app, 
a.index, tx_hash, tx_sender, receiver_address, 
ARRAY_AGG(w.address) as wallet_address,
COALESCE(gas_used, 0) as gas_used, COALESCE(gas_price, 0) as gas_price, 
COALESCE(transaction_fee, 0) as transaction_fee, 
version,

FROM tradelogs AS a
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN wallet as w on a.wallet_address_id = w.id
WHERE a.timestamp >= $1 and a.timestamp <= $2
GROUP BY a.id;
`

const selectTradeLogsWithTxHashQuery = `
SELECT 
a.timestamp AS timestamp, 
a.block_number, 
a.eth_amount, 
original_eth_amount, 
eth_usd_rate, 
ARRAY_AGG(d.address) AS user_address,
ARRAY_AGG(e.address) AS src_address, 
ARRAY_AGG(f.address) AS dst_address,
a.src_amount, 
a.dst_amount, 
ip, 
country, 
integration_app, 
a.index, 
tx_hash, 
ARRAY_AGG(w.address) AS wallet_address, 
tx_sender, 
receiver_address, 
COALESCE(gas_used, 0) as gas_used, 
COALESCE(gas_price, 0) as gas_price, 
COALESCE(transaction_fee, 0) as transaction_fee,
version,

FROM tradelogs AS a
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN wallet AS w ON a.wallet_address_id = w.id
LEFT JOIN fee ON fee.trade_id = a.id
LEFT JOIN split ON split.trade_id = a.id
INNER JOIN reserve sr ON sr.id = split.reserve_id
WHERE a.tx_hash=$1
GROUP BY a.id;
`
