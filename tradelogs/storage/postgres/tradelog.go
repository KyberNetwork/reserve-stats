package postgres

import (
	"math/big"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/pgsql"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/KyberNetwork/reserve-stats/tradelogs/storage/postgres/schema"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
)

func (tldb *TradeLogDB) saveReserve(reserves []common.Reserve) error {
	var (
		logger                               = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		addresses, reserveIDs, rebateWallets []string
		blockNumbers, reserveTypes           []uint64
	)
	query := `INSERT INTO reserve (address, reserve_id, rebate_wallet, block_number, reserve_type)
	VALUES(
		UNNEST($1::TEXT[]),
		UNNEST($2::TEXT[]),
		UNNEST($3::TEXT[]),
		UNNEST($4::INTEGER[]),
		UNNEST($5::INTEGER[])
	) ON CONFLICT (address, reserve_id, block_number) DO NOTHING;`
	logger.Infow("save reserve", "query", query)
	for _, r := range reserves {
		addresses = append(addresses, r.Address.Hex())
		reserveIDs = append(reserveIDs, ethereum.BytesToHash(r.ReserveID[:]).Hex())
		rebateWallets = append(rebateWallets, r.RebateWallet.Hex())
		blockNumbers = append(blockNumbers, r.BlockNumber)
		reserveTypes = append(reserveTypes, r.ReserveType)
	}
	if _, err := tldb.db.Exec(query, pq.StringArray(addresses), pq.StringArray(reserveIDs), pq.StringArray(rebateWallets),
		pq.Array(blockNumbers), pq.Array(reserveTypes)); err != nil {
		logger.Errorw("failed to add reserve into db", "error", err)
		return err
	}
	return nil
}

func (tldb *TradeLogDB) updateRebateWallet(reserves []common.Reserve) error {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	)
	query := `INSERT INTO reserve(address, reserve_id, rebate_wallet, block_number, reserve_type)
		VALUES (
		(SELECT address FROM reserve WHERE reserve_id = $1 order by block_number desc limit 1),
		$1,
		$2,
		$3, 
		(SELECT reserve_type FROM reserve WHERE reserve_id = $1 order by block_number desc limit 1)
		) ON CONFLICT (address, reserve_id, block_number) DO NOTHING;`
	logger.Infow("query update rebate wallet", "value", query)
	tx, err := tldb.db.Beginx()
	if err != nil {
		return err
	}
	defer pgsql.CommitOrRollback(tx, logger, &err)
	for _, r := range reserves {
		if _, err := tx.Exec(query, ethereum.BytesToHash(r.ReserveID[:]).Hex(), r.RebateWallet.Hex(), r.BlockNumber); err != nil {
			return err
		}
	}
	return nil
}

func (tldb *TradeLogDB) calculateDstAmount(srcAddress, dstAddress string, srcAmount, rate float64) (float64, error) {
	var (
		srcDecimals, dstDecimals int64
		err                      error
		dstAmount                float64
	)
	srcDecimals, err = tldb.tokenAmountFormatter.GetDecimals(ethereum.HexToAddress(srcAddress))
	if err != nil {
		return dstAmount, err
	}
	dstDecimals, err = tldb.tokenAmountFormatter.GetDecimals(ethereum.HexToAddress(dstAddress))
	if err != nil {
		return dstAmount, err
	}
	srcAmountBig, err := tldb.tokenAmountFormatter.ToWei(ethereum.HexToAddress(srcAddress), srcAmount)
	if err != nil {
		return dstAmount, err
	}

	rateBig, err := tldb.tokenAmountFormatter.ToWei(blockchain.ETHAddr, rate)
	if err != nil {
		return dstAmount, err
	}
	dstAmountTmp := big.NewInt(0).Mul(srcAmountBig, rateBig)
	// this formula is base on https://github.com/KyberNetwork/smart-contracts/blob/Katalyst/contracts/sol6/utils/Utils5.sol#L88
	if dstDecimals >= srcDecimals {
		precision := new(big.Float).SetInt(new(big.Int).Exp(
			big.NewInt(10), big.NewInt(18), nil,
		))
		exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(dstDecimals-srcDecimals), nil)
		tmp := big.NewInt(0).Mul(dstAmountTmp, exp)
		dstAmountInt, _ := new(big.Float).Quo(new(big.Float).SetInt(tmp), precision).Int(nil)
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(ethereum.HexToAddress(dstAddress), dstAmountInt)
		if err != nil {
			return dstAmount, err
		}
	} else {
		precision := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil)
		exp := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(srcDecimals-dstDecimals), nil)
		tmp := big.NewInt(0).Mul(exp, precision)
		dstAmountInt, _ := new(big.Float).Quo(new(big.Float).SetInt(dstAmountTmp), new(big.Float).SetInt(tmp)).Int(nil)
		dstAmount, err = tldb.tokenAmountFormatter.FromWei(ethereum.HexToAddress(dstAddress), dstAmountInt)
		if err != nil {
			return dstAmount, err
		}
	}
	return dstAmount, nil
}

