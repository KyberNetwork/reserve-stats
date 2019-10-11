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
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
)

// TradeLogDB is storage of tradelog data
type TradeLogDB struct {
	sugar                *zap.SugaredLogger
	db                   *sqlx.DB
	tokenAmountFormatter blockchain.TokenAmountFormatterInterface
}

//NewTradeLogDB create a new instance of TradeLogDB
func NewTradeLogDB(sugar *zap.SugaredLogger, db *sqlx.DB, tokenAmountFormatter blockchain.TokenAmountFormatterInterface) (*TradeLogDB, error) {
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

func (tldb *TradeLogDB) saveReserveAddress(tx *sqlx.Tx, reserveAddressArray []string) error {
	var logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	query := fmt.Sprintf(insertionAddressTemplate, schema.ReserveTableName)
	logger.Debugw("updating rsv...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(reserveAddressArray))
	return err
}

func (tldb *TradeLogDB) saveTokens(tx *sqlx.Tx, tokensArray []string) error {
	var logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	query := fmt.Sprintf(insertionAddressTemplate, schema.TokenTableName)
	logger.Debugw("updating rsv...", "query", query)
	_, err := tx.Exec(query, pq.StringArray(tokensArray))
	return err
}

// SaveTradeLogs persist trade logs to DB
func (tldb *TradeLogDB) SaveTradeLogs(logs []common.TradeLog) (err error) {
	var (
		logger              = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		reserveAddress      = make(map[string]struct{})
		reserveAddressArray []string
		tokens              = make(map[string]struct{})
		tokensArray         []string
		records             []*record

		users = make(map[ethereum.Address]struct{})
	)
	for _, log := range logs {
		r, err := tldb.recordFromTradeLog(log)
		if err != nil {
			return err
		}

		if _, ok := users[log.UserAddress]; ok {
			r.IsFirstTrade = false
		} else {
			isFirstTrade, err := tldb.isFirstTrade(log.UserAddress)
			if err != nil {
				return err
			}
			r.IsFirstTrade = isFirstTrade
		}
		records = append(records, r)
		users[log.UserAddress] = struct{}{}
	}

	for _, r := range records {
		reserve := r.SrcReserveAddress
		if _, ok := reserveAddress[reserve]; !ok {
			reserveAddress[reserve] = struct{}{}
			reserveAddressArray = append(reserveAddressArray, reserve)
		}
		reserve = r.DstReserveAddress
		if _, ok := reserveAddress[reserve]; !ok {
			reserveAddress[reserve] = struct{}{}
			reserveAddressArray = append(reserveAddressArray, reserve)
		}
		token := r.SrcAddress
		if _, ok := tokens[reserve]; !ok {
			tokens[token] = struct{}{}
			tokensArray = append(tokensArray, token)
		}
		token = r.DestAddress
		if _, ok := tokens[reserve]; !ok {
			tokens[token] = struct{}{}
			tokensArray = append(tokensArray, token)
		}
	}

	tx, err := tldb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)

	err = tldb.saveReserveAddress(tx, reserveAddressArray)
	if err != nil {
		return err
	}

	err = tldb.saveTokens(tx, tokensArray)
	if err != nil {
		return err
	}

	for _, r := range records {
		logger.Debugw("Record", "record", r)
		_, err = tx.NamedExec(insertionUserTemplate, r)
		if err != nil {
			logger.Debugw("Error while add users", "error", err)
			return err
		}

		_, err = tx.NamedExec(insertionWalletTemplate, r)
		if err != nil {
			logger.Debugw("Error while add wallet", "error", err)
			return err
		}

		_, err = tx.NamedExec(insertionTradelogsTemplate, r)
		if err != nil {
			logger.Debugw("Error while add trade logs", "error", err)
			return err
		}
	}

	return err
}

func (tldb *TradeLogDB) isFirstTrade(userAddr ethereum.Address) (bool, error) {
	query := `SELECT NOT EXISTS(SELECT NULL FROM "` + schema.UserTableName + `" WHERE address=$1);`
	row := tldb.db.QueryRow(query, userAddr.Hex())
	var result bool
	if err := row.Scan(&result); err != nil {
		tldb.sugar.Error(err)
		return false, err
	}
	return result, nil
}

func (tldb *TradeLogDB) LoadTradeLogs(from, to time.Time) ([]common.TradeLog, error) {
	var logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	var queryResult []struct {
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
	}
	err := tldb.db.Select(&queryResult, selectTradeLogsQuery, from, to)
	if err != nil {
		return nil, err
	}

	if len(queryResult) == 0 {
		logger.Debugw("empty result returned", "query", selectTradeLogsQuery)
		return nil, nil
	}

	result := make([]common.TradeLog, 0)

	for _, r := range queryResult {

		var (
			ethAmountInWei         *big.Int
			srcAmountInWei         *big.Int
			dstAmountInWei         *big.Int
			originalEthAmountInWei *big.Int
		)

		if ethAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.EthAmount); err != nil {
			return nil, err
		}

		if originalEthAmountInWei, err = tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, r.OriginalEthAmount); err != nil {
			return nil, err
		}
		SrcAddress := ethereum.HexToAddress(r.SrcAddress)
		if srcAmountInWei, err = tldb.tokenAmountFormatter.ToWei(SrcAddress, r.SrcAmount); err != nil {
			return nil, err
		}
		DstAddress := ethereum.HexToAddress(r.DstAddress)
		if dstAmountInWei, err = tldb.tokenAmountFormatter.ToWei(DstAddress, r.SrcAmount); err != nil {
			return nil, err
		}

		tradeLog := common.TradeLog{
			TransactionHash:   ethereum.HexToHash(r.TxHash),
			Index:             r.LogIndex,
			Timestamp:         r.Timestamp,
			BlockNumber:       r.BlockNumber,
			EthAmount:         ethAmountInWei,
			OriginalEthAmount: originalEthAmountInWei,
			UserAddress:       ethereum.HexToAddress(r.UserAddress),
			SrcAddress:        SrcAddress,
			DestAddress:       DstAddress,
			SrcAmount:         srcAmountInWei,
			DestAmount:        dstAmountInWei,
			SrcReserveAddress: ethereum.HexToAddress(r.SrcReserveAddress),
			DstReserveAddress: ethereum.HexToAddress(r.DstReserveAddress),
			IP:                r.IP.String,
			Country:           r.Country.String,
			IntegrationApp:    r.IntegrationApp,
			FiatAmount:        r.EthAmount * r.EthUsdRate,
			TxSender:          ethereum.HexToAddress(r.TxSender),
			ReceiverAddress:   ethereum.HexToAddress(r.ReceiverAddr),
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
const insertionTradelogsTemplate string = `
INSERT INTO "` + schema.TradeLogsTableName + `"(
	timestamp,
 	block_number,
 	tx_hash,
 	eth_amount,
	original_eth_amount,
 	user_address_id,
 	src_address_id,
 	dst_address_id,
 	src_reserve_address_id,
 	dst_reserve_address_id,
 	src_amount,
 	dst_amount,
 	wallet_address_id,
 	src_burn_amount,
 	dst_burn_amount,
 	src_wallet_fee_amount,
 	dst_wallet_fee_amount,
 	integration_app,
 	ip,
 	country,
 	eth_usd_rate,
 	eth_usd_provider,
	index,
	kyced,
	is_first_trade,
	tx_sender,
	receiver_address
) VALUES (
 	:timestamp,
 	:block_number,
 	:tx_hash,
 	:eth_amount,
	:original_eth_amount,
 	(SELECT id FROM users WHERE address=:user_address),
 	(SELECT id FROM token WHERE address=:src_address),
 	(SELECT id FROM token WHERE address=:dst_address),
 	(SELECT id FROM reserve WHERE address=:src_reserve_address),
 	(SELECT id FROM reserve WHERE address=:dst_reserve_address),
 	:src_amount,
 	:dst_amount,
 	(SELECT id FROM wallet WHERE address=:wallet_address),
 	:src_burn_amount,
 	:dst_burn_amount,
 	:src_wallet_fee_amount,
 	:dst_wallet_fee_amount,
 	:integration_app,
 	:ip,
 	:country,
 	:eth_usd_rate,
 	:eth_usd_provider,
 	:index,
	:kyced,
	:is_first_trade,
	:tx_sender,
 	:receiver_address
)
ON CONFLICT
DO NOTHING;`

const selectTradeLogsQuery = `
SELECT a.timestamp AS timestamp, block_number, eth_amount, original_eth_amount, eth_usd_rate, d.address AS user_address,
e.address AS src_address, f.address AS dst_address,
src_amount, dst_amount, ip, country, integration_app, src_burn_amount, dst_burn_amount,
index, tx_hash, b.address AS src_rsv_address, c.address AS dst_rsv_address, src_wallet_fee_amount, dst_wallet_fee_amount,
g.address AS wallet_addr, tx_sender, receiver_address
FROM "` + schema.TradeLogsTableName + `" AS a
INNER JOIN reserve AS b ON a.src_reserve_address_id = b.id
INNER JOIN reserve AS c ON a.dst_reserve_address_id = c.id
INNER JOIN users AS d ON a.user_address_id = d.id
INNER JOIN token AS e ON a.src_address_id = e.id
INNER JOIN token AS f ON a.dst_address_id = f.id
INNER JOIN wallet AS g ON a.wallet_address_id = g.id
WHERE a.timestamp >= $1 and a.timestamp <= $2;
`