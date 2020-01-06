package influx

import (
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	burnVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/burnfee_volume"
	logschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/tradelog"
	volSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/volume"
	walletFeeVolumeSchema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/influx/schema/walletfee_volume"
)

func (is *Storage) rowToAggregatedBurnFee(row []interface{}, idxs map[burnVolumeSchema.FieldName]int) (time.Time, float64, ethereum.Address, error) {
	var (
		ts      time.Time
		burnFee float64
		reserve ethereum.Address
	)

	ts, err := influxdb.GetTimeFromInterface(row[idxs[burnVolumeSchema.Time]])
	if err != nil {
		return ts, burnFee, reserve, err
	}

	burnFee, err = influxdb.GetFloat64FromInterface(row[idxs[burnVolumeSchema.SumAmount]])
	if err != nil {
		return ts, burnFee, reserve, err
	}

	switch {
	case row[idxs[burnVolumeSchema.SrcReserveAddr]] != nil && row[idxs[burnVolumeSchema.DstReserveAddr]] != nil:
		panic("Logic fault : there should not be a record with both source and dest reserve address")
	case row[idxs[burnVolumeSchema.SrcReserveAddr]] != nil:
		reserve, err = influxdb.GetAddressFromInterface(row[idxs[burnVolumeSchema.SrcReserveAddr]])
		if err != nil {
			return ts, burnFee, reserve, err
		}
	case row[idxs[burnVolumeSchema.DstReserveAddr]] != nil:
		reserve, err = influxdb.GetAddressFromInterface(row[idxs[burnVolumeSchema.DstReserveAddr]])
		if err != nil {
			return ts, burnFee, reserve, err
		}
	default:
		panic("Logic fault : there should not be a record with nil source and dest reserve address")
	}

	return ts, burnFee, reserve, nil
}

func (is *Storage) rowToAggregatedUserVolume(row []interface{}, idxs volSchema.FieldsRegistrar) (time.Time, float64, float64, error) {
	var (
		ts        time.Time
		ethAmount float64
		usdAmount float64
		err       error
	)
	ts, err = influxdb.GetTimeFromInterface(row[idxs[volSchema.Time]])
	if err != nil {
		return ts, ethAmount, usdAmount, err
	}

	ethAmount, err = influxdb.GetFloat64FromInterface(row[idxs[volSchema.ETHVolume]])
	if err != nil {
		return ts, ethAmount, usdAmount, err
	}

	usdAmount, err = influxdb.GetFloat64FromInterface(row[idxs[volSchema.USDVolume]])
	if err != nil {
		return ts, ethAmount, usdAmount, err
	}
	return ts, ethAmount, usdAmount, err
}

func (is *Storage) rowToUserInfo(row []interface{}) (float64, float64, error) {
	var (
		ethAmount, usdAmount float64
		err                  error
	)
	ethAmount, err = influxdb.GetFloat64FromInterface(row[1])
	if err != nil {
		return ethAmount, usdAmount, err
	}
	usdAmount, err = influxdb.GetFloat64FromInterface(row[2])
	if err != nil {
		return ethAmount, usdAmount, err
	}
	return ethAmount, usdAmount, nil
}

//this function can also work for burnFee and walletFee
func (is *Storage) rowToAggregatedFee(row []interface{}, idxs walletFeeVolumeSchema.FieldsRegistrar) (time.Time, float64, error) {
	var (
		ts  time.Time
		fee float64
	)
	if len(row) != 2 {
		return ts, fee, fmt.Errorf("query row len should be 2 but got %d", len(row))
	}
	ts, err := influxdb.GetTimeFromInterface(row[idxs[walletFeeVolumeSchema.Time]])
	if err != nil {
		return ts, fee, err
	}
	fee, err = influxdb.GetFloat64FromInterface(row[idxs[walletFeeVolumeSchema.SumAmount]])
	if err != nil {
		return ts, fee, err
	}
	return ts, fee, nil
}

