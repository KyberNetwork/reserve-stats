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
	lbdCommon "github.com/KyberNetwork/reserve-stats/lib/lastblockdaily/common"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
)

const (
	reserveTableName = "reserves"
	tokenTableName   = "tokens"
	quoteTableName   = "quotes"
	rateTableName    = "token_rates"
	usdTableName     = "usd_rates"
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
		rdb.tableNames[quoteTableName],
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

//getTokenQuote from pair return token and quote from quote-token pair
func getTokenQuoteFromPair(pair string) (string, string, error) {
	ids := strings.Split(strings.TrimSpace(pair), "-")
	if len(ids) != 2 {
		return "", "", fmt.Errorf("pair %s is malformed. Expected quote-token, got", pair)
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

func (rdb *RatesStorage) updateQuotes(tx *sqlx.Tx, quotes []string) error {
	const bsStmt = `INSERT INTO %[1]s(symbol) 
	VALUES(unnest($1::TEXT[]))
	ON CONFLICT ON CONSTRAINT %[1]s_symbol_key DO NOTHING`
	var logger = rdb.sugar.With(
		"func", "reserverates/storage/postgres/RateStorage.updateQuotes",
	)

	query := fmt.Sprintf(bsStmt, rdb.tableNames[quoteTableName])
	logger.Debugw("updating quotes...", "query", query)

	_, err := tx.Exec(query, pq.StringArray(quotes))
	return err
}

//UpdateRatesRecords update mutiple rate records from a block with mutiple reserve address into the DB
func (rdb *RatesStorage) UpdateRatesRecords(blockInfo lbdCommon.BlockInfo, rateRecords map[string]map[string]float64, ethusdRate float64) error {
	var (
		rsvAddrs    = make(map[string]bool)
		rsvAddrsArr []string
		tokens      = make(map[string]bool)
		tokensArr   []string
		quotes      = make(map[string]bool)
		quotesArr   []string
		logger      = rdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.UpdateRatesRecords",
			"block_number", blockInfo.Block,
			"timestamp", blockInfo.Timestamp.String(),
		)
		nRecord = 0
	)

	const rtStmt = `INSERT INTO %[1]s(time, token_id, quote_id, block, rate, reserve_id)
	VALUES ($1, 
		(SELECT id FROM %[2]s as tk WHERE tk.symbol= $2),
		(SELECT id FROM %[3]s as bs WHERE bs.symbol= $3),
		$4, 
		$5, 
		(SELECT id FROM %[4]s as rs WHERE rs.Address= $6)
	)
	ON CONFLICT ON CONSTRAINT %[1]s_no_duplicate DO NOTHING;`

	for rsvAddr, rateRecord := range rateRecords {
		if _, ok := rsvAddrs[rsvAddr]; !ok {
			rsvAddrs[rsvAddr] = true
			rsvAddrsArr = append(rsvAddrsArr, rsvAddr)
		}
		for pair := range rateRecord {
			token, quote, err := getTokenQuoteFromPair(pair)
			if err != nil {
				return err
			}
			if _, ok := tokens[token]; !ok {
				tokens[token] = true
				tokensArr = append(tokensArr, token)
			}

			if _, ok := quotes[quote]; !ok {
				quotes[quote] = true
				quotesArr = append(quotesArr, quote)
			}
		}
	}
	tx, err := rdb.db.Beginx()
	defer pgsql.CommitOrRollback(tx, rdb.sugar, &err)

	if err != nil {
		return err
	}

	if err = rdb.updateETHUSDPrice(blockInfo, ethusdRate, tx); err != nil {
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

	sort.Strings(quotesArr)
	if err := rdb.updateQuotes(tx, quotesArr); err != nil {
		return err
	}

	query := fmt.Sprintf(rtStmt, rdb.tableNames[rateTableName], rdb.tableNames[tokenTableName], rdb.tableNames[quoteTableName], rdb.tableNames[reserveTableName])
	logger.Debugw("updating rates...", "query", query)
	for rsvAddr, rateRecord := range rateRecords {
		for pair, rate := range rateRecord {
			token, quote, err := getTokenQuoteFromPair(pair)
			if err != nil {
				return err
			}
			_, err = tx.Exec(query,
				blockInfo.Timestamp.UTC(),
				token,
				quote,
				blockInfo.Block,
				rate,
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

//GetETHUSDRates  return the ETH/USD  rate in  a period of times
func (rdb *RatesStorage) GetETHUSDRates(from, to time.Time) (storage.AccountingReserveRates, error) {
	var (
		result = make(storage.AccountingReserveRates)

		dbResult []usdRateFromDB
		logger   = rdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.GetUSDRate",
			"from", from.String(),
			"to", to.String(),
		)
	)
	const (
		selectStmt = `SELECT time,rate FROM %[1]s WHERE time>=$1 AND time <$2`
	)
	query := fmt.Sprintf(selectStmt, rdb.tableNames[usdTableName])
	logger.Debug("Queryingrate...", "query", query)
	if err := rdb.db.Select(&dbResult, query, from, to); err != nil {
		return result, err
	}
	for _, record := range dbResult {
		result[record.Time.UTC()] = map[string]map[string]float64{
			"USD": {
				"ETH": record.Rate,
			},
		}
	}
	return result, nil
}

//GetRates return the reserve rates in a period of times
func (rdb *RatesStorage) GetRates(from, to time.Time) (map[string]storage.AccountingReserveRates, error) {
	var (
		result  = make(map[string]storage.AccountingReserveRates)
		rowData = ratesFromDB{}
		logger  = rdb.sugar.With(
			"func", "reserverates/storage/postgres/RateStorage.GetRates",
			"from", from.String(),
			"to", to.String(),
		)
	)
	const (
		selectStmt = `SELECT * FROM  rates_view
		WHERE  time>=$1 AND time<$2`
	)
	logger.Debugw("Querying rate...", "query", selectStmt)

	rows, err := rdb.db.Queryx(selectStmt, from.UTC(), to.UTC())
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.StructScan(&rowData); err != nil {
			return nil, err
		}

		timestamp := rowData.Time.UTC()

		if _, ok := result[rowData.Reserve]; !ok {
			result[rowData.Reserve] = map[time.Time]map[string]map[string]float64{
				timestamp: {
					rowData.Quote: {
						rowData.Token: rowData.Rate,
					},
				},
			}
		}
		if _, ok := result[rowData.Reserve][timestamp]; !ok {
			result[rowData.Reserve][timestamp] = map[string]map[string]float64{
				rowData.Quote: {
					rowData.Token: rowData.Rate,
				},
			}
		}
		if _, ok := result[rowData.Reserve][timestamp][rowData.Quote]; !ok {
			result[rowData.Reserve][timestamp][rowData.Quote] = map[string]float64{
				rowData.Token: rowData.Rate,
			}
		}
		result[rowData.Reserve][timestamp][rowData.Quote][rowData.Token] = rowData.Rate
	}
	return result, nil
}

//updateETHUSDPrice store the ETHUSD rate at that blockInfo
func (rdb *RatesStorage) updateETHUSDPrice(blockInfo lbdCommon.BlockInfo, ethusdRate float64, tx *sqlx.Tx) error {
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
	_, err := tx.Exec(query,
		blockInfo.Timestamp.UTC(),
		blockInfo.Block,
		ethusdRate,
	)

	return err
}

//GetLastResolvedBlockInfo return block info of the rate with latest timestamp
func (rdb *RatesStorage) GetLastResolvedBlockInfo(reserveAddr ethereum.Address) (lbdCommon.BlockInfo, error) {
	const (
		selectStmt = `SELECT time,block FROM %[1]s WHERE time=
		(SELECT MAX(time) FROM %[1]s) AND reserve_id=(SELECT id FROM %[2]s WHERE %[2]s.address=$1) LIMIT 1`
	)
	var (
		rateTableResult = lbdCommon.BlockInfo{}
		logger          = rdb.sugar.With("func", "accounting/reserve-rate/storage/postgres/RatesStorage.GetLastResolvedBlockInfo")
	)

	query := fmt.Sprintf(selectStmt, rdb.tableNames[rateTableName], rdb.tableNames[reserveTableName])
	logger.Debugw("Querying last resolved block from rates table...", "query", query)
	if err := rdb.db.Get(&rateTableResult, query, reserveAddr.Hex()); err != nil {
		return rateTableResult, err
	}
	rateTableResult.Timestamp = rateTableResult.Timestamp.UTC()

	return rateTableResult, nil
}
