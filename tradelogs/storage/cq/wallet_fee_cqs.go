package cq

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	walletschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletfee"
	walletFeeVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletfee_volume"
)

// CreateWalletFeeCqs return a set of cqs required for burnfee aggregation
func CreateWalletFeeCqs(dbName string) ([]*cq.ContinuousQuery, error) {
	var (
		result []*cq.ContinuousQuery
	)
	srcWalletFeeHourCqs, err := cq.NewContinuousQuery(
		"src_wallet_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
<<<<<<< HEAD
		"SELECT SUM(src_wallet_fee_amount) as sum_amount INTO wallet_fee_hour FROM trades WHERE src_rsv_addr!='' GROUP BY src_rsv_addr, wallet_addr",
=======
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s INTO %[3]s FROM %[4]s GROUP BY %[5]s, %[6]s`,
			walletschema.Amount.String(),
			walletFeeVolumeSchema.SumAmount.String(),
			common.WalletFeeVolumeMeasurementHour,
			common.WalletMeasurementName,
			walletschema.ReserveAddr.String(),
			walletschema.WalletAddr.String(),
		),
>>>>>>> 9860731... add wallet fee volume schema
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, srcWalletFeeHourCqs)

	dstWalletFeeHourCqs, err := cq.NewContinuousQuery(
		"dst_wallet_fee_hour",
		dbName,
		hourResampleInterval,
		hourResampleFor,
		"SELECT SUM(dst_wallet_fee_amount) as sum_amount INTO wallet_fee_hour FROM trades WHERE dst_rsv_addr!='' GROUP BY dst_rsv_addr, wallet_addr",
		"1h",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, dstWalletFeeHourCqs)
	srcWalletFeeDayCqs, err := cq.NewContinuousQuery(
		"src_wallet_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
<<<<<<< HEAD
		"SELECT SUM(src_wallet_fee_amount) as sum_amount INTO wallet_fee_day FROM trades WHERE src_rsv_addr!='' GROUP BY src_rsv_addr, wallet_addr",
		"1d",
=======
		fmt.Sprintf(`SELECT SUM(%[1]s) AS %[2]s INTO %[3]s FROM %[4]s GROUP BY %[5]s, %[6]s`,
			walletschema.Amount.String(),
			walletFeeVolumeSchema.SumAmount.String(),
			common.WalletFeeVolumeMeasurementDay,
			common.WalletMeasurementName,
			walletschema.ReserveAddr.String(),
			walletschema.WalletAddr.String(),
		), "1d",
>>>>>>> 9860731... add wallet fee volume schema
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, srcWalletFeeDayCqs)
	dstWalletFeeDayCqs, err := cq.NewContinuousQuery(
		"dst_wallet_fee_day",
		dbName,
		dayResampleInterval,
		dayResampleFor,
		"SELECT SUM(dst_wallet_fee_amount) as sum_amount INTO wallet_fee_day FROM trades WHERE dst_rsv_addr!='' GROUP BY dst_rsv_addr, wallet_addr",
		"1d",
		[]string{},
	)
	if err != nil {
		return nil, err
	}
	result = append(result, dstWalletFeeDayCqs)

	return result, nil
}
