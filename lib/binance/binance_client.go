package binance

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	endpointPrefix = "https://api.binance.com"
)

//Client represent a binance api client
type Client struct {
	APIKey      string
	SecretKey   string
	Sugar       *zap.SugaredLogger
	rateLimiter Limiter
}

//Option sets the initialization behavior for binance instance
type Option func(cl *Client)

//WithRateLimiter alter ratelimiter of binance client
func WithRateLimiter(limiter Limiter) Option {
	return func(cl *Client) {
		cl.rateLimiter = limiter
	}
}

//NewBinance return a new client for binance api
func NewBinance(apiKey, secretKey string, sugar *zap.SugaredLogger, options ...Option) *Client {
	clnt := &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Sugar:     sugar,
	}
	for _, opt := range options {
		opt(clnt)
	}
	//Set Default rate limiter to the limit spefified by https://api.binance.com/api/v1/exchangeInfo
	if clnt.rateLimiter == nil {
		const binanceDefaultRateLimit = 20

		clnt.rateLimiter = rate.NewLimiter(rate.Limit(binanceDefaultRateLimit), 5)
	}
	return clnt
}

//waitN mimic the leaky bucket algorithm to wait for n drop
func (bc *Client) waitN(n int) error {
	// for i := 0; i < n; i++ {
	// 	if err := bc.rateLimiter.WaitN(context.Background(), 1); err != nil {
	// 		return err
	// 	}
	// }
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

func (bc *Client) sendRequest(method, endpoint string, params map[string]string, signNeeded bool,
	timepoint time.Time) ([]byte, error) {

	var (
		respBody []byte
		logger   = bc.Sugar.With("func", "binance_client/sendRequest")
	)
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
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

	resp, err := client.Do(req)
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
		break
	case 418:
		err = errors.New("ip has been auto-banned by binance for continuing to send requests after receiving 429 codes")
		break
	case 500:
		err = errors.New("500 from Binance, its fault")
		break
	case 401:
		err = errors.New("binance api key not valid")
		break
	case 200:
		respBody, err = ioutil.ReadAll(resp.Body)
		break
	default:
		var response APIResponse
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
func (bc *Client) GetTradeHistory(symbol string, fromID int64, fromTime, toTime time.Time) ([]TradeHistory, error) {
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

	if fromID != -1 {
		params["fromId"] = strconv.FormatInt(fromID, 10)
	}

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
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"startTime": strconv.FormatUint(timeutil.TimeToTimestampMs(fromTime), 10),
			"endTime":   strconv.FormatUint(timeutil.TimeToTimestampMs(toTime), 10),
		},
		true,
		time.Now(),
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
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

	endpoint := fmt.Sprintf("%s/api/v1/exchangeInfo", endpointPrefix)
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
