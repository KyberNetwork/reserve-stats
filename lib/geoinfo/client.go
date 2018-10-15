package geoinfo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var errResponseFalse = errors.New("Server return success false")

// Client is the the real implementation of geoinfo client interface
type Client struct {
	host   string
	sugar  *zap.SugaredLogger
	client *http.Client
}

type tradeLogGeoInfoResp struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
	Data    struct {
		IP      string `json:"IP"`
		Country string `json:"Country"`
	} `json:"data"`
}

const timeout = time.Minute * 5

// NewClient creates a new geoinfo client instance.
func NewClient(sugar *zap.SugaredLogger, host string) (*Client, error) {
	return &Client{
		host:   host,
		sugar:  sugar,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// GetTxInfo get ip, country info of a tx
func (g Client) GetTxInfo(tx string) (ip string, country string, err error) {
	url := fmt.Sprintf("%s/get-tx-info/%s", g.host, tx)
	resp, err := g.client.Get(url)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			g.sugar.Debugw("Response body close error", "err", cErr.Error())
		}
	}()
	response := tradeLogGeoInfoResp{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", "", err
	}
	if response.Success != true {
		g.sugar.Debugw("Get error while get info of tx", "tx", tx, "err", response.Err)
		return "", "", nil
	}
	return response.Data.IP, response.Data.Country, nil
}
