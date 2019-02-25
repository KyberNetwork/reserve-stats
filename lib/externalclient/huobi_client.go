package externalclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

//HuobiClient represent a huobi client for
//calling to huobi endpoint
type HuobiClient struct {
	APIKey    string
	SecretKey string
}

//Account represent a huobi account
type Account struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	State  string `json:"state"`
	UserID string `json:"user-id"`
}

const (
	//HuobiEndpoint is base on
	HuobiEndpoint = "https://api.huobi.pro"
)

//NewHuobiClient return a new HuobiClient instance
func NewHuobiClient(apiKey, secretKey string) *HuobiClient {
	return &HuobiClient{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}
}

//Sign return sign of request
func (hc *HuobiClient) Sign(msg string) string {
	mac := hmac.New(sha256.New, []byte(hc.SecretKey))
	if _, err := mac.Write([]byte(msg)); err != nil {
		log.Printf("Encode message error: %s", err.Error())
	}
	result := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return result
}

func (hc *HuobiClient) fillRequest(req *http.Request, signNeeded bool) {
	if req.Method == "POST" || req.Method == "PUT" || req.Method == "DELETE" {
		req.Header.Add("Content-Type", "application/json")
	} else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if signNeeded {
		q := req.URL.Query()
		sig := url.Values{}

		method := req.Method
		auth := q.Encode()
		hostname := req.URL.Hostname()
		path := req.URL.Path
		payload := strings.Join([]string{method, hostname, path, auth}, "\n")
		sig.Set("Signature", hc.Sign(payload))
		req.URL.RawQuery = q.Encode() + "&" + sig.Encode()
	}
}

func (hc *HuobiClient) sendRequest(method, requestURL string, params map[string]string,
	signNeeded bool) ([]byte, error) {

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	reqBody, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, requestURL, nil)
	if err != nil {
		return nil, err
	}
	if method == "POST" {
		req.Body = ioutil.NopCloser(strings.NewReader(string(reqBody)))
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	if signNeeded {
		timestamp := fmt.Sprintf("%s", time.Now().Format("2006-01-02T15:04:05"))
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
			log.Printf("Response body close error: %s", cErr.Error())
		}
	}()
	switch resp.StatusCode {
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
func (hc *HuobiClient) GetAccounts() ([]Account, error) {
	var (
		result []Account
	)
	endpoint := fmt.Sprintf("%s/v1/account/accounts", HuobiEndpoint)
	res, err := hc.sendRequest(
		"GET",
		endpoint,
		map[string]string{},
		true,
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}
