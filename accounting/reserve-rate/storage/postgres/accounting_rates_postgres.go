package postgres

import (
	"fmt"
	"sort"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/reserverates/common"
)

const (
	reserveTableName = "reserves"
	tokenTableName   = "tokens"
	baseTableName    = "bases"
	rateTableName    = "rates"
	usdTableName     = "usds"
)

//RatesStorage defines the object to store rates
type RatesStorage struct {
	sugar      *zap.SugaredLogger
	db         *sqlx.DB
	tableNames map[string]string
}

//TearDown removes all the tables
func (rdb *RatesStorage) TearDown() error {
	const dropFMT = `
	DROP VIEW rates_view;
	DROP TABLE %[1]s,%[2]s,%[3]s,%[4]s,%[5]s CASCADE;
	`
	tx, err := rdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, rdb.sugar, &err)
	query := fmt.Sprintf(dropFMT,
		rdb.tableNames[reserveTableName],
		rdb.tableNames[baseTableName],
		rdb.tableNames[tokenTableName],
		rdb.tableNames[rateTableName],
		rdb.tableNames[usdTableName],
	)
	_, err = tx.Exec(query)
	return err
}

//Close close DB connection
func (rdb *RatesStorage) Close() error {
	if rdb.db != nil {
		return rdb.db.Close()
	}
	return nil
}

//getTokenBase from pair return token and base from base-token pair
func getTokenBaseFromPair(pair string) (string, string, error) {
	ids := strings.Split(strings.TrimSpace(pair), "-")
	if len(ids) != 2 {
		return "", "", fmt.Errorf("pair %s is malformed. Expected base-token, got", pair)
	}
	return ids[1], ids[0], nil
}

func (rdb *RatesStorage) updateRsvAddrs(tx *sqlx.Tx, rsvs []string) error {
	const rsvStmt = `INSERT INTO %[1]s(address)
	VALUES(unnest($1::TEXT[]))
	ON CONFLICT ON CONSTRAINT %[1]s_address_key DO NOTHING`
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.updateRsvAddr",
	)

	query := fmt.Sprintf(rsvStmt, rdb.tableNames[reserveTableName])
	logger.Debugw("updating rsv...", "query", query)

	_, err := tx.Exec(query, pq.StringArray(rsvs))
	return err
}

func (rdb *RatesStorage) updateTokens(tx *sqlx.Tx, tokens []string) error {
	const tkStmt = `INSERT INTO %[1]s(symbol)
	VALUES(unnest($1::TEXT[]))
	ON CONFLICT ON CONSTRAINT %[1]s_symbol_key DO NOTHING`
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.updateToken",
	)

	query := fmt.Sprintf(tkStmt, rdb.tableNames[tokenTableName])
	logger.Debugw("updating tokens...", "query", query)

	_, err := tx.Exec(query, pq.StringArray(tokens))
	return err
}

func (rdb *RatesStorage) updateBases(tx *sqlx.Tx, bases []string) error {
	const bsStmt = `INSERT INTO %[1]s(symbol) 
	VALUES(unnest($1::TEXT[]))
	ON CONFLICT ON CONSTRAINT %[1]s_symbol_key DO NOTHING`
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.updateBases",
	)

	query := fmt.Sprintf(bsStmt, rdb.tableNames[baseTableName])
	logger.Debugw("updating bases...", "query", query)

	_, err := tx.Exec(query, pq.StringArray(bases))
	return err
}

//UpdateRatesRecords update mutiple rate records from a block with mutiple reserve address into the DB
func (rdb *RatesStorage) UpdateRatesRecords(blockInfo lastblockdaily.BlockInfo, rateRecords map[string]map[string]common.ReserveRateEntry) error {
	var (
		rsvAddrs    = make(map[string]bool)
		rsvAddrsArr []string
		tokens      = make(map[string]bool)
		tokensArr   []string
		bases       = make(map[string]bool)
		basesArr    []string
		logger      = rdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"block_number", blockInfo.Block,
			"timestamp", blockInfo.Timestamp.String(),
		)
		nRecord = 0
	)

	const rtStmt = `INSERT INTO %[1]s(time, token_id, base_id, block, buy_reserve_rate, sell_reserve_rate, reserve_id)
	VALUES ($1, 
		(SELECT id FROM %[2]s as tk WHERE tk.symbol= $2),
		(SELECT id FROM %[3]s as bs WHERE bs.symbol= $3),
		$4, 
		$5, 
		$6, 
		(SELECT id FROM %[4]s as rs WHERE rs.Address= $7)
	)
	ON CONFLICT ON CONSTRAINT %[1]s_no_duplicate DO NOTHING;`

	for rsvAddr, rateRecord := range rateRecords {
		if _, ok := rsvAddrs[rsvAddr]; !ok {
			rsvAddrs[rsvAddr] = true
			rsvAddrsArr = append(rsvAddrsArr, rsvAddr)
		}
		for pair := range rateRecord {
			token, base, err := getTokenBaseFromPair(pair)
			if err != nil {
				return err
			}
			if _, ok := tokens[token]; !ok {
				tokens[token] = true
				tokensArr = append(tokensArr, token)
			}

			if _, ok := bases[base]; !ok {
				bases[base] = true
				basesArr = append(basesArr, base)
			}
		}
	}
	tx, err := rdb.db.Beginx()
	defer pgsql.CommitOrRollback(tx, rdb.sugar, &err)

	if err != nil {
		return err
	}
	sort.Strings(rsvAddrsArr)
	if err := rdb.updateRsvAddrs(tx, rsvAddrsArr); err != nil {
		return err
	}

	sort.Strings(tokensArr)
	if err := rdb.updateTokens(tx, tokensArr); err != nil {
		return err
	}

	sort.Strings(basesArr)
	if err := rdb.updateBases(tx, basesArr); err != nil {
		return err
	}

	query := fmt.Sprintf(rtStmt, rdb.tableNames[rateTableName], rdb.tableNames[tokenTableName], rdb.tableNames[baseTableName], rdb.tableNames[reserveTableName])
	logger.Debugw("updating rates...", "query", query)
	for rsvAddr, rateRecord := range rateRecords {
		for pair, rate := range rateRecord {

			token, base, err := getTokenBaseFromPair(pair)
			if err != nil {
				return err
			}
			_, err = tx.Exec(query,
				blockInfo.Timestamp,
				token,
				base,
				blockInfo.Block,
				rate.BuyReserveRate,
				rate.SellReserveRate,
				rsvAddr,
			)
			if err != nil {
				return err
			}
			nRecord++
		}
	}
	logger.Debugw("updating rates finished", "total records", nRecord)

	return err
}

