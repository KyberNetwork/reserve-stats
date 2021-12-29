package storage

import (
	"database/sql"
	_ "embed"
	"fmt"
	"math/big"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/zerox"
)

// ZeroxStorage ...
type ZeroxStorage struct {
	sugar *zap.SugaredLogger
	db    *sqlx.DB
}

//go:embed schema.sql
var schema string

// NewZeroxStorage ...
func NewZeroxStorage(db *sqlx.DB, sugar *zap.SugaredLogger) (*ZeroxStorage, error) {
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &ZeroxStorage{
		db:    db,
		sugar: sugar,
	}, nil
}

// InsertTradelogs ...
func (zs *ZeroxStorage) InsertTradelogs(tradelogs []zerox.Tradelog) error {
	var (
		tx, inputTokenSymbol, inputTokenAddress, outputTokenSymbol, outputTokenAddress, takerAddress []string
		timestamp                                                                                    []int64
		inputTokenAmount, outputTokenAmount                                                          []float64
	)
	query := `INSERT INTO tradelogs (tx, timestamp, input_token_symbol, input_token_address, input_token_amount, output_token_symbol, output_token_address, output_token_amount, taker_address)
	VALUES(
		unnest($1::TEXT[]),
		unnest($2::BIGINT[]),
		unnest($3::TEXT[]),
		unnest($4::TEXT[]),
		unnest($5::FLOAT[]),
		unnest($6::TEXT[]),
		unnest($7::TEXT[]),
		unnest($8::FLOAT[]),
		unnest($9::TEXT[])
	) ON CONFLICT (tx) DO NOTHING;`

	for _, trade := range tradelogs {
		tx = append(tx, trade.Transaction.ID)
		ts, err := strconv.ParseInt(trade.Timestamp, 10, 64)
		if err != nil {
			zs.sugar.Errorw("failed to parse timestamp", "error", err)
			return err
		}
		timestamp = append(timestamp, ts)
		inputTokenSymbol = append(inputTokenSymbol, trade.InputToken.Symbol)
		inputTokenAddress = append(inputTokenAddress, trade.InputToken.ID)
		inputAmount, err := zs.convertAmount(trade.InputTokenAmount, trade.InputToken.Decimals)
		if err != nil {
			zs.sugar.Errorw("faield to get input amount", "error", err)
			return err
		}
		inputTokenAmount = append(inputTokenAmount, inputAmount)
		outputTokenSymbol = append(outputTokenSymbol, trade.OutputToken.Symbol)
		outputTokenAddress = append(outputTokenAddress, trade.OutputToken.ID)
		outputAmount, err := zs.convertAmount(trade.OutputTokenAmount, trade.OutputToken.Decimals)
		if err != nil {
			zs.sugar.Errorw("faield to get output amount", "error", err)
			return err
		}
		outputTokenAmount = append(outputTokenAmount, outputAmount)
		takerAddress = append(takerAddress, trade.Taker.ID)
	}

	if _, err := zs.db.Exec(query, pq.Array(tx), pq.Array(timestamp), pq.Array(inputTokenSymbol), pq.Array(inputTokenAddress), pq.Array(inputTokenAmount),
		pq.Array(outputTokenSymbol), pq.Array(outputTokenAddress), pq.Array(outputTokenAmount), pq.Array(takerAddress)); err != nil {
		zs.sugar.Errorw("failed to insert trades", "error", err)
		return err
	}
	return nil
}

func (zs *ZeroxStorage) convertAmount(amountStr, decimalsStr string) (float64, error) {
	decimals, err := strconv.ParseInt(decimalsStr, 10, 64)
	if err != nil {
		zs.sugar.Errorw("failed to parse decimals", "error", err)
		return 0, err
	}
	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	amountFloat, ok := big.NewFloat(0).SetString(amountStr)
	if !ok {
		zs.sugar.Error("failed to parse amount")
		return 0, fmt.Errorf("failed to parse amount")
	}
	amountF := new(big.Float).Quo(amountFloat, big.NewFloat(0).SetInt(pow))
	amount, _ := amountF.Float64()
	return amount, nil
}

// GetLastTradeTimestamp ...
func (zs *ZeroxStorage) GetLastTradeTimestamp() (int64, error) {
	var (
		timestamp int64
	)
	query := `SELECT COALESCE(MAX(timestamp), 0) FROM tradelogs;`
	if err := zs.db.Get(&timestamp, query); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return timestamp, nil
}

