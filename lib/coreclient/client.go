package coreclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/binance"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
	"go.uber.org/zap"
)

// CoreClient client for core
type CoreClient struct {
	Endpoint  string `json:"endpoint"`
	APIKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
	sugar     *zap.SugaredLogger
	client    *http.Client
}

// NewCoreClient return new core client
func NewCoreClient(endpoint, apiKey, secretKey string, sugar *zap.SugaredLogger) *CoreClient {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &CoreClient{
		Endpoint:  endpoint,
		APIKey:    apiKey,
		SecretKey: secretKey,
		sugar:     sugar,
		client:    client,
	}
}

// TradingPairSymbols is a pair of token trading
type TradingPairSymbols struct {
	TradingPair
	BaseSymbol  string `json:"base_symbol"`
	QuoteSymbol string `json:"quote_symbol"`
}

// TradingPair ...
type TradingPair struct {
	ID              uint64  `json:"id"`
	Base            uint64  `json:"base"`
	Quote           uint64  `json:"quote"`
	PricePrecision  uint64  `json:"price_precision"`
	AmountPrecision uint64  `json:"amount_precision"`
	AmountLimitMin  float64 `json:"amount_limit_min"`
	AmountLimitMax  float64 `json:"amount_limit_max"`
	PriceLimitMin   float64 `json:"price_limit_min"`
	PriceLimitMax   float64 `json:"price_limit_max"`
	MinNotional     float64 `json:"min_notional"`
}

func (cc *CoreClient) sendRequest(method, endpoint string, params map[string]string, signNeeded bool,
	timepoint time.Time) ([]byte, error) {

	var (
		respBody []byte
		logger   = cc.sugar.With("func", "coreclient.sendRequest")
	)
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	logger.Infow("request", "url", req.URL)

	req, err = httputil.Sign(req, cc.APIKey, cc.SecretKey)
	if err != nil {
		return respBody, err
	}

	resp, err := cc.client.Do(req)
	if err != nil {
		return respBody, err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			logger.Errorw("Response body close error", "error", cErr.Error())
		}
	}()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorw("failed to read response body", "error", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info("get token pair successfully")
	default:
		err = fmt.Errorf("return with code: %d, content: %s", resp.StatusCode, respBody)
	}
	return respBody, err
}

// GetBinanceSupportedPairs ...
func (cc *CoreClient) GetBinanceSupportedPairs(exchangeID int64) ([]binance.Symbol, error) {
	var (
		coreResponse struct {
			Data []TradingPairSymbols `json:"data"`
		}
		result []binance.Symbol
	)
	endpoint := fmt.Sprintf("%s/v3/trading-pair", cc.Endpoint)
	params := map[string]string{
		"exchange_id": strconv.FormatInt(exchangeID, 10),
	}

	res, err := cc.sendRequest(
		http.MethodGet,
		endpoint,
		params,
		true,
		time.Now(),
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &coreResponse)
	for _, pair := range coreResponse.Data {
		result = append(result, binance.Symbol{
			Symbol:     pair.BaseSymbol + pair.QuoteSymbol,
			BaseAsset:  pair.BaseSymbol,
			QuoteAsset: pair.QuoteSymbol,
		})
	}
	return result, err
}
