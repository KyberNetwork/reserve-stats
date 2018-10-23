package storage

import (
	"fmt"
	"strconv"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

// rowToBurnFee converts the result of InfluxDB query to BurnFee event
// The query is:
// SELECT time, tx_hash, reserve_addr, amount FROM burn_fees WHERE_clause
func (is *InfluxStorage) rowToBurnFee(row []interface{}) (ethereum.Hash, common.BurnFee, error) {
	var (
		burnFee common.BurnFee
		txHash  ethereum.Hash
	)

	txHash, err := influxdb.GetTxHashFromInterface(row[1])
	if err != nil {
		return txHash, burnFee, err
	}

	reserveAddr, err := influxdb.GetAddressFromInterface(row[2])
	if err != nil {
		return txHash, burnFee, err
	}

	humanizedAmount, err := influxdb.GetFloat64FromInterface(row[3])
	if err != nil {
		return txHash, burnFee, err
	}

	weiAmount, err := is.coreClient.ToWei(blockchain.KNCAddr, humanizedAmount)
	if err != nil {
		return txHash, burnFee, err
	}

	burnFee = common.BurnFee{
		ReserveAddress: reserveAddr,
		Amount:         weiAmount,
	}

	return txHash, burnFee, nil
}

// rowToWalletFee converts the result of InfluxDB query to FeeToWallet event
// The query is:
// SELECT time, tx_hash, reserve_addr, wallet_addr, amount FROM wallet_fees WHERE_clause
func (is *InfluxStorage) rowToWalletFee(row []interface{}) (ethereum.Hash, common.WalletFee, error) {
	var (
		walletFee common.WalletFee
		txHash    ethereum.Hash
	)

	txHash, err := influxdb.GetTxHashFromInterface(row[1])
	if err != nil {
		return txHash, walletFee, err
	}

	reserveAddr, err := influxdb.GetAddressFromInterface(row[2])
	if err != nil {
		return txHash, walletFee, err
	}

	walletAddr, err := influxdb.GetAddressFromInterface(row[3])
	if err != nil {
		return txHash, walletFee, err
	}

	humanizedAmount, err := influxdb.GetFloat64FromInterface(row[4])
	if err != nil {
		return txHash, walletFee, err
	}

	weiAmount, err := is.coreClient.ToWei(blockchain.KNCAddr, humanizedAmount)
	if err != nil {
		return txHash, walletFee, err
	}

	walletFee = common.WalletFee{
		ReserveAddress: reserveAddr,
		WalletAddress:  walletAddr,
		Amount:         weiAmount,
	}

	return txHash, walletFee, nil
}

// rowToTradeLog converts the result of InfluxDB query from to TradeLog event.
// The query is:
// SELECT time, block_number, tx_hash,
// eth_receival_sender, eth_receival_amount,
// user_addr, src_addr, dst_addr, src_amount, dst_amount, (eth_amount * eth_usd_rate) as fiat_amount,
// ip, country FROM trades WHERE_clause
func (is *InfluxStorage) rowToTradeLog(row []interface{},
	burnFeesByTxHash map[ethereum.Hash][]common.BurnFee,
	walletFeesByTxHash map[ethereum.Hash][]common.WalletFee) (common.TradeLog, error) {

	var tradeLog common.TradeLog

	timestamp, err := influxdb.GetTimeFromInterface(row[0])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get timestamp: %s", err)
	}

	blockNumber, err := strconv.ParseUint(row[1].(string), 10, 64)

	txHash, err := influxdb.GetTxHashFromInterface(row[2])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get tx_hash: %s", err)
	}

	ethReceivalAddr, err := influxdb.GetAddressFromInterface(row[3])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get eth_receival_addr: %s", err)
	}

	humanizedEthReceival, err := influxdb.GetFloat64FromInterface(row[4])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get eth_receival_amount: %s", err)
	}

	ethReceivalAmountInWei, err := is.coreClient.ToWei(blockchain.ETHAddr, humanizedEthReceival)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert eth_receival_amount: %s", err)
	}

	userAddr, err := influxdb.GetAddressFromInterface(row[5])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get user_addr: %s", err)
	}

	srcAddress, err := influxdb.GetAddressFromInterface(row[6])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_addr: %s", err)
	}

	humanizedSrcAmount, err := influxdb.GetFloat64FromInterface(row[8])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_amount: %s", err)
	}

	srcAmountInWei, err := is.coreClient.ToWei(srcAddress, humanizedSrcAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert src_amount: %s", err)
	}

	dstAddress, err := influxdb.GetAddressFromInterface(row[7])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_addr: %s", err)
	}

	humanizedDstAmount, err := influxdb.GetFloat64FromInterface(row[9])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_amount: %s", err)
	}

	dstAmountInWei, err := is.coreClient.ToWei(dstAddress, humanizedDstAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert dst_amount: %s", err)
	}

	fiatAmount, err := influxdb.GetFloat64FromInterface(row[10])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get fiat_amount: %s", err)
	}

	ip, ok := row[11].(string)
	if !ok {
		ip = ""
	}

	country, ok := row[12].(string)
	if !ok {
		country = ""
	}

	tradeLog = common.TradeLog{
		Timestamp:       timestamp,
		BlockNumber:     blockNumber,
		TransactionHash: txHash,

		EtherReceivalSender: ethReceivalAddr,
		EtherReceivalAmount: ethReceivalAmountInWei,

		UserAddress: userAddr,
		SrcAddress:  srcAddress,
		DestAddress: dstAddress,
		SrcAmount:   srcAmountInWei,
		DestAmount:  dstAmountInWei,
		FiatAmount:  fiatAmount,

		BurnFees:   burnFeesByTxHash[txHash],
		WalletFees: walletFeesByTxHash[txHash],

		IP:      ip,
		Country: country,
	}

	return tradeLog, nil
}
