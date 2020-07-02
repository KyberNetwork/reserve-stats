package postgres

import (
	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	"github.com/jmoiron/sqlx"
)

// GetFeeByTradelogID return fee by tradelog id
func (tldb *TradeLogDB) GetFeeByTradelogID(tradelogID uint64) (common.TradelogFee, error) {
	var (
		fee common.TradelogFee
	)
	return fee, nil
}

// SaveFee save fee by tradelog
func (tldb *TradeLogDB) SaveFee(tx *sqlx.Tx, fees []common.TradelogFee, tradelogID int64) error {
	var (
		logger = tldb.sugar.With("func", caller.GetCurrentFunctionName())
	)
	logger.Info("save fee")
	query := `INSERT INTO fee(trade_id, reserve_address, wallet_address, wallet_fee, platform_fee, burn, rebate, reward)
	VALUES(
		$1,
		$2,
		$3, 
		$4, 
		$5, 
		$6,
		$7,
		$8
	);`
	for _, f := range fees {
		var (
			walletFee, platformFee, burn, rebate, reward float64
			err                                          error
		)
		if f.WalletFee != nil {
			walletFee, err = tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.WalletFee)
			if err != nil {
				logger.Debugw("failed to parse wallet fee", "error", err)
				return err
			}
		}
		if f.PlatformFee != nil {
			platformFee, err = tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.PlatformFee)
			if err != nil {
				logger.Debugw("failed to parse wallet fee", "error", err)
				return err
			}
		}
		if f.Burn != nil {
			burn, err = tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Burn)
			if err != nil {
				logger.Debugw("failed to parse wallet fee", "error", err)
				return err
			}
		}
		if f.Rebate != nil {
			rebate, err = tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Rebate)
			if err != nil {
				logger.Debugw("failed to parse rebate", "error", err)
				return err
			}
		}
		if f.Reward != nil {
			reward, err = tldb.tokenAmountFormatter.FromWei(blockchain.ETHAddr, f.Reward)
			if err != nil {
				logger.Debugw("failed to parse reward", "error", err)
				return err
			}
		}
		logger.Infow("Query params",
			"tradelogID", tradelogID,
			"reserve address", f.ReserveAddr.Hex(),
			"wallet address", f.PlatformWallet.Hex(),
			"wallet fee", walletFee,
			"platform fee", platformFee,
			"burn", burn,
			"rebate", rebate,
			"reward", reward,
		)
		if _, err := tx.Exec(query, tradelogID, f.ReserveAddr.Hex(), f.PlatformWallet.Hex(), walletFee, platformFee, burn, rebate, reward); err != nil {
			return err
		}
	}
	return nil
}
