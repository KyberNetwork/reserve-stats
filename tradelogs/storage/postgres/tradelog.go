package postgres

import (
	"encoding/json"
	"fmt"

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
	logger.Infow("query", "value", query)
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
		(SELECT address FROM reserve WHERE reserve_id = $1),
		$1,
		$2,
		$3, 
		(SELECT reserve_type FROM reserve WHERE reserve_id = $1)
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

func (tldb *TradeLogDB) prepareFeeRecords(r *record) ([]string, []string, []float64, []float64, []float64, []float64, error) {
	var (
		reserveAddresses, platformWallets     []string
		burns, rebates, rewards, platformFees []float64
	)
	for _, f := range r.Fee {
		reserveAddresses = append(reserveAddresses, f.ReserveAddr.Hex())
		platformWallets = append(platformWallets, f.PlatformWallet.Hex())
		platformFee, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Burn)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, err
		}
		platformFees = append(platformFees, platformFee)
		burn, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Burn)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, err
		}
		burns = append(burns, burn)
		rebate, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Rebate)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, err
		}
		rebates = append(rebates, rebate)
		reward, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Reward)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, err
		}
		rewards = append(rewards, reward)
	}
	return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, nil
}

// func (tldb *TradeLogDB) prepareSplitRecords(r *record) ([]string, []string, []string, []float64, []float64, error) {
// 	var (
// 		reserveAddressIDs, srcAddresses, destAddresses []string
// 		srcAmounts, rates                              []float64
// 	)
// 	for index, s := range r.T2EReserves {
// 		reserveAddressIDs = append(reserveAddressIDs, ethereum.Bytes2Hex(s[:]))
// 		srcAddresses = append(srcAddresses, r.SrcAddress)
// 		destAddresses = append(destAddresses, blockchain.ETHAddr.Hex())
// 		srcAmounts = append(srcAmounts, r.T2ESrcAmount[index])
// 		rates = append(rates, r.T2ERates[index])
// 	}
// 	return reserveAddressIDs, srcAddresses, destAddresses, srcAmounts, rates, nil
// }

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
				$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25,
				$26, $27, $28, $29, $30, $31
			);`
			var tradelogID uint64
			reserveAddresses, platformWallets, platformFees, burns, rebates, rewards, err := tldb.prepareFeeRecords(r)
			// reserveAddressIds,

			if err != nil {
				logger.Debugw("failed to prepare fees record", "error", err)
				return err
			}
			if err := tx.Get(&tradelogID, query, 0,
				r.Timestamp, r.BlockNumber, r.TransactionHash,
				r.EthAmount, r.OriginalEthAmount, r.UserAddress, r.SrcAddress, r.DestAddress,
				r.SrcAmount, r.DestAmount,
				r.IntegrationApp,
				r.IP,
				r.Country,
				r.ETHUSDRate,
				r.ETHUSDProvider,
				r.Index,
				r.Kyced,
				r.IsFirstTrade,
				r.TxSender,
				r.ReceiverAddress,
				r.GasUsed,
				r.GasPrice,
				r.TransactionFee,
				r.Version,
				pq.StringArray(reserveAddresses),
				pq.StringArray(platformWallets),
				pq.Array(platformFees),
				pq.Array(burns),
				pq.Array(rebates),
				pq.Array(rewards),
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
