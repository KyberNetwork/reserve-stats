package storage

import (
	"fmt"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/influxdata/influxdb/models"

	"github.com/KyberNetwork/reserve-stats/lib/blockchain"
	"github.com/KyberNetwork/reserve-stats/lib/influxdb"
	"github.com/KyberNetwork/reserve-stats/tradelogs/common"
)

func (is *InfluxStorage) rowToAggregatedBurnFee(row []interface{}) (time.Time, float64, ethereum.Address, error) {
	var (
		ts      time.Time
		burnFee float64
		reserve ethereum.Address
	)

	ts, err := influxdb.GetTimeFromInterface(row[0])
	if err != nil {
		return ts, burnFee, reserve, err
	}

	burnFee, err = influxdb.GetFloat64FromInterface(row[1])
	if err != nil {
		return ts, burnFee, reserve, err
	}

	reserve, err = influxdb.GetAddressFromInterface(row[2])

	return ts, burnFee, reserve, nil
}

//this function can also work for burnFee and walletFee
func (is *InfluxStorage) rowToAggregatedFee(row []interface{}) (time.Time, float64, error) {
	var (
		ts  time.Time
		fee float64
	)
	if len(row) != 2 {
		return ts, fee, fmt.Errorf("query row len should be 2 but got %d", len(row))
	}
	ts, err := influxdb.GetTimeFromInterface(row[0])
	if err != nil {
		return ts, fee, err
	}
	fee, err = influxdb.GetFloat64FromInterface(row[1])
	if err != nil {
		return ts, fee, err
	}
	return ts, fee, nil
}

// rowToBurnFees converts the result of InfluxDB query to BurnFee event
// The query is:
// SELECT time, reserve_addr, amount, log_index FROM burn_fees WHERE_clause GROUP BY tx_hash, trade_log_index
func (is *InfluxStorage) rowToBurnFees(row models.Row) (ethereum.Hash, uint64, []common.BurnFee, error) {
	var (
		burnFees      []common.BurnFee
		txHash        ethereum.Hash
		tradeLogIndex uint64
	)

	txHash, err := influxdb.GetTxHashFromInterface(row.Tags["tx_hash"])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	tradeLogIndex, err = influxdb.GetUint64FromTagValue(row.Tags["trade_log_index"])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	for _, value := range row.Values {
		reserveAddr, err := influxdb.GetAddressFromInterface(value[1])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		humanizedAmount, err := influxdb.GetFloat64FromInterface(value[2])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		weiAmount, err := is.coreClient.ToWei(blockchain.KNCAddr, humanizedAmount)
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		logIndex, err := influxdb.GetUint64FromTagValue(value[3])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		burnFee := common.BurnFee{
			ReserveAddress: reserveAddr,
			Amount:         weiAmount,
			Index:          uint(logIndex),
		}

		burnFees = append(burnFees, burnFee)
	}

	return txHash, tradeLogIndex, burnFees, nil
}

// rowToWalletFees converts the result of InfluxDB query to FeeToWallet event
// The query is:
// SELECT time, reserve_addr, wallet_addr, amount, log_index FROM wallet_fees WHERE_clause GROUP BY tx_hash, trade_log_index
func (is *InfluxStorage) rowToWalletFees(row models.Row) (ethereum.Hash, uint64, []common.WalletFee, error) {
	var (
		walletFees    []common.WalletFee
		txHash        ethereum.Hash
		tradeLogIndex uint64
	)

	txHash, err := influxdb.GetTxHashFromInterface(row.Tags["tx_hash"])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	tradeLogIndex, err = influxdb.GetUint64FromTagValue(row.Tags["trade_log_index"])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	for _, value := range row.Values {
		reserveAddr, err := influxdb.GetAddressFromInterface(value[1])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		walletAddr, err := influxdb.GetAddressFromInterface(value[2])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		humanizedAmount, err := influxdb.GetFloat64FromInterface(value[3])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		weiAmount, err := is.coreClient.ToWei(blockchain.KNCAddr, humanizedAmount)
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		logIndex, err := influxdb.GetUint64FromTagValue(value[4])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		walletFee := common.WalletFee{
			ReserveAddress: reserveAddr,
			WalletAddress:  walletAddr,
			Amount:         weiAmount,
			Index:          uint(logIndex),
		}

		walletFees = append(walletFees, walletFee)
	}

	return txHash, tradeLogIndex, walletFees, nil
}

// rowToTradeLog converts the result of InfluxDB query from to TradeLog event.
// The query is:
// SELECT time, block_number,
// eth_receival_sender, eth_receival_amount,
// user_addr, src_addr, dst_addr, src_amount, dst_amount, (eth_amount * eth_usd_rate) as fiat_amount,
// ip, country FROM trades WHERE_clause GROUP BY tx_hash, log_index
func (is *InfluxStorage) rowToTradeLog(row models.Row,
	burnFeesByTxHash map[ethereum.Hash]map[uint][]common.BurnFee,
	walletFeesByTxHash map[ethereum.Hash]map[uint][]common.WalletFee) (common.TradeLog, error) {

	var tradeLog common.TradeLog

	txHash, err := influxdb.GetTxHashFromInterface(row.Tags["tx_hash"])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get tx_hash: %s", err)
	}

	logIndex, err := influxdb.GetUint64FromTagValue(row.Tags["log_index"])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get trade log index: %s", err)
	}

	value := row.Values[0]

	timestamp, err := influxdb.GetTimeFromInterface(value[0])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get timestamp: %s", err)
	}

	blockNumber, err := strconv.ParseUint(value[1].(string), 10, 64)

	ethReceivalAddr, err := influxdb.GetAddressFromInterface(value[2])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get eth_receival_addr: %s", err)
	}

	humanizedEthReceival, err := influxdb.GetFloat64FromInterface(value[3])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get eth_receival_amount: %s", err)
	}

	ethReceivalAmountInWei, err := is.coreClient.ToWei(blockchain.ETHAddr, humanizedEthReceival)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert eth_receival_amount: %s", err)
	}

	userAddr, err := influxdb.GetAddressFromInterface(value[4])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get user_addr: %s", err)
	}

	srcAddress, err := influxdb.GetAddressFromInterface(value[5])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_addr: %s", err)
	}

	humanizedSrcAmount, err := influxdb.GetFloat64FromInterface(value[7])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get src_amount: %s", err)
	}

	srcAmountInWei, err := is.coreClient.ToWei(srcAddress, humanizedSrcAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert src_amount: %s", err)
	}

	dstAddress, err := influxdb.GetAddressFromInterface(value[6])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_addr: %s", err)
	}

	humanizedDstAmount, err := influxdb.GetFloat64FromInterface(value[8])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get dst_amount: %s", err)
	}

	dstAmountInWei, err := is.coreClient.ToWei(dstAddress, humanizedDstAmount)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert dst_amount: %s", err)
	}

	fiatAmount, err := influxdb.GetFloat64FromInterface(value[9])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get fiat_amount: %s", err)
	}

	ip, ok := value[10].(string)
	if !ok {
		ip = ""
	}

	country, ok := value[11].(string)
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

		BurnFees:   burnFeesByTxHash[txHash][uint(logIndex)],
		WalletFees: walletFeesByTxHash[txHash][uint(logIndex)],

		IP:      ip,
		Country: country,

		Index: uint(logIndex),
	}

	return tradeLog, nil
}
