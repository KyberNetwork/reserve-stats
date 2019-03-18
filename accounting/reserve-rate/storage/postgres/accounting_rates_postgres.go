package postgres

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/reserve-rate/storage"
	"github.com/KyberNetwork/reserve-stats/lib/lastblockdaily"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
)

const (
	reserveTableName = "reserves"
	tokenTableName   = "tokens"
	// TODO: change this to quotes.
	baseTableName = "bases"
	rateTableName = "token_rates"
	usdTableName  = "usd_rates"
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
// TODO: parsing the pair string from storage doesn't look right,
//  as all the database implementations will need to redo it.
//  Maybe the fetcher should do it before passing to storage?
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
func (rdb *RatesStorage) UpdateRatesRecords(blockInfo lastblockdaily.BlockInfo, rateRecords map[string]map[string]float64) error {
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

	const rtStmt = `INSERT INTO %[1]s(time, token_id, base_id, block, rate, reserve_id)
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
				blockInfo.Timestamp.UTC(),
				token,
				base,
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
					rowData.Base: {
						rowData.Token: rowData.Rate,
					},
				},
			}
		}
		if _, ok := result[rowData.Reserve][timestamp]; !ok {
			result[rowData.Reserve][timestamp] = map[string]map[string]float64{
				rowData.Base: {
					rowData.Token: rowData.Rate,
				},
			}
		}
		if _, ok := result[rowData.Reserve][timestamp][rowData.Base]; !ok {
			result[rowData.Reserve][timestamp][rowData.Base] = map[string]float64{
				rowData.Token: rowData.Rate,
			}
		}
		result[rowData.Reserve][timestamp][rowData.Base][rowData.Token] = rowData.Rate
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
		blockInfo.Timestamp.UTC(),
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
	logger.Debugw("Querying last resolved block from rates table...", "query", query)
	if err := rdb.db.Get(&rateTableResult, query); err != nil {
		return rateTableResult, err
	}
	rateTableResult.Timestamp = rateTableResult.Timestamp.UTC()

	query = fmt.Sprintf(selectStmt, rdb.tableNames[usdTableName])
	logger.Debugw("Querying last resolved block from usd table...", "query", query)
	if err := rdb.db.Get(&usdTableResult, query); err != nil {
		return usdTableResult, err
	}
	// TODO: this won't work when we add a new reserve address and need to fetch all rates from its creation date.
	//  In other words, we need to track last inserted data of each reserve.
	usdTableResult.Timestamp = usdTableResult.Timestamp.UTC()

	if usdTableResult.Timestamp.Before(rateTableResult.Timestamp) {
		return usdTableResult, nil
	}
	return rateTableResult, nil
}