//GetRates return the reserve rates in a period of times
func (rdb *RatesStorage) GetRates(reserves []ethereum.Address, from, to time.Time) (map[string]storage.AccountingReserveRates, error) {
	var (
		result    = make(map[string]storage.AccountingReserveRates)
		rowData   = ratesFromDB{}
		rsvsAddrs []string
		logger    = rdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.GetRates",
			"from", from.String(),
			"to", to.String(),
			"reserves", reserves,
		)
	)
	const (
		selectStmt = `SELECT * FROM  rates_view
		WHERE  time>=$1 AND time<$2 AND reserve IN (SELECT unnest($3::text[]));`
	)
	for _, rsv := range reserves {
		rsvsAddrs = append(rsvsAddrs, rsv.Hex())
	}
	logger.Debugw("Querrying rate...", "query", selectStmt)

	rows, err := rdb.db.Queryx(selectStmt, from, to, pq.StringArray(rsvsAddrs))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.StructScan(&rowData); err != nil {
			return nil, err
		}
		if _, ok := result[rowData.Reserve]; !ok {
			result[rowData.Reserve] = map[time.Time]map[string]map[string]float64{
				rowData.Time: {
					rowData.Base: {
						rowData.Token: rowData.BuyRate,
					},
				},
			}
		}
		if _, ok := result[rowData.Reserve][rowData.Time]; !ok {
			result[rowData.Reserve][rowData.Time] = map[string]map[string]float64{
				rowData.Base: {
					rowData.Token: rowData.BuyRate,
				},
			}
		}
		if _, ok := result[rowData.Reserve][rowData.Time][rowData.Base]; !ok {
			result[rowData.Reserve][rowData.Time][rowData.Base] = map[string]float64{
				rowData.Token: rowData.BuyRate,
			}
		}
		result[rowData.Reserve][rowData.Time][rowData.Base][rowData.Token] = rowData.BuyRate
	}
	return result, nil
}

//UpdateETHUSDPrice store the ETHUSD rate at that blockInfo
func (rdb *RatesStorage) UpdateETHUSDPrice(blockInfo lastblockdaily.BlockInfo, ethusdRate float64) error {
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
		"block_number", blockInfo.Block,
		"timestamp", blockInfo.Timestamp.String(),
	)
	const updateStmt = `INSERT INTO %[1]s(time, block, rate)
	VALUES ( 
		$1,
		$2, 
		$3)
	ON CONFLICT (time,block) DO UPDATE SET rate=EXCLUDED.rate;`
	query := fmt.Sprintf(updateStmt,
		rdb.tableNames[usdTableName],
	)

	logger.Debugw("updating eth-usdrates...", "query", query)
	tx, err := rdb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, rdb.sugar, &err)
	_, err = tx.Exec(query,
		blockInfo.Timestamp,
		blockInfo.Block,
		ethusdRate,
	)
	if err != nil {
		return err
	}

	return err
}

//GetLastResolvedBlockInfo return block info of the rate with latest timestamp
func (rdb *RatesStorage) GetLastResolvedBlockInfo() (lastblockdaily.BlockInfo, error) {
	const (
		selectStmt = `SELECT time,block FROM %[1]s WHERE time=
		(SELECT MAX(time) FROM %[1]s) LIMIT 1`
	)
	var (
		usdTableResult  = lastblockdaily.BlockInfo{}
		rateTableResult = lastblockdaily.BlockInfo{}
		logger          = rdb.sugar.With("func", "accounting/reserve-rate/storage/postgres/accounting_rates_postgres.GetLastResolvedBlockInfo")
	)

	query := fmt.Sprintf(selectStmt, rdb.tableNames[rateTableName])
	logger.Debugw("Querrying last resolved block from rates table...", "query", query)
	if err := rdb.db.Get(&rateTableResult, query); err != nil {
		return rateTableResult, err
	}
	query = fmt.Sprintf(selectStmt, rdb.tableNames[usdTableName])
	logger.Debugw("Querrying last resolved block from usd table...", "query", query)
	if err := rdb.db.Get(&usdTableResult, query); err != nil {
		return usdTableResult, err
	}
	if usdTableResult.Timestamp.Before(rateTableResult.Timestamp) {
		return usdTableResult, nil
	}
	return rateTableResult, nil
}
