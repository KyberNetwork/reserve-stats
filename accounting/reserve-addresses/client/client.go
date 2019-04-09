package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-stats/accounting/common"
	rcommon "github.com/KyberNetwork/reserve-stats/accounting/reserve-addresses/common"
	"github.com/KyberNetwork/reserve-stats/lib/httputil"
)

// Client is the the implementation to query user kyced status info.
type Client struct {
	sugar  *zap.SugaredLogger
	client *http.Client
	url    string
}

// NewClient creates a new user kyc lookup client instance.
func NewClient(sugar *zap.SugaredLogger, url string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar,
		url:    url,
		client: client,
	}, nil
}

// ReserveAddresses Will return all the current reserve addresses in DB
func (c *Client) ReserveAddresses(filterTypes ...common.AddressType) ([]common.ReserveAddress, error) {
	const endpoint = "/addresses"
	var (
		result        rcommon.AllAddressesResponse
		reserveResult []common.ReserveAddress
		filter        = make(map[common.AddressType]struct{})
	)

	req, err := httputil.NewRequest(http.MethodGet, endpoint, c.url, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cErr := rsp.Body.Close(); cErr != nil {
			c.sugar.Errorf("failed to close body: err=%s", cErr.Error())
		}
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}
	if err = json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(filterTypes) != 0 {
		for _, typ := range filterTypes {
			filter[typ] = struct{}{}
		}

		for _, addr := range result.Data {
			if _, ok := filter[addr.Type]; ok {
				reserveResult = append(reserveResult, addr)
			}
		}
	} else {
		reserveResult = result.Data
	}

	if err = rsp.Body.Close(); err != nil {
		return nil, err
	}

	return reserveResult, nil
}
