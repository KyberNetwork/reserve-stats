package postgres

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
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
	ID                 uint64         `db:"id"`
	Timestamp          time.Time      `db:"timestamp"`
	BlockNumber        uint64         `db:"block_number"`
	EthAmount          float64        `db:"eth_amount"`
	OriginalEthAmount  float64        `db:"original_eth_amount"`
	EthUsdRate         float64        `db:"eth_usd_rate"`
	UserAddress        string         `db:"user_address"`
	SrcAddress         string         `db:"src_address"`
	DstAddress         string         `db:"dst_address"`
	SrcAmount          float64        `db:"src_amount"`
	DstAmount          float64        `db:"dst_amount"`
	LogIndex           uint           `db:"index"`
	TxHash             string         `db:"tx_hash"`
	IP                 sql.NullString `db:"ip"`
	Country            sql.NullString `db:"country"`
	IntegrationApp     string         `db:"integration_app"`
	SrcBurnAmount      float64        `db:"src_burn_amount"`
	DstBurnAmount      float64        `db:"dst_burn_amount"`
	SrcReserveAddress  string         `db:"src_rsv_address"`
	DstReserveAddress  string         `db:"dst_rsv_address"`
	SrcWalletFeeAmount float64        `db:"src_wallet_fee_amount"`
	DstWalletFeeAmount float64        `db:"dst_wallet_fee_amount"`
	WalletAddress      string         `db:"wallet_addr"`
	TxSender           string         `db:"tx_sender"`
	ReceiverAddr       string         `db:"receiver_address"`
	GasUsed            uint64         `db:"gas_used"`
	GasPrice           float64        `db:"gas_price"`
	TransactionFee     float64        `db:"transaction_fee"`
	Version            uint           `db:"version"`
}

type feeRecord struct {
	ID             uint64  `db:"id"`
	TradeID        uint64  `db:"trade_id"`
	ReserveAddress string  `db:"reserve_address"`
	WalletAddress  string  `db:"wallet_address"`
	WalletFee      float64 `db:"wallet_fee"`
	PlatformFee    float64 `db:"platform_fee"`
	Burn           float64 `db:"burn"`
	Rebate         float64 `db:"rebate"`
	Reward         float64 `db:"reward"`
	Version        uint64  `db:"version"`
	// RebateWallets  []string  `db:"rebatewallets"`
	// RebatePercents []float64 `db:"rebatepercents"`
}

type splitRecord struct {
	ReserveAddress string  `db:"address"`
	SrcToken       string  `db:"src"`
	DstToken       string  `db:"dst"`
	SrcAmount      float64 `db:"src_amount"`
	Rate           float64 `db:"rate"`
}

func (tldb *TradeLogDB) tradeLogFromDBData(r tradeLogDBData, f []feeRecord, s []splitRecord) (common.TradelogV4, error) {
	var (
		tradeLog common.TradelogV4
		err      error

		ethAmountInWei                     *big.Int
		srcAmountInWei                     *big.Int
		dstAmountInWei                     *big.Int
		originalEthAmountInWei             *big.Int
		gasPriceInWei, transactionFeeInWei *big.Int
		fees                               []common.TradelogFee
		split                              []common.TradeSplit

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
	SrcAddress := ethereum.HexToAddress(r.SrcAddress)
	if srcAmountInWei, err = tldb.tokenAmountFormatter.ToWei(SrcAddress, r.SrcAmount); err != nil {
		logger.Debugw("failed to parse src amount", "error", err)
		return tradeLog, err
	}
	DstAddress := ethereum.HexToAddress(r.DstAddress)
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

	for _, fee := range f {
		platformFee, err := tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, fee.PlatformFee)
		if err != nil {
			return tradeLog, err
		}
		walletFee, err := tldb.tokenAmountFormatter.ToWei(blockchain.KNCAddr, fee.WalletFee)
		if err != nil {
			return tradeLog, err
		}
		burn, err := tldb.tokenAmountFormatter.ToWei(blockchain.KNCAddr, fee.Burn)
		if err != nil {
			return tradeLog, err
		}
		rebate, err := tldb.tokenAmountFormatter.ToWei(blockchain.KNCAddr, fee.Rebate)
		if err != nil {
			return tradeLog, err
		}
		reward, err := tldb.tokenAmountFormatter.ToWei(blockchain.KNCAddr, fee.Reward)
		if err != nil {
			return tradeLog, err
		}
		fees = append(fees, common.TradelogFee{
			ReserveAddr:    ethereum.HexToAddress(fee.ReserveAddress),
			PlatformWallet: ethereum.HexToAddress(fee.WalletAddress),
			PlatformFee:    platformFee,
			WalletFee:      walletFee,
			Burn:           burn,
			Rebate:         rebate,
			Reward:         reward,
		})
	}

	for _, sp := range s {
		srcAmount, err := tldb.tokenAmountFormatter.ToWei(ethereum.HexToAddress(sp.SrcToken), sp.SrcAmount)
		if err != nil {
			return tradeLog, err
		}
		rate, err := tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, sp.Rate)
		if err != nil {
			return tradeLog, err
		}
		split = append(split, common.TradeSplit{
			ReserveAddress: ethereum.HexToAddress(sp.ReserveAddress),
			SrcToken:       ethereum.HexToAddress(sp.SrcToken),
			DstToken:       ethereum.HexToAddress(sp.DstToken),
			SrcAmount:      srcAmount,
			Rate:           rate,
		})
	}

	tradeLog = common.TradelogV4{
		TransactionHash:   ethereum.HexToHash(r.TxHash),
		Index:             r.LogIndex,
		Timestamp:         r.Timestamp,
		BlockNumber:       r.BlockNumber,
		EthAmount:         ethAmountInWei,
		OriginalEthAmount: originalEthAmountInWei,
		User: common.KyberUserInfo{
			UserAddress: ethereum.HexToAddress(r.UserAddress),
			IP:          r.IP.String,
			Country:     r.Country.String,
		},
		TokenInfo: common.TradeTokenInfo{
			SrcAddress:  SrcAddress,
			DestAddress: DstAddress,
		},
		SrcAmount:       srcAmountInWei,
		DestAmount:      dstAmountInWei,
		IntegrationApp:  r.IntegrationApp,
		FiatAmount:      r.EthAmount * r.EthUsdRate,
		WalletAddress:   ethereum.HexToAddress(r.WalletAddress),
		ReceiverAddress: ethereum.HexToAddress(r.ReceiverAddr),
		ETHUSDRate:      r.EthUsdRate,
		TxDetail: common.TxDetail{
			GasUsed:        r.GasUsed,
			GasPrice:       gasPriceInWei,
			TransactionFee: transactionFeeInWei,
			TxSender:       ethereum.HexToAddress(r.TxSender),
		},
		Fees:    fees,
		Split:   split,
		Version: r.Version,
	}
	return tradeLog, nil
}

