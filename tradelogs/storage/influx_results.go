package storage

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// rowToTradeLog converts the result of InfluxDB query from to TradeLog event.
// A trade event has is stored in database with following fields:
// -  time
// -  block_number burn_fee
// -  dst_addr
// -  dst_amount
// -  eth_receival_amount
// -  eth_receival_sender
// -  fiat_amount
// -  reserve_addr
// -  src_addr
// -  src_amount
// -  tx_hash
// -  user_addr
// -  wallet_addr
// -  wallet_fee
func (is *InfluxStorage) rowToTradeLog(row []interface{}) (common.TradeLog, error) {
	var tradeLog common.TradeLog

	timestamp, err := time.Parse(time.RFC3339, row[0].(string))
	if err != nil {
		return tradeLog, err
	}
	blockNumber, err := strconv.ParseUint(row[1].(string), 10, 64)

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
		TransactionHash: ethereum.HexToHash(row[2].(string)),

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
