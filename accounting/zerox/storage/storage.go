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
		tx, inputTokenSymbol, inputTokenAddress, outputTokenSymbol, outputTokenAddress []string
		timestamp                                                                      []int64
		inputTokenAmount, outputTokenAmount                                            []float64
	)
	query := `INSERT INTO tradelogs (tx, timestamp, input_token_symbol, input_token_address, input_token_amount, output_token_symbol, output_token_address, output_token_amount)
	VALUES(
		unnest($1::TEXT[]),
		unnest($2::BIGINT[]),
		unnest($3::TEXT[]),
		unnest($4::TEXT[]),
		unnest($5::FLOAT[]),
		unnest($6::TEXT[]),
		unnest($7::TEXT[]),
		unnest($8::FLOAT[])
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
	}

	if _, err := zs.db.Exec(query, pq.Array(tx), pq.Array(timestamp), pq.Array(inputTokenSymbol), pq.Array(inputTokenAddress), pq.Array(inputTokenAmount),
		pq.Array(outputTokenSymbol), pq.Array(outputTokenAddress), pq.Array(outputTokenAmount)); err != nil {
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
	query := `INSERT INTO convert_trades(original_symbol, symbol, price, timestamp, original_trade, convert_trade)
	VALUES (
		unnest($1::TEXT[]),
		unnest($2::TEXT[]),
		unnest($3::FLOAT[]),
		unnest($4::BIGINT[]),
		unnest($5::JSONB[]),
		unnest($6::JSONB[])
	);`
	if _, err := zs.db.Exec(query, pq.Array(convertTrades.OriginalSymbols), pq.Array(convertTrades.Symbols),
		pq.Array(convertTrades.Prices), pq.Array(convertTrades.Timestamps), pq.Array(convertTrades.OriginalTrades), pq.Array(convertTrades.Trades)); err != nil {
		zs.sugar.Errorw("failed to insert convert trades", "error", err)
		return err
	}
	return nil
}

// GetConvertTrades ...
func (zs *ZeroxStorage) GetConvertTrades(fromTime, toTime int64) ([]zerox.ConvertTrade, error) {
	var (
		result []zerox.ConvertTrade
		err    error
	)
	const query = `SELECT symbol, price, timestamp FROM convert_trades WHERE timestamp >= $1 AND timestamp <= $2;`
	if err := zs.db.Select(&result, query, fromTime, toTime); err != nil {
		zs.sugar.Errorw("failed to get convert eth price", "error", err)
		return result, err
	}
	return result, err
}
