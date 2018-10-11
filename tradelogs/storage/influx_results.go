package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// rowToTradeLog converts the row result of InfluxDB trade log query.
// return fields:
//- time
//- block_number
//- tx_hash
//- eth_receival_sender
//- eth_receival_amount
//- user_addr,
//- src_addr
//- dst_addr
//- src_amount
//- dst_amount
//- fiat_amount
//- reserve_addr
//- wallet_addr
//- wallet_fee
//- burn_fee
//- ip
//- country
func (is *InfluxStorage) rowToTradeLog(row []interface{}) (common.TradeLog, error) {
	// fieldsNum is number of fields returned by query, will be changed
	const fieldsNum = 17
	var tradeLog common.TradeLog

	if len(row) != fieldsNum {
		return tradeLog, fmt.Errorf("expected %d fields, got: %d", fieldsNum, len(row))
	}

	// field: time
	timeField, ok := row[0].(string)
	if !ok {
		return tradeLog, fmt.Errorf("time field is not string")
	}
	timestamp, err := time.Parse(time.RFC3339, timeField)
	if err != nil {
		return tradeLog, err
	}

	// field: block_number
	blockNumberField, ok := row[1].(string)
	if !ok {
		return tradeLog, fmt.Errorf("block number field is not string")
	}
	blockNumber, err := strconv.ParseUint(blockNumberField, 10, 64)
	if err != nil {
		return tradeLog, err
	}

	// field: tx_hash
	txHash, ok := row[2].(string)
	if !ok {
		return tradeLog, fmt.Errorf("")
	}

	ethReceivalAmount, err := row[4].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	ethReceivalAmountInWei, err := is.amountFmt.ToWei(blockchain.ETHAddr, ethReceivalAmount)
	if err != nil {
		return tradeLog, err
	}

	srcAddress := ethereum.HexToAddress(row[6].(string))
	srcAmount, err := row[8].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	srcAmountInWei, err := is.amountFmt.ToWei(srcAddress, srcAmount)
	if err != nil {
		return tradeLog, err
	}

	dstAddress := ethereum.HexToAddress(row[7].(string))
	dstAmount, err := row[9].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	dstAmountInWei, err := is.amountFmt.ToWei(dstAddress, dstAmount)
	if err != nil {
		return tradeLog, err
	}

	fiatAmount, err := row[10].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}

	walletFee, err := row[13].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	walletFeeInWei, err := is.amountFmt.ToWei(blockchain.KNCAddr, walletFee)
	if err != nil {
		return tradeLog, err
	}

	burnFee, err := row[14].(json.Number).Float64()
	if err != nil {
		return tradeLog, err
	}
	burnFeeInWei, err := is.amountFmt.ToWei(blockchain.KNCAddr, burnFee)
	if err != nil {
		return tradeLog, err
	}

	ip, ok := row[15].(string)
	if !ok {
		ip = ""
	}

	country, ok := row[16].(string)
	if err != nil {
		country = ""
	}

	tradeLog = common.TradeLog{
		Timestamp:       timestamp,
		BlockNumber:     blockNumber,
		TransactionHash: ethereum.HexToHash(txHash),

		EtherReceivalSender: ethereum.HexToAddress(row[3].(string)),
		EtherReceivalAmount: ethReceivalAmountInWei,

		UserAddress: ethereum.HexToAddress(row[5].(string)),
		SrcAddress:  srcAddress,
		DestAddress: dstAddress,
		SrcAmount:   srcAmountInWei,
		DestAmount:  dstAmountInWei,
		FiatAmount:  fiatAmount,

		ReserveAddress: ethereum.HexToAddress(row[11].(string)),
		WalletAddress:  ethereum.HexToAddress(row[12].(string)),
		WalletFee:      walletFeeInWei,
		BurnFee:        burnFeeInWei,

		IP:      ip,
		Country: country,
	}

	return tradeLog, nil
}
