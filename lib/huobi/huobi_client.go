package huobi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/KyberNetwork/reserve-stats/lib/caller"
	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	//HuobiEndpoint is base on
	huobiEndpoint = "https://api.huobi.pro"
)

//Client represent a huobi client for
//calling to huobi endpoint
type Client struct {
	APIKey      string
	SecretKey   string
	sugar       *zap.SugaredLogger
	rateLimiter Limiter
}

//Option sets the initialization behavior for binance instance
type Option func(cl *Client) error

//WithRateLimiter alter ratelimiter of binance client
func WithRateLimiter(limiter Limiter) Option {
	return func(cl *Client) error {
		cl.rateLimiter = limiter
		return nil
	}
}

//WithValidation check if API key is valid by calling GetAccounts with its key
func WithValidation() Option {
	return func(cl *Client) error {
		_, err := cl.GetAccounts()
		if err != nil {
			return fmt.Errorf("failed to validate Huobi API key by calling GetAccountInfo API: err=%s", err.Error())
		}
		return nil
	}
}

//NewClient return a new HuobiClient instance
func NewClient(apiKey, secretKey string, sugar *zap.SugaredLogger, options ...Option) (*Client, error) {
	clnt := &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		sugar:     sugar,
	}
	for _, opt := range options {
		if err := opt(clnt); err != nil {
			return nil, err
		}
	}
	//Set Default rate limiter to the limit spefified by https://github.com/huobiapi/API_Docs_en/wiki/Request_Process
	if clnt.rateLimiter == nil {
		const huobiDefaultRateLimit = 10
		clnt.rateLimiter = rate.NewLimiter(rate.Limit(huobiDefaultRateLimit), 1)
	}
	return clnt, nil

}

func (hc *Client) sign(msg string) (string, error) {
	mac := hmac.New(sha256.New, []byte(hc.SecretKey))
	if _, err := mac.Write([]byte(msg)); err != nil {
		hc.sugar.Errorw("encode message error", "error", err.Error())
		return "", err
	}
	result := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return result, nil
}

func (hc *Client) fillRequest(req *http.Request, signNeeded bool) error {
	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodDelete {
		req.Header.Add("Content-Type", "application/json")
	}
	if signNeeded {
		q := req.URL.Query()
		method := req.Method
		auth := q.Encode()
		hostname := req.URL.Hostname()
		path := req.URL.Path
		payload := strings.Join([]string{method, hostname, path, auth}, "\n")
		signature, err := hc.sign(payload)
		if err != nil {
			return err
		}
		q.Set("Signature", signature)
		req.URL.RawQuery = q.Encode()
	}
	return nil
}

func (hc *Client) sendRequest(method, requestURL string, params map[string]string,
	signNeeded bool) ([]byte, error) {
	var (
		logger = hc.sugar.With("func", caller.GetCurrentFunctionName())
	)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, requestURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	//Wait before sign
	if err := hc.rateLimiter.WaitN(context.Background(), 1); err != nil {
		return []byte{}, err
	}
	if signNeeded {
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")
		params["SignatureMethod"] = "HmacSHA256"
		params["SignatureVersion"] = "2"
		params["AccessKeyId"] = hc.APIKey
		params["Timestamp"] = timestamp
	}
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	var respBody []byte
	if err := hc.fillRequest(req, signNeeded); err != nil {
		return respBody, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return respBody, err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			logger.Error("response body close error", "error", cErr.Error())
		}
	}()
	switch resp.StatusCode {
	case 404:
		err = errors.New("api not found")
	case 429:
		err = errors.New("breaking Huobi request rate limit")
	case 500:
		err = errors.New("500 from Huobi, its fault")
	case 200:
		respBody, err = ioutil.ReadAll(resp.Body)
	}
	return respBody, err
}

