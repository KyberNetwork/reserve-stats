package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	endpointPrefix      = "https://api.binance.com"
	badAPIKeyFormatCode = -2014
	rejfectedMbxKeyCode = -2015
)

var (
	// ErrBadAPIKeyFormat is the error to returns in
	// https://github.com/binance-exchange/binance-official-api-docs/blob/master/errors.md#-2014-bad_api_key_fmt
	ErrBadAPIKeyFormat = errors.New("API-key format invalid")
	// ErrRejectedMBxKey is the error to returns in
	// https://github.com/binance-exchange/binance-official-api-docs/blob/master/errors.md#-2015-rejected_mbx_key
	ErrRejectedMBxKey = errors.New("invalid API-key, IP, or permissions for action")

	// Err500 ...
	Err500 = errors.New("500 from Binance, its fault")
)

//Client represent a binance api client
type Client struct {
	APIKey      string
	SecretKey   string
	sugar       *zap.SugaredLogger
	rateLimiter Limiter
	client      *http.Client
}

//Option sets the initialization behavior for binance instance
type Option func(cl *Client) error

//WithRateLimiter alter rate limiter of binance client
func WithRateLimiter(limiter Limiter) Option {
	return func(cl *Client) error {
		cl.rateLimiter = limiter
		return nil
	}
}

//WithValidation check if API key is valid by calling GetAccountInfo with its key
func WithValidation() Option {
	return func(cl *Client) error {
		_, err := cl.GetAccountInfo()
		if err != nil {
			return fmt.Errorf("failed to validate Binance API key by calling GetAccountInfo API: err=%s", err.Error())
		}
		return nil
	}
}

//NewBinance return a new client for binance api
func NewBinance(apiKey, secretKey string, sugar *zap.SugaredLogger, options ...Option) (*Client, error) {
	client := &http.Client{
		Transport: NewTransportRateLimiter(&http.Client{Timeout: time.Second * 30}),
	}
	clnt := &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		sugar:     sugar,
		client:    client,
	}
	for _, opt := range options {
		if err := opt(clnt); err != nil {
			return nil, err
		}
	}
	//Set Default rate limiter to the limit spefified by https://api.binance.com/api/v1/exchangeInfo
	if clnt.rateLimiter == nil {
		clnt.rateLimiter = NewRateLimiter(defaultHardLimit)
	}
	return clnt, nil
}

//waitN mimic the leaky bucket algorithm to wait for n drop
func (bc *Client) waitN(n int) error {
	return bc.rateLimiter.WaitN(context.Background(), n)
}

func (bc *Client) fillRequest(req *http.Request, signNeeded bool, timepoint time.Time) error {
	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodDelete {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("User-Agent", "binance/go")
	}
	req.Header.Add("Accept", "application/json")
	if signNeeded {
		q := req.URL.Query()
		sig := url.Values{}
		req.Header.Set("X-MBX-APIKEY", bc.APIKey)
		q.Set("timestamp", fmt.Sprintf("%d", timeutil.TimeToTimestampMs(timepoint)))
		q.Set("recvWindow", "5000")
		signature, err := bc.sign(q.Encode())
		if err != nil {
			return err
		}
		sig.Set("signature", signature)
		// Using separated values map for signature to ensure it is at the end
		// of the query. This is required for /wapi apis from binance without
		// any damn documentation about it!!!
		req.URL.RawQuery = q.Encode() + "&" + sig.Encode()
	}
	return nil
}

//Sign key for authenticated api
func (bc *Client) sign(msg string) (string, error) {
	mac := hmac.New(sha256.New, []byte(bc.SecretKey))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return "", fmt.Errorf("encode message error: %s", err.Error())
	}
	result := ethereum.Bytes2Hex(mac.Sum(nil))
	return result, nil
}

func decodeErrorResponse(body io.Reader) (*ErrorResponse, error) {
	var response = &ErrorResponse{}
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

//ErrorResponse response a basic response from binance
type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (bc *Client) sendRequest(method, endpoint string, params map[string]string, signNeeded bool,
	timepoint time.Time) ([]byte, error) {

	var (
		respBody []byte
		logger   = bc.sugar.With("func", caller.GetCurrentFunctionName())
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
	if err := bc.fillRequest(req, signNeeded, timepoint); err != nil {
		return respBody, err
	}

	resp, err := bc.client.Do(req)
	if err != nil {
		return respBody, err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			logger.Errorw("Response body close error", "error", cErr.Error())
		}
	}()
	switch resp.StatusCode {
	case 429:
		err = errors.New("breaking binance request rate limit")
	case 418:
		err = errors.New("ip has been auto-banned by binance for continuing to send requests after receiving 429 codes")
	case 500:
		err = Err500
		errRsp, err := decodeErrorResponse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected response: failed to decode error response: err=%s", err.Error())
		}
		logger.Errorw("unexpected response from Binance API",
			"code", errRsp.Code,
			"msg", errRsp.Msg,
		)
	case 401:
		errRsp, err := decodeErrorResponse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected response: http_code=%d, failed to decode error response: err=%s",
				http.StatusUnauthorized,
				err.Error())
		}
		logger.Errorw("unexpected response from Binance API",
			"http_status", http.StatusUnauthorized,
			"code", errRsp.Code,
			"msg", errRsp.Msg,
		)
		// https://github.com/binance-exchange/binance-official-api-docs/blob/master/errors.md#-2014-bad_api_key_fmt
		switch errRsp.Code {
		case badAPIKeyFormatCode:
			return nil, ErrBadAPIKeyFormat
		case rejfectedMbxKeyCode:
			return nil, ErrRejectedMBxKey
		default:
			return nil, fmt.Errorf("code=%d msg=%s", errRsp.Code, errRsp.Msg)
		}
	case 200:
		respBody, err = ioutil.ReadAll(resp.Body)
	default:
		var response ErrorResponse
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			logger.Errorw("request body decode error", "error", err)
			break
		}
		err = fmt.Errorf("binance return with code: %d - %s", resp.StatusCode, response.Msg)
	}
	return respBody, err
}

