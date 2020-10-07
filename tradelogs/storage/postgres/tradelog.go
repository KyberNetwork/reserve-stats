package postgres

import (
	"encoding/json"
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

func (tldb *TradeLogDB) prepareFeeRecords(r *record) ([]string, []string, []float64, []float64, []float64, []float64, []float64, []uint, [][]byte, [][]byte, error) {
	var (
		reserveAddresses, platformWallets                 []string
		walletFees, burns, rebates, rewards, platformFees []float64
		rebateWallets, rebatePercents                     [][]byte
		indexes                                           []uint
	)
	for _, f := range r.Fee {
		reserveAddresses = append(reserveAddresses, f.ReserveAddr.Hex())
		platformWallets = append(platformWallets, f.PlatformWallet.Hex())
		walletFee, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.WalletFee)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, err
		}
		walletFees = append(walletFees, walletFee)
		platformFee, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.PlatformFee)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, err
		}
		platformFees = append(platformFees, platformFee)
		burn, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Burn)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, err
		}
		burns = append(burns, burn)
		rebate, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Rebate)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, err
		}
		rebates = append(rebates, rebate)
		reward, err := tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Reward)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, err
		}
		rewards = append(rewards, reward)
		indexes = append(indexes, f.Index)
		var (
			wallets  []string
			percents []uint64
		)
		for _, wallet := range f.RebateWallets {
			wallets = append(wallets, wallet.Hex())
		}
		walletsJSON, err := json.Marshal(wallets)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, err
		}
		rebateWallets = append(rebateWallets, walletsJSON)
		for _, percent := range f.RebatePercentBpsPerWallet {
			percents = append(percents, percent.Uint64())
		}
		percentsJSON, err := json.Marshal(percents)
		if err != nil {
			return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, err
		}
		rebatePercents = append(rebatePercents, percentsJSON)
	}
	return reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, indexes, rebateWallets, rebatePercents, nil
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

func (tldb *TradeLogDB) prepareSplitRecords(r *record) ([]string, []string, []string, []float64, []float64, []float64, []uint, error) {
	var (
		reserveAddressIDs, srcAddresses, destAddresses []string
		srcAmounts, rates, dstAmounts                  []float64
		indexes                                        []uint
	)
	// before katalyst
	if r.Version != 4 {
		var index uint
		if ethereum.HexToAddress(r.SrcAddress) != blockchain.ETHAddr {
			reserveAddressIDs = append(reserveAddressIDs, r.SrcReserveAddress)
			srcAddresses = append(srcAddresses, r.SrcAddress)
			destAddresses = append(destAddresses, blockchain.ETHAddr.Hex())
			srcAmounts = append(srcAmounts, r.SrcAmount)
			dstAmounts = append(dstAmounts, r.OriginalEthAmount)
			rates = append(rates, 0) // TODO: need to calculate rate here
			indexes = append(indexes, index+1)
			index++
		}

		if ethereum.HexToAddress(r.DestAddress) != blockchain.ETHAddr {
			reserveAddressIDs = append(reserveAddressIDs, r.DstReserveAddress)
			srcAddresses = append(srcAddresses, blockchain.ETHAddr.Hex())
			destAddresses = append(destAddresses, r.DestAddress)
			srcAmounts = append(srcAmounts, r.OriginalEthAmount)
			dstAmounts = append(dstAmounts, r.DestAmount)
			rates = append(rates, 0) // TODO: need to calculate rate here
			indexes = append(indexes, index+1)
		}
		return reserveAddressIDs, srcAddresses, destAddresses, srcAmounts, rates, dstAmounts, indexes, nil
	}

	// After katalyst
	var splitIndex uint
	for index, s := range r.T2EReserves {
		reserveAddressIDs = append(reserveAddressIDs, ethereum.BytesToHash(s[:]).Hex())
		srcAddresses = append(srcAddresses, r.SrcAddress)
		destAddresses = append(destAddresses, blockchain.ETHAddr.Hex())
		srcAmounts = append(srcAmounts, r.T2ESrcAmount[index])
		dstAmount, err := tldb.calculateDstAmount(r.SrcAddress, blockchain.ETHAddr.Hex(), r.T2ESrcAmount[index], r.T2ERates[index])
		if err != nil {
			return reserveAddressIDs, srcAddresses, destAddresses, srcAmounts, rates, dstAmounts, indexes, nil
		}
		dstAmounts = append(dstAmounts, dstAmount)
		rates = append(rates, r.T2ERates[index])
		indexes = append(indexes, splitIndex+1)
		splitIndex++
	}

	for index, s := range r.E2TReserves {
		reserveAddressIDs = append(reserveAddressIDs, ethereum.BytesToHash(s[:]).Hex())
		srcAddresses = append(srcAddresses, blockchain.ETHAddr.Hex())
		destAddresses = append(destAddresses, r.DestAddress)
		srcAmounts = append(srcAmounts, r.E2TSrcAmount[index])
		dstAmount, err := tldb.calculateDstAmount(blockchain.ETHAddr.Hex(), r.DestAddress, r.E2TSrcAmount[index], r.E2TRates[index])
		if err != nil {
			return reserveAddressIDs, srcAddresses, destAddresses, srcAmounts, rates, dstAmounts, indexes, nil
		}
		dstAmounts = append(dstAmounts, dstAmount)
		rates = append(rates, r.E2TRates[index])
		indexes = append(indexes, splitIndex+1)
		splitIndex++
	}

	return reserveAddressIDs, srcAddresses, destAddresses, srcAmounts, rates, dstAmounts, indexes, nil
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
				$13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25,
				$26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43
			);`
			var tradelogID uint64
			reserveAddresses, platformWallets, burns, rebates, rewards, platformFees, walletFees, feeIndexes, rebateWallets, rebatePercents, err := tldb.prepareFeeRecords(r)
			if err != nil {
				logger.Debugw("failed to prepare fee records", "error", err)
				return err
			}
			reserveAddressIds, srcAddresses, dstAddresses, srcAmounts, rates, dstAmounts, splitIndexes, err := tldb.prepareSplitRecords(r)
			if err != nil {
				logger.Debugw("failed to prepared split records", "error", err)
				return err
			}

			logger.Debugw("reserve", "ids", reserveAddressIds)

			if err != nil {
				logger.Debugw("failed to prepare fees record", "error", err)
				return err
			}
			if err := tx.Get(&tradelogID, query, 0,
				r.Timestamp, r.BlockNumber, r.TransactionHash,
				r.EthAmount, r.OriginalEthAmount, r.UserAddress, r.SrcAddress, r.DestAddress,
				r.WalletAddress,
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
				pq.Array(walletFees),
				pq.Array(platformFees),
				pq.Array(burns),
				pq.Array(rebates),
				pq.Array(rewards),
				pq.Array(feeIndexes),
				pq.Array(rebateWallets),
				pq.Array(rebatePercents),
				pq.StringArray(reserveAddressIds),
				pq.StringArray(srcAddresses),
				pq.StringArray(dstAddresses),
				pq.Array(srcAmounts),
				pq.Array(rates),
				pq.Array(dstAmounts),
				pq.Array(splitIndexes),
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
