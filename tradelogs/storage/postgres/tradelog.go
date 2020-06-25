package postgres

import (
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// SaveTradelogV4 save tradelog for v4
func (tldb *TradeLogDB) SaveTradelogV4() error {
	return nil
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
