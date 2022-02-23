package http

import (
	"fmt"
	"math"
	"strings"
)

func convertRateToBinance(inAmount, outAmount float64, inToken, outToken string) (string, string, float64) {
	var (
		in, out     = math.MaxInt64, math.MaxInt64
		quoteTokens = []string{"DAI", "USDT", "BUSD", "USDC", "BTC", "WBTC", "WETH", "ETH"}
	)
	for i, t := range quoteTokens {
		if inToken == t {
			in = i
			continue
		}
		if outToken == t {
			out = i
			continue
		}
	}
	if in == out {
		return "", "", 0
	}
	if in < out {
		symbol := fmt.Sprintf("%s%s", outToken, inToken)
		side := "ask"
		rate := inAmount / outAmount
		return symbol, side, rate
	}
	symbol := fmt.Sprintf("%s%s", inToken, outToken)
	side := "bid"
	rate := outAmount / inAmount
	return symbol, side, rate
}

// avgPrice calculate average price of a slice of trade
func avgPrice(trades []*ConvertTrade) float64 {
	var (
		rate float64
	)
	if len(trades) == 0 {
		return 0
	}
	for _, t := range trades {
		rate += t.Rate
	}
	return rate / float64(len(trades))
}

// return ethChange, tokenChange and ty
func getAmountAndType(symbol, side string, ethAmount, qty float64) (string, float64, float64) {
	var (
		ethChange, tokenChange float64
	)
	tradeType := sellType
	tokenChange = qty * -1
	ethChange = ethAmount
	if side == askSide {
		tradeType = buyType
		if strings.HasSuffix(symbol, eth) {
			ethChange = ethAmount * -1
			tokenChange = qty
		}
	} else if strings.HasPrefix(symbol, eth) {
		ethChange = ethAmount * -1
		tokenChange = qty
	}
	return tradeType, ethChange, tokenChange
}

// return tradetype, ethChange, tokenChange for onchain trades
// same with get amount and type offchain and invert
func getAmountAndTypeOnchain(symbol, side string, ethAmount, qty float64) (string, float64, float64) {
	tradeType, ethChange, tokenChange := getAmountAndType(symbol, side, ethAmount, qty)
	if tradeType == buyType {
		tradeType = sellType
	} else {
		tradeType = buyType
	}
	return tradeType, ethChange * -1, tokenChange * -1
}