// SaveTradeLogs persist trade logs to DB
func (tldb *TradeLogDB) SaveTradeLogs(crResult *common.CrawlResult) (err error) {
	var (
		logger              = tldb.sugar.With("func", caller.GetCurrentFunctionName())
		reserveAddress      = make(map[string]struct{})
		reserveAddressArray []string
		tokens              = make(map[string]struct{})
		tokensArray         []string
		records             []*record

		users = make(map[ethereum.Address]struct{})
	)
	if crResult != nil {
		if len(crResult.Reserves) > 0 {
			if err := tldb.saveReserve(crResult.Reserves); err != nil {
				return err
			}
		}

		if len(crResult.UpdateWallets) > 0 {
			if err := tldb.updateRebateWallet(crResult.UpdateWallets); err != nil {
				return err
			}
		}

		logs := crResult.Trades
		for _, log := range logs {
			r, err := tldb.recordFromTradeLog(log)
			if err != nil {
				return err
			}

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
			reserve := r.SrcReserveAddress
			if reserve != "" {
				if _, ok := reserveAddress[reserve]; !ok {
					reserveAddress[reserve] = struct{}{}
					reserveAddressArray = append(reserveAddressArray, reserve)
				}
			}
			reserve = r.DstReserveAddress
			if reserve != "" {
				if _, ok := reserveAddress[reserve]; !ok {
					reserveAddress[reserve] = struct{}{}
					reserveAddressArray = append(reserveAddressArray, reserve)
				}
			}
		}

		tx, err := tldb.db.Beginx()
		if err != nil {
			return err
		}
		defer pgsql.CommitOrRollback(tx, logger, &err)
		if len(reserveAddressArray) > 0 {
			err = tldb.saveReserveAddress(tx, reserveAddressArray)
			if err != nil {
				logger.Debugw("failed to save reserve address", "error", err)
				return err
			}
		}

		err = tldb.saveTokens(tx, tokensArray)
		if err != nil {
			logger.Debugw("failed to save token", "error", err)
			return err
		}

		for _, r := range records {
			logger.Debugw("Record", "record", r)
			_, err = tx.NamedExec(insertionUserTemplate, r)
			if err != nil {
				logger.Infow("user", "address", r.UserAddress)
				logger.Debugw("Error while add users", "error", err)
				return err
			}

			_, err = tx.NamedExec(insertionWalletTemplate, r)
			if err != nil {
				logger.Debugw("Error while add wallet", "error", err)
				return err
			}

			query := `SELECT _id as id FROM 
			create_or_update_tradelogs(
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
				$13, $14, $15, $16, $17, $18, $19, $20, $21
			);`
			var tradelogID uint64
			if err != nil {
				logger.Debugw("failed to prepare fees record", "error", err)
				return err
			}
			if err := tx.Get(&tradelogID, query, 0,
				r.Timestamp, r.BlockNumber, r.TransactionHash,
				r.EthAmount, r.OriginalEthAmount, r.UserAddress, r.SrcAddress, r.DestAddress,
				r.SrcAmount, r.DestAmount,
				r.ETHUSDRate,
				r.ETHUSDProvider,
				r.Index,
				r.IsFirstTrade,
				r.TxSender,
				r.ReceiverAddress,
				r.GasUsed,
				r.GasPrice,
				r.TransactionFee,
				r.Version,
			); err != nil {
				logger.Debugw("failed to save tradelogs", "error", err)
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