// rowToTradeLog converts the result of InfluxDB query from to TradeLog event.
// The query is:
// SELECT time, block_number,
// eth_receival_sender, eth_receival_amount,
// user_addr, src_addr, dst_addr, src_amount, dst_amount, (eth_amount * eth_usd_rate) as fiat_amount,
// ip, country FROM trades WHERE_clause GROUP BY tx_hash, log_index
func (is *Storage) rowToTradeLog(value []interface{},
	idxs logschema.FieldsRegistrar) (common.TradeLog, error) {

	var (
		tradeLog          common.TradeLog
		dstReserveAddress ethereum.Address
		srcReserveAddress ethereum.Address
	)

	txHash, err := influxdb.GetTxHashFromInterface(value[idxs[logschema.TxHash]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get tx_hash: %s", err)
	}

	logIndex, err := influxdb.GetUint64FromTagValue(value[idxs[logschema.LogIndex]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get trade log index: %s", err)
	}

	timestamp, err := influxdb.GetTimeFromInterface(value[idxs[logschema.Time]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get timestamp: %s", err)
	}

	blockNumber, err := influxdb.GetUint64FromInterface(value[idxs[logschema.BlockNumber]])

	if err != nil {
		return tradeLog, fmt.Errorf("failed to get blockNumber: %s", err)
	}

	ethAmount, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.EthAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get ethAmount: %s", err)
	}

	ethAmountInWei, err := is.tokenAmountFormatter.ToWei(blockchain.ETHAddr, ethAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get ethReceivalAmount: %s", err)
	}

	originalEthAmount, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.OriginalEthAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get original ethAmount: %s", err)
	}

	originalEthAmountInWei, err := is.tokenAmountFormatter.ToWei(blockchain.ETHAddr, originalEthAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get original ethReceivalAmount: %s", err)
	}

	userAddr, err := influxdb.GetAddressFromInterface(value[idxs[logschema.UserAddr]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get user_addr: %s", err)
	}

	srcAddress, err := influxdb.GetAddressFromInterface(value[idxs[logschema.SrcAddr]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_addr: %s", err)
	}

	dstAddress, err := influxdb.GetAddressFromInterface(value[idxs[logschema.DstAddr]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_addr: %s", err)
	}
	humanizedSrcAmount, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.SrcAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_amount: %s", err)
	}

	srcAmountInWei, err := is.tokenAmountFormatter.ToWei(srcAddress, humanizedSrcAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert src_amount: %s", err)
	}

	humanizedDstAmount, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.DstAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_amount: %s", err)
	}

	if value[idxs[logschema.SrcReserveAddr]] != nil {
		srcReserveAddress, err = influxdb.GetAddressFromInterface(value[idxs[logschema.SrcReserveAddr]])
		if err != nil {
			return tradeLog, fmt.Errorf("failed to get src_reserve_addr: %s", err.Error())
		}
	}

	if value[idxs[logschema.DstReserveAddr]] != nil {
		dstReserveAddress, err = influxdb.GetAddressFromInterface(value[idxs[logschema.DstReserveAddr]])
		if err != nil {
			return tradeLog, fmt.Errorf("failed to get dst_reserve_addr: %s", err.Error())
		}
	}

	dstAmountInWei, err := is.tokenAmountFormatter.ToWei(dstAddress, humanizedDstAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert dst_amount: %s", err)
	}

	ip, ok := value[idxs[logschema.IP]].(string)
	if !ok {
		ip = ""
	}

	country, ok := value[idxs[logschema.Country]].(string)
	if !ok {
		country = ""
	}

	uid, ok := value[idxs[logschema.UID]].(string)
	if !ok {
		uid = ""
	}

	appName, ok := value[idxs[logschema.IntegrationApp]].(string)
	if !ok {
		appName = ""
	}

	fiatAmount, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.FiatAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get fiat_amount: %s", err)
	}
	srcBurnFee, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.SourceBurnAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_burn_amount: %s", err)
	}

	dstBurnFee, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.DestBurnAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_burn_amount: %s", err)
	}

	srcWalletFee, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.SourceWalletFeeAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_wallet_fee_amount: %s", err)
	}

	dstWalletFee, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.DestWalletFeeAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_wallet_fee_amount: %s", err)
	}

	walletAddr, err := influxdb.GetAddressFromInterface(value[idxs[logschema.WalletAddress]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get wallet_addr: %s", err)
	}

	walletName, ok := value[idxs[logschema.WalletName]].(string)
	if !ok {
		walletName = ""
	}

	txSender, err := influxdb.GetAddressFromInterface(value[idxs[logschema.TxSender]])
	if err != nil {
		return tradeLog, fmt.Errorf("failded to get tx_sender: %s", err)
	}
	receiverAddr, err := influxdb.GetAddressFromInterface(value[idxs[logschema.ReceiverAddress]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get receiver_addr: %s", err)
	}

	tradeLog = common.TradeLog{
		Timestamp:       timestamp,
		BlockNumber:     blockNumber,
		TransactionHash: txHash,

		EthAmount:         ethAmountInWei,
		OriginalEthAmount: originalEthAmountInWei,

		UserAddress:       userAddr,
		SrcAddress:        srcAddress,
		DestAddress:       dstAddress,
		SrcReserveAddress: srcReserveAddress,
		DstReserveAddress: dstReserveAddress,
		SrcAmount:         srcAmountInWei,
		DestAmount:        dstAmountInWei,
		FiatAmount:        fiatAmount,
		WalletAddress:     walletAddr,
		WalletName:        walletName,

		SrcBurnAmount:      srcBurnFee,
		DstBurnAmount:      dstBurnFee,
		SrcWalletFeeAmount: srcWalletFee,
		DstWalletFeeAmount: dstWalletFee,

		IP:             ip,
		Country:        country,
		UID:            uid,
		IntegrationApp: appName,
		Index:          uint(logIndex),

		TxSender:        txSender,
		ReceiverAddress: receiverAddr,
	}

	return tradeLog, nil
}
