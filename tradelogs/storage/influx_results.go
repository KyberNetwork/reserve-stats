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
	burnschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/burnfee"
	logschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/tradelog"
	walletschema "github.com/KyberNetwork/reserve-stats/tradelogs/storage/schema/walletfee"
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

func (is *InfluxStorage) rowToAggregatedUserVolume(row []interface{}) (time.Time, float64, float64, error) {
	var (
		ts        time.Time
		ethAmount float64
		usdAmount float64
		err       error
	)
	ts, err = influxdb.GetTimeFromInterface(row[0])
	if err != nil {
		return ts, ethAmount, usdAmount, err
	}

	ethAmount, err = influxdb.GetFloat64FromInterface(row[1])
	if err != nil {
		return ts, ethAmount, usdAmount, err
	}

	usdAmount, err = influxdb.GetFloat64FromInterface(row[2])
	if err != nil {
		return ts, ethAmount, usdAmount, err
	}
	return ts, ethAmount, usdAmount, err
}

func (is *InfluxStorage) rowToUserInfo(row []interface{}) (float64, float64, error) {
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

	txHash, err := influxdb.GetTxHashFromInterface(row.Tags[burnschema.TxHash.String()])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	tradeLogIndex, err = influxdb.GetUint64FromTagValue(row.Tags[burnschema.TradeLogIndex.String()])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	idxs, err := burnschema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}
	for _, value := range row.Values {
		reserveAddr, err := influxdb.GetAddressFromInterface(value[idxs[burnschema.ReserveAddr]])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		humanizedAmount, err := influxdb.GetFloat64FromInterface(value[idxs[burnschema.Amount]])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		weiAmount, err := is.tokenAmountFormatter.ToWei(blockchain.KNCAddr, humanizedAmount)
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		logIndex, err := influxdb.GetUint64FromTagValue(value[idxs[burnschema.LogIndex]])
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

	txHash, err := influxdb.GetTxHashFromInterface(row.Tags[walletschema.TxHash.String()])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	tradeLogIndex, err = influxdb.GetUint64FromTagValue(row.Tags[walletschema.TradeLogIndex.String()])
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	idxs, err := walletschema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return txHash, tradeLogIndex, nil, err
	}

	for _, value := range row.Values {
		reserveAddr, err := influxdb.GetAddressFromInterface(value[idxs[walletschema.ReserveAddr]])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		walletAddr, err := influxdb.GetAddressFromInterface(value[idxs[walletschema.WalletAddr]])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		humanizedAmount, err := influxdb.GetFloat64FromInterface(value[idxs[walletschema.Amount]])
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		weiAmount, err := is.tokenAmountFormatter.ToWei(blockchain.KNCAddr, humanizedAmount)
		if err != nil {
			return txHash, tradeLogIndex, nil, err
		}

		logIndex, err := influxdb.GetUint64FromTagValue(value[idxs[walletschema.LogIndex]])
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

func (is *InfluxStorage) rowToCountryStats(row []interface{}) (time.Time, common.CountryStats, error) {
	var (
		ts           time.Time
		countryStats common.CountryStats
		err          error
	)
	return ts, countryStats, err
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

	txHash, err := influxdb.GetTxHashFromInterface(row.Tags[logschema.TxHash.String()])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get tx_hash: %s", err)
	}

	logIndex, err := influxdb.GetUint64FromTagValue(row.Tags[logschema.LogIndex.String()])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get trade log index: %s", err)
	}

	value := row.Values[0]
	idxs, err := logschema.NewFieldsRegistrar(row.Columns)
	if err != nil {
		return tradeLog, err
	}

	timestamp, err := influxdb.GetTimeFromInterface(value[idxs[logschema.Time]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get timestamp: %s", err)
	}

	blockNumber, err := strconv.ParseUint(value[idxs[logschema.BlockNumber]].(string), 10, 64)

	ethReceivalAddr, err := influxdb.GetAddressFromInterface(value[idxs[logschema.EthReceivalSender]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get eth_receival_addr: %s", err)
	}

	humanizedEthReceival, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.EthReceivalAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get eth_receival_amount: %s", err)
	}

	ethReceivalAmountInWei, err := is.tokenAmountFormatter.ToWei(blockchain.ETHAddr, humanizedEthReceival)
	if err != nil {
		return tradeLog, fmt.Errorf("failed to convert eth_receival_amount: %s", err)
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

	appName, ok := value[idxs[logschema.IntegrationApp]].(string)
	if !ok {
		appName = ""
	}

	fiatAmount, err := influxdb.GetFloat64FromInterface(value[idxs[logschema.FiatAmount]])
	if err != nil {
		return tradeLog, fmt.Errorf("failed to get fiat_amount: %s", err)
	}
	burnFees := burnFeesByTxHash[txHash][uint(logIndex)]
	if burnFees == nil {
		burnFees = []common.BurnFee{}
	}

	walletFees := walletFeesByTxHash[txHash][uint(logIndex)]
	if walletFees == nil {
		walletFees = []common.WalletFee{}
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

		BurnFees:   burnFees,
		WalletFees: walletFees,

		IP:             ip,
		Country:        country,
		IntegrationApp: appName,
		Index:          uint(logIndex),
	}

	return tradeLog, nil
}