// InsertConvertTrades ...
func (zs *ZeroxStorage) InsertConvertTrades(convertTrades zerox.ConvertTrades) error {
	query := `INSERT INTO convert_trades(original_symbol, symbol, price, timestamp, original_trade, convert_trade, in_token, in_token_amount, out_token, out_token_amount)
	VALUES (
		unnest($1::TEXT[]),
		unnest($2::TEXT[]),
		unnest($3::FLOAT[]),
		unnest($4::BIGINT[]),
		unnest($5::JSONB[]),
		unnest($6::JSONB[]),
		unnest($7::TEXT[]),
		unnest($8::FLOAT[]),
		unnest($9::TEXT[]),
		unnest($10::FLOAT[])
	) ON CONFLICT (symbol, timestamp) DO NOTHING;`
	if _, err := zs.db.Exec(query, pq.Array(convertTrades.OriginalSymbols), pq.Array(convertTrades.Symbols),
		pq.Array(convertTrades.Prices), pq.Array(convertTrades.Timestamps), pq.Array(convertTrades.OriginalTrades), pq.Array(convertTrades.Trades),
		pq.Array(convertTrades.InToken), pq.Array(convertTrades.InTokenAmount), pq.Array(convertTrades.OutToken), pq.Array(convertTrades.OutTokenAmount)); err != nil {
		zs.sugar.Errorw("failed to insert convert trades", "error", err)
		return err
	}
	return nil
}

// GetConvertTrades ...
func (zs *ZeroxStorage) GetConvertTrades(fromTime, toTime int64) ([]zerox.ConvertTrade, error) {
	var (
		result []zerox.ConvertTrade
	)
	const query = `SELECT symbol, price, timestamp FROM convert_trades WHERE timestamp >= $1 AND timestamp <= $2;`
	err := zs.db.Select(&result, query, fromTime, toTime)
	return result, err
}

// Get0xTrades ...
func (zs *ZeroxStorage) Get0xTrades(fromTime, toTime int64) ([]zerox.SimpleTradelog, error) {
	var (
		result []zerox.SimpleTradelog
	)
	query := `SELECT tx, timestamp, input_token_symbol, input_token_amount, output_token_symbol, output_token_amount, taker_address FROM tradelogs WHERE timestamp * 1000 >= $1 AND timestamp * 1000 <= $2;`
	if err := zs.db.Select(&result, query, fromTime, toTime); err != nil {
		return nil, err
	}
	return result, nil
}

// GetConvertTradeInfo ...
func (zs *ZeroxStorage) GetConvertTradeInfo(fromTime, toTime int64) ([]zerox.ConvertTradeInfo, error) {
	var (
		result []zerox.ConvertTradeInfo
	)
	const query = `WITH 
intoken AS (SELECT price AS in_token_rate, timestamp FROM convert_trades WHERE symbol = concat(original_trade->'inputToken'->>'symbol','USDT')),
outtoken AS (SELECT price as out_token_rate, timestamp FROM convert_trades WHERE symbol = concat(original_trade->'outputToken'->>'symbol','USDT')),
ethtoken AS (SELECT symbol, price as eth_usdt_rate, timestamp,in_token, in_token_amount, out_token, out_token_amount,
	original_trade->'transaction'->>'id' as tx_hash, original_trade->'taker'->>'id' as taker  FROM convert_trades WHERE symbol = 'ETHUSDT')
SELECT in_token, COALESCE(in_token_rate, 0) AS in_token_rate, in_token_amount, eth_usdt_rate as eth_rate, out_token, out_token_amount, COALESCE(out_token_rate, 0) AS out_token_rate, ethtoken.timestamp,
tx_hash, taker, '0xRFQ' as account_name
FROM ethtoken
FULL JOIN intoken ON intoken.timestamp = ethtoken.timestamp
FULL JOIN outtoken ON ethtoken.timestamp = outtoken.timestamp
WHERE ethtoken.timestamp >= $1 AND ethtoken.timestamp <= $2;`
	err := zs.db.Select(&result, query, fromTime, toTime)
	return result, err
}

// GetBinanceConvertTradeInfo ...
func (zs *ZeroxStorage) GetBinanceConvertTradeInfo(fromTime, toTime int64) ([]zerox.ConvertTradeInfo, error) {
	var (
		temp   []zerox.CexConvertTradeInfo
		result []zerox.ConvertTradeInfo
	)
	const query = `
SELECT original_symbol AS in_token, price as eth_rate, binance_convert_to_eth_price.timestamp as timestamp, original_trade->>'qty' AS in_token_amount, original_trade->>'price' AS in_token_rate,
original_trade->>'isBuyer' as is_buyer, account as account_name
FROM binance_convert_to_eth_price
JOIN binance_trades 
ON binance_convert_to_eth_price.original_trade->'id' = binance_trades."data"->'id'
AND binance_convert_to_eth_price.original_symbol = binance_trades.symbol
WHERE binance_convert_to_eth_price.timestamp >= $1 AND binance_convert_to_eth_price.timestamp <= $2;`
	err := zs.db.Select(&temp, query, fromTime, toTime)
	for _, t := range temp {
		inTokenAmount, err := strconv.ParseFloat(t.InTokenAmount, 64)
		if err != nil {
			return nil, err
		}
		inTokenRate, err := strconv.ParseFloat(t.InTokenRate, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, zerox.ConvertTradeInfo{
			AccountName:   t.AccountName,
			Timestamp:     t.Timestamp,
			ETHRate:       t.ETHRate,
			InToken:       t.InToken,
			InTokenAmount: inTokenAmount,
			InTokenRate:   inTokenRate,
			IsBuyer:       t.IsBuyer,
			TxHash:        t.TxHash,
		})
	}
	return result, err
}
