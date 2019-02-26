package huobi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
)

const (
	//HuobiEndpoint is base on
	HuobiEndpoint = "https://api.huobi.pro"
)

//Client represent a huobi client for
//calling to huobi endpoint
type Client struct {
	APIKey    string
	SecretKey string
	sugar     *zap.SugaredLogger
}

//NewHuobiClient return a new HuobiClient instance
func NewHuobiClient(apiKey, secretKey string, sugar *zap.SugaredLogger) *Client {
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		sugar:     sugar,
	}
}

func (hc *Client) sign(msg string) string {
	mac := hmac.New(sha256.New, []byte(hc.SecretKey))
	if _, err := mac.Write([]byte(msg)); err != nil {
		hc.sugar.Errorw("encode message error", "error", err.Error())
	}
	result := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return result
}

func (hc *Client) fillRequest(req *http.Request, signNeeded bool) {
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
		q.Set("Signature", hc.sign(payload))
		req.URL.RawQuery = q.Encode()
	}
}

func (hc *Client) sendRequest(method, requestURL string, params map[string]string,
	signNeeded bool) ([]byte, error) {
	var (
		logger = hc.sugar.With("func", "huobi_client/sendRequest")
	)
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
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
	if signNeeded {
		timestamp := fmt.Sprintf("%s", time.Now().UTC().Format("2006-01-02T15:04:05"))
		params["SignatureMethod"] = "HmacSHA256"
		params["SignatureVersion"] = "2"
		params["AccessKeyId"] = hc.APIKey
		params["Timestamp"] = timestamp
	}
	var sortedParams []string
	for k := range params {
		sortedParams = append(sortedParams, k)
	}
	sort.Strings(sortedParams)
	for _, k := range sortedParams {
		q.Add(k, params[k])
	}
	req.URL.RawQuery = q.Encode()

	hc.fillRequest(req, signNeeded)
	var respBody []byte
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
		break
	case 429:
		err = errors.New("breaking Huobi request rate limit")
		break
	case 500:
		err = errors.New("500 from Huobi, its fault")
		break
	case 200:
		respBody, err = ioutil.ReadAll(resp.Body)
		break
	}
	return respBody, err
}

//GetAccounts return list of accounts in huobi
func (hc *Client) GetAccounts() ([]Account, error) {
	var (
		result AccountResponse
	)
	endpoint := fmt.Sprintf("%s/v1/account/accounts", HuobiEndpoint)
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
	return result.Data, err
}

//GetTradeHistory return trade history of an account
func (hc *Client) GetTradeHistory(symbol string, startDate, endDate time.Time) (TradeHistoryList, error) {
	var (
		result TradeHistoryList
	)
	endpoint := fmt.Sprintf("%s/v1/order/orders", HuobiEndpoint)
	res, err := hc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"states":     "filled",
			"symbol":     strings.ToLower(symbol),
			"start-date": startDate.Format("2016-06-09"),
			"endDate":    endDate.Format("2018-03-08"),
		},
		true,
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}

//GetWithdrawHistory return withdraw history of an account
func (hc *Client) GetWithdrawHistory() (WithdrawHistoryList, error) {
	var (
		result WithdrawHistoryList
	)
	endpoint := fmt.Sprintf("%s/v1/query/finances/", HuobiEndpoint)
	res, err := hc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"types": "withdraw-group",
			"size":  "20",
		},
		true,
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}