//GetTradeHistory return history of trading on binance
//if fromID = -1 then function will not put fromId as a param
func (bc *Client) GetTradeHistory(symbol string, fromID uint64) ([]TradeHistory, error) {
	var (
		result []TradeHistory
	)
	const weight = 5
	//Wait before creating the request to avoid timestamp request outside the recWindow
	if err := bc.waitN(weight); err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("%s/api/v3/myTrades", endpointPrefix)
	params := map[string]string{
		"symbol": symbol,
	}

	params["fromId"] = strconv.FormatUint(fromID, 10)

	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		params,
		true,
		time.Now(),
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}

//GetAssetDetail return detail of asset
func (bc *Client) GetAssetDetail() (AssetDetailResponse, error) {
	var (
		result AssetDetailResponse
	)
	const weight = 1
	//Wait before creating the request to avoid timestamp request outside the recWindow
	if err := bc.waitN(weight); err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("%s/wapi/v3/assetDetail.html", endpointPrefix)
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{},
		true,
		time.Now(),
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}

//GetWithdrawalHistory return withdrawal history of an account
func (bc *Client) GetWithdrawalHistory(fromTime, toTime time.Time) (WithdrawHistoryList, error) {
	var (
		result WithdrawHistoryList
	)
	const weight = 1
	//Wait before creating the request to avoid timestamp request outside the recWindow
	if err := bc.waitN(weight); err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("%s/wapi/v3/withdrawHistory.html", endpointPrefix)

	params := map[string]string{}
	if !fromTime.IsZero() {
		params["startTime"] = strconv.FormatUint(timeutil.TimeToTimestampMs(fromTime), 10)
	}

	if !toTime.IsZero() {
		params["endTime"] = strconv.FormatUint(timeutil.TimeToTimestampMs(toTime), 10)
	}

	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		params,
		true,
		time.Now(),
	)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(res, &result)
	if err != nil {
		return result, err
	}
	if !result.Success {
		return result, fmt.Errorf("failed to get binance withdrawal history, reason: %s", result.Message)
	}
	return result, err
}

//GetExchangeInfo return exchange info
func (bc *Client) GetExchangeInfo() (ExchangeInfo, error) {
	var (
		result ExchangeInfo
	)
	const weight = 1
	//Wait before creating the request to avoid timestamp request outside the recWindow
	if err := bc.rateLimiter.WaitN(context.Background(), weight); err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("%s/api/v3/exchangeInfo", endpointPrefix)
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{},
		false,
		time.Now(),
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}

//GetAccountInfo return account infos
func (bc *Client) GetAccountInfo() (AccountInfo, error) {
	var (
		result AccountInfo
	)
	const weight = 5
	//Wait before creating the request to avoid timestamp request outside the recWindow
	if err := bc.waitN(weight); err != nil {
		return result, err
	}

	endpoint := fmt.Sprintf("%s/api/v3/account", endpointPrefix)

	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		nil,
		true,
		time.Now(),
	)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(res, &result)
	return result, err
}

// GetMarginTradeHistory return margin trade history
func (bc *Client) GetMarginTradeHistory(symbol string, fromID uint64) ([]TradeHistory, error) {
	var (
		result []TradeHistory
		err    error
	)
	endpoint := fmt.Sprintf("%s/sapi/v1/margin/myTrades", endpointPrefix)
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"symbol": symbol,
			"fromId": strconv.FormatUint(fromID, 10),
		},
		true,
		time.Now(),
	)
	if err == Err500 { // currently if we does not enabled margin trade, binance will return 500, ignore it (Spiros said)
		return result, nil
	}
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}

// AggregatedTrade ...
type AggregatedTrade struct {
	AggregateTradeID uint64 `json:"a"`
	Price            string `json:"p"`
	Quantity         string `json:"q"`
	FirstTradeID     uint64 `json:"f"`
	LastTradeID      uint64 `json:"l"`
	Timestamp        uint64 `json:"T"`
}

// GetAggregatedTrades ...
func (bc *Client) GetAggregatedTrades(symbol string, startTime, endTime uint64) ([]AggregatedTrade, error) {
	var (
		result []AggregatedTrade
		err    error
	)
	endpoint := fmt.Sprintf("%s/api/v3/aggTrades", endpointPrefix)
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"symbol":    symbol,
			"startTime": strconv.FormatUint(startTime, 10),
			"endTime":   strconv.FormatUint(endTime, 10),
			"limit":     "1",
		},
		false,
		time.Now(),
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}
