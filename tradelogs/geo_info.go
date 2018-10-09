package tradelogs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var errResponseFalse = errors.New("Server return success false")

type geoInfo struct {
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

func newGeoInfo(sugar *zap.SugaredLogger, host string) (*geoInfo, error) {
	return &geoInfo{
		host:   host,
		sugar:  sugar,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// GetTxInfo get ip, country info of a tx
func (g geoInfo) GetTxInfo(tx string) (string, string, error) {
	url := fmt.Sprintf("%s/get-tx-info/%s", g.host, tx)
	resp, err := g.client.Get(url)
	if err != nil {
		return "", "", err
	}
	response := tradeLogGeoInfoResp{}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			g.sugar.Debugw("Response body close error", "err", cErr.Error())
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", "", err
	}
	if response.Success != true {
		g.sugar.Debugw("Get error while get info of tx", "tx", tx, "err", response.Err)
		return "", "", errResponseFalse
	}
	return response.Data.IP, response.Data.Country, nil
}
