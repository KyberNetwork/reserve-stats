package userprofile

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-stats/lib/httputil/signer"
	ethereum "github.com/ethereum/go-ethereum/common"
)

// Client is the the implementation to query userprofile info.
type Client struct {
	sugar      *zap.SugaredLogger
	client     *http.Client
	url        string
	signingKey string
}

// NewClient creates a new core client instance.
func NewClient(sugar *zap.SugaredLogger, url, signingKey string) (*Client, error) {
	const timeout = time.Minute
	client := &http.Client{Timeout: timeout}
	return &Client{sugar: sugar, url: url, client: client, signingKey: signingKey}, nil
}

// LookUpUserProfile will look for the UserProfile of input addr from server
func (c *Client) LookUpUserProfile(addr ethereum.Address) (UserProfile, error) {
	const endpoint = "/api/wallet_info"
	var (
		params = make(map[string]string)
		result = userClientReply{
			Data: UserProfile{},
		}
	)

	params["nonce"] = signer.GenerateNonce()
	params["wallet_address"] = addr.Hex()
	req, err := signer.NewRequest(c.sugar, c.url, http.MethodGet, c.signingKey, endpoint, params)
	if err != nil {
		return result.Data, err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return result.Data, err
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return result.Data, fmt.Errorf("unexpected return code: %d", rsp.StatusCode)
	}
	if err = json.NewDecoder(rsp.Body).Decode(&result); err != nil {
		return result.Data, err
	}
	if !result.Success {
		return result.Data, fmt.Errorf("failed to get  user profile, reason : %v", result.Reason)
	}
	return result.Data, nil
}
