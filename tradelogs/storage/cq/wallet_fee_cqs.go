package cq

import (
	"bytes"
	"text/template"

	"github.com/KyberNetwork/reserve-stats/lib/cq"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	walletschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletfee"
	walletFeeVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletfee_volume"
)

func executeWalletFeeQueryTemplate(templateString string) (string, error) {
	tmpl, err := template.New("walletFeeQuery").Parse(templateString)
	if err != nil {
		return "", err
	}
	var queryBuf bytes.Buffer
	if err := tmpl.Execute(&queryBuf, struct {
		WalletAmount                 string
		WalletSumAmount              string
		WalletFeeHourMeasurementName string
		WalletFeeDayMeasurementName  string
		WalletFeeMeasurementName     string
		ReserveAddr                  string
		WalletAddr                   string
	}{
		WalletAmount:                 walletschema.Amount.String(),
		WalletSumAmount:              walletFeeVolumeSchema.SumAmount.String(),
		WalletFeeHourMeasurementName: common.WalletFeeVolumeMeasurementHour,
		WalletFeeDayMeasurementName:  common.WalletFeeVolumeMeasurementDay,
		WalletFeeMeasurementName:     common.WalletMeasurementName,
		ReserveAddr:                  walletschema.ReserveAddr.String(),
		WalletAddr:                   walletschema.WalletAddr.String(),
	}); err != nil {
		return "", err
	}
	return "", nil
}

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
		"SELECT SUM(src_wallet_fee_amount) as sum_amount INTO wallet_fee_hour FROM trades WHERE src_rsv_addr!='' GROUP BY src_rsv_addr, wallet_addr",
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
		"SELECT SUM(src_wallet_fee_amount) as sum_amount INTO wallet_fee_day FROM trades WHERE src_rsv_addr!='' GROUP BY src_rsv_addr, wallet_addr",
		"1d",
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