//GetAccounts return list of accounts in huobi
func (hc *Client) GetAccounts() ([]Account, error) {
	var (
		result AccountResponse
	)
	endpoint := fmt.Sprintf("%s/v1/account/accounts", huobiEndpoint)
	res, err := hc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{},
		true,
	)
	if err != nil {
		return result.Data, err
	}
	err = json.Unmarshal(res, &result)
	if result.Status != StatusOK.String() {
		return result.Data, fmt.Errorf("received unexpect status: err=%s code=%s msg=%s",
			result.Status,
			result.ErrorCode,
			result.ErrorMessage)
	}
	return result.Data, err
}

//GetTradeHistory return trade history of an account
//extras  params included fromID for further querrying.
//details at https://github.com/huobiapi/API_Docs_en/wiki/REST_Reference#get-v1orderorders--get-order-list
func (hc *Client) GetTradeHistory(symbol string, startDate, endDate time.Time, extras ...ExtrasTradeHistoryParams) (TradeHistoryList, error) {

	var (
		result TradeHistoryList
		params = map[string]string{
			"states":     "filled",
			"symbol":     strings.ToLower(symbol),
			"start-time": strconv.FormatUint(timeutil.TimeToTimestampMs(startDate), 10),
			"end-time":   strconv.FormatUint(timeutil.TimeToTimestampMs(endDate), 10),
		}
	)
	endpoint := fmt.Sprintf("%s/v1/order/orders", huobiEndpoint)
	res, err := hc.sendRequest(
		http.MethodGet,
		endpoint,
		params,
		true,
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	if err != nil {
		return result, err
	}

	if result.Status != StatusOK.String() {
		return result, fmt.Errorf("received unexpect status: err=%s code=%s msg=%s",
			result.Status,
			result.ErrorCode,
			result.ErrorMessage)
	}
	return result, err
}

//GetWithdrawHistory return withdraw history of an account
func (hc *Client) GetWithdrawHistory(currency string, fromID uint64) (WithdrawHistoryList, error) {
	var (
		result WithdrawHistoryList
	)
	endpoint := fmt.Sprintf("%s/v1/query/deposit-withdraw", huobiEndpoint)
	res, err := hc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"type":     "withdraw",
			"size":     "20",
			"from":     strconv.FormatUint(fromID, 10),
			"currency": strings.ToLower(currency),
		},
		true,
	)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(res, &result)
	if err != nil {
		return result, err
	}
	if result.Status != StatusOK.String() {
		return result, fmt.Errorf("received unexpect status: err=%s code=%s msg=%s",
			result.Status,
			result.ErrorCode,
			result.ErrorMessage)
	}
	return result, err
}

//GetSymbolsPair return list of pairs for Huobi's data
func (hc *Client) GetSymbolsPair() ([]Symbol, error) {
	var (
		symbolReply SymbolsReply
	)
	endpoint := fmt.Sprintf("%s/v1/common/symbols", huobiEndpoint)
	res, err := hc.sendRequest(
		http.MethodGet,
		endpoint,
		nil,
		false,
	)
	if err != nil {
		return symbolReply.Data, err
	}
	err = json.Unmarshal(res, &symbolReply)
	if err != nil {
		return symbolReply.Data, err
	}
	if symbolReply.Status != StatusOK.String() {
		return symbolReply.Data, fmt.Errorf("unexpected reply status %s", symbolReply.Status)
	}
	return symbolReply.Data, nil
}

//GetCurrencies return list of Currencies supported by Huobi
func (hc *Client) GetCurrencies() ([]string, error) {
	var (
		reply CurrenciesReply
	)
	endpoint := fmt.Sprintf("%s/v1/common/currencys", huobiEndpoint)
	res, err := hc.sendRequest(
		http.MethodGet,
		endpoint,
		nil,
		false,
	)
	if err != nil {
		return reply.Data, err
	}
	err = json.Unmarshal(res, &reply)
	if err != nil {
		return reply.Data, err
	}
	if reply.Status != StatusOK.String() {
		return reply.Data, fmt.Errorf("unexpected reply status %s", reply.Status)
	}
	return reply.Data, nil
}