// LoadTradeLogsByTxHash get list of tradelogs by tx hash
func (tldb *TradeLogDB) LoadTradeLogsByTxHash(tx ethereum.Hash) ([]common.TradelogV4, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult []tradeLogDBData
		result      = make([]common.TradelogV4, 0)
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
		var (
			feeResult []feeRecord
		)
		if err := tldb.db.Select(&feeResult, selectFeeByTradelogID, r.ID); err != nil {
			logger.Debugw("failed to get fee from db", "error", err)
			return result, err
		}
		tradeLog, err := tldb.tradeLogFromDBData(r, feeResult, nil)
		if err != nil {
			logger.Errorw("cannot parse db data to trade log", "error", err)
			return nil, err
		}
		result = append(result, tradeLog)
	}
	return result, nil
}

// LoadTradeLogs get list of tradelogs by timestamp from time to time
func (tldb *TradeLogDB) LoadTradeLogs(from, to time.Time) ([]common.TradelogV4, error) {
	var (
		logger      = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		queryResult []tradeLogDBData
		result      = make([]common.TradelogV4, 0)
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
		var (
			feeResult   []feeRecord
			splitResult []splitRecord
		)
		if err := tldb.db.Select(&feeResult, selectFeeByTradelogID, r.ID); err != nil {
			logger.Debugw("failed to get fee from db", "error", err)
			return nil, err
		}
		if err := tldb.db.Select(&splitResult, selectSplitByTradelogID, r.ID); err != nil {
			logger.Debugw("failed to get split from db", "error", err)
			return nil, err
		}
		tradeLog, err := tldb.tradeLogFromDBData(r, feeResult, splitResult)
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

const insertionWalletTemplate string = `
INSERT INTO wallet(
	address,
	name
) VALUES (
	:wallet_address,
	:wallet_name
)
ON CONFLICT (address) 
DO NOTHING;`

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
SELECT a.id, a.timestamp AS timestamp, a.block_number, eth_amount, original_eth_amount, eth_usd_rate, d.address AS user_address,
e.address AS src_address, f.address AS dst_address,
src_amount, dst_amount, ip, country, integration_app, 
index, tx_hash, tx_sender, receiver_address, 
COALESCE(gas_used, 0) as gas_used, COALESCE(gas_price, 0) as gas_price, COALESCE(transaction_fee, 0) as transaction_fee, version
FROM "` + schema.TradeLogsTableName + `" AS a
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
WHERE a.timestamp >= $1 and a.timestamp <= $2;
`

const selectFeeByTradelogID = `SELECT id, trade_id, reserve_address, wallet_address, wallet_fee,
platform_fee, burn, rebate, reward FROM fee WHERE trade_id = $1;`

const selectSplitByTradelogID = `SELECT reserve.address, split.src, split.dst, split.src_amount, split.rate FROM split 
JOIN reserve ON reserve.id = split.reserve_id 
WHERE trade_id = $1;`

const selectTradeLogsWithTxHashQuery = `
SELECT a.timestamp AS timestamp, a.block_number, eth_amount, original_eth_amount, eth_usd_rate, d.address AS user_address,
e.address AS src_address, f.address AS dst_address,
src_amount, dst_amount, ip, country, integration_app, src_burn_amount, dst_burn_amount,
index, tx_hash, b.address AS src_rsv_address, c.address AS dst_rsv_address, src_wallet_fee_amount, dst_wallet_fee_amount,
g.address AS wallet_addr, tx_sender, receiver_address, COALESCE(gas_used, 0) as gas_used, COALESCE(gas_price, 0) as gas_price, COALESCE(transaction_fee, 0) as transaction_fee
FROM "` + schema.TradeLogsTableName + `" AS a
INNER JOIN reserve AS b ON a.src_reserve_address_id = b.id
INNER JOIN reserve AS c ON a.dst_reserve_address_id = c.id
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN wallet AS g ON a.wallet_address_id = g.id
WHERE a.tx_hash=$1;
`
