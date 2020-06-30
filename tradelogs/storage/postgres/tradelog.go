package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
)

// SaveTradelogV4 save tradelog for v4
func (tldb *TradeLogDB) SaveTradelogV4() error {
	return nil
}

func (tldb *TradeLogDB) saveReserve(reserves []common.Reserve) error {
	var (
		logger                               = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		addresses, reserveIDs, rebateWallets []string
		blockNumbers                         []uint64
	)
	query := `INSERT INTO reserve (address, reserve_id, rebate_wallet, block_number)
	VALUES(
		UNNEST($1::TEXT[]),
		UNNEST($2::TEXT[]),
		UNNEST($3::TEXT[]),
		UNNEST($4::INTEGER[])
	) ON CONFLICT (address, reserve_id, block_number) DO NOTHING;`
	logger.Infow("query", "value", query)
	for _, r := range reserves {
		addresses = append(addresses, r.Address.Hex())
		reserveIDs = append(reserveIDs, ethereum.BytesToHash(r.ReserveID[:]).Hex())
		rebateWallets = append(rebateWallets, r.RebateWallet.Hex())
		blockNumbers = append(blockNumbers, r.BlockNumber)
	}
	if _, err := tldb.db.Exec(query, pq.StringArray(addresses), pq.StringArray(reserveIDs), pq.StringArray(rebateWallets), pq.Array(blockNumbers)); err != nil {
		logger.Errorw("failed to add reserve into db", "error", err)
		return err
	}
	return nil
}

// SaveTradeLogs persist trade logs to DB
func (tldb *TradeLogDB) SaveTradeLogs(crResult *common.CrawlResult) (err error) {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		// reserveAddress      = make(map[string]struct{})
		reserveAddressArray []string
		tokens              = make(map[string]struct{})
		tokensArray         []string
		records             []*record

		users = make(map[ethereum.Address]struct{})
	)
	if crResult != nil {
		if len(crResult.Reserves) > 0 {
			tldb.saveReserve(crResult.Reserves)
		}

		logs := crResult.Trades
		for _, log := range logs {
			r, err := tldb.recordFromTradeLog(log)
			if err != nil {
				return err
			}

			byteRecords, _ := json.Marshal(r)
			fmt.Printf("record: %s\n", byteRecords)

			if _, ok := users[log.User.UserAddress]; ok {
				r.IsFirstTrade = false
			} else {
				isFirstTrade, err := tldb.isFirstTrade(log.User.UserAddress)
				if err != nil {
					return err
				}
				r.IsFirstTrade = isFirstTrade
			}
			records = append(records, r)
			users[log.User.UserAddress] = struct{}{}
		}

		for _, r := range records {
			token := r.SrcAddress
			if _, ok := tokens[token]; !ok {
				tokens[token] = struct{}{}
				tokensArray = append(tokensArray, token)
			}
			token = r.DestAddress
			if _, ok := tokens[token]; !ok {
				tokens[token] = struct{}{}
				tokensArray = append(tokensArray, token)
			}
		}

		// 	for _, reserve := range r.E2TReserves {
		// 		if _, ok := reserveAddress[reserve]; !ok {
		// 			reserveAddress[reserve] = struct{}{}
		// 			reserveAddressArray = append(reserveAddressArray, reserve)
		// 		}
		// 		token := r.DestAddress
		// 		if _, ok := tokens[reserve]; !ok {
		// 			tokens[token] = struct{}{}
		// 			tokensArray = append(tokensArray, token)
		// 		}
		// 	}
		// }

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
	return nil
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