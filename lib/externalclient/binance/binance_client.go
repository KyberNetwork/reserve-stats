package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/lib/timeutil"
)

const (
	endpointPrefix = "https://api.binance.com"
)

//Client represent a binance api client
type Client struct {
	APIKey    string
	SecretKey string
	sugar     *zap.SugaredLogger
}

//NewBinanceClient return a new client for binance api
func NewBinanceClient(apiKey, secretKey string, sugar *zap.SugaredLogger) *Client {
	return &Client{
		APIKey:    apiKey,
		SecretKey: secretKey,
		sugar:     sugar,
	}
}

func (bc *Client) fillRequest(req *http.Request, signNeeded bool, timepoint uint64) {
	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodDelete {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("User-Agent", "binance/go")
	}
	req.Header.Add("Accept", "application/json")
	if signNeeded {
		q := req.URL.Query()
		sig := url.Values{}
		req.Header.Set("X-MBX-APIKEY", bc.APIKey)
		q.Set("timestamp", fmt.Sprintf("%d", int64(timepoint)))
		q.Set("recvWindow", "5000")
		signature, err := bc.sign(q.Encode())
		if err != nil {
			log.Fatal(err)
		}
		sig.Set("signature", signature)
		req.URL.RawQuery = q.Encode() + "&" + sig.Encode()
	}
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
	timepoint uint64) ([]byte, error) {

	var (
		respBody []byte
		logger   = bc.sugar.With("func", "binance_client/sendRequest")
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
	bc.fillRequest(req, signNeeded, timepoint)

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
			break
		}
		err = fmt.Errorf("binance return with code: %d - %s", resp.StatusCode, response.Msg)
	}
	return respBody, err
}

//GetTradeHistory return history of trading on binance
func (bc *Client) GetTradeHistory(symbol string, fromID int64) ([]TradeHistory, error) {
	var (
		result []TradeHistory
	)
	endpoint := fmt.Sprintf("%s/api/v3/myTrades", endpointPrefix)
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"symbol": symbol,
			"fromId": strconv.FormatInt(fromID, 10),
		},
		true,
		timeutil.UnixMilliSecond(),
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
	endpoint := fmt.Sprintf("%s/wapi/v3/assetDetail.html", endpointPrefix)
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{},
		true,
		timeutil.UnixMilliSecond(),
	)
	if err != nil {
		return result, err
	}
	log.Printf("%s", res)
	err = json.Unmarshal(res, &result)
	return result, err
}

//GetWithdrawalHistory return withdrawal history of an account
func (bc *Client) GetWithdrawalHistory(fromTime, toTime time.Time) (WithdrawHistoryList, error) {
	var (
		result WithdrawHistoryList
	)
	endpoint := fmt.Sprintf("%s/wapi/v3/withdrawHistory.html", endpointPrefix)
	res, err := bc.sendRequest(
		http.MethodGet,
		endpoint,
		map[string]string{
			"startTime": strconv.FormatUint(timeutil.TimeToTimestampMs(fromTime), 10),
			"endTime":   strconv.FormatUint(timeutil.TimeToTimestampMs(toTime), 10),
		},
		true,
		timeutil.UnixMilliSecond(),
	)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(res, &result)
	return result, err
}
